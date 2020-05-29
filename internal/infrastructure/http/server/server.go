package server

import (
	"context"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/disk/file"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/gpio/relay"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/gpio/temperature"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/rainsensor"
	logger2 "github.com/bruli/raspberryWaterSystem/internal/infrastructure/log/logger"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/mysql"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/telegram"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	"github.com/bruli/raspberryWaterSystem/internal/rain"
	relay2 "github.com/bruli/raspberryWaterSystem/internal/relay"
	"github.com/bruli/raspberryWaterSystem/internal/status"
	"github.com/bruli/raspberryWaterSystem/internal/weather"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	wait        time.Duration
	srv         *http.Server
	daemon      *daemon
	executionCh chan executionData
	log         logger.Logger
	deactivate  *relay2.Deactivate
}

func NewServer(conf *Config) *Server {
	log := getLogger(conf)
	log.Info("Water system started...")
	st := status.New()
	zoneRepository := file.NewZoneRepository(conf.ZonesFile)
	executionRepo := file.NewExecutionRepository(conf.ExecutionsFile)
	mysqlConfig := mysql.NewConfig(conf.MysqlHost, conf.MysqlPort, conf.MysqlUser, conf.MysqlPass, conf.MysqlDatabase)
	mysqlRepository := mysql.NewRepository(mysqlConfig)
	relayManager := relay.NewManager()
	executionLogRepository := mysql.NewExecutionLogRepository(mysqlRepository)
	notificationSender := notificationSender(conf, log)

	homepage := newHomePage(log, st)
	createZone := newCreateZone(zone.NewCreator(
		zoneRepository,
		relay.NewZoneRelayRepository(),
		log),
		log)
	getZones := newGetZones(zone.NewGetter(zoneRepository), log)
	getExecutionLogs := newGetExecutionLogs(
		execution.NewReadLogs(mysql.NewExecutionLogRepository(mysqlRepository), log),
		log)
	createExecution := newCreateExecution(
		execution.NewCreator(executionRepo, log),
		log)
	getExecutions := newGetExecutions(execution.NewGetter(executionRepo, log), log)
	executionCh := make(chan executionData)
	executionWater := NewExecutionWater(executionCh, log)
	weatherRepo := getWeatherRepository(conf, log)
	getter := weather.NewGetter(weatherRepo)
	temperat := newGetTemperature(getter, log, st)
	remover := zone.NewRemover(zoneRepository, log)
	removeZon := newRemoveZone(remover, log)

	rout := newRouter(homepage,
		createZone,
		getZones,
		getExecutionLogs,
		createExecution,
		getExecutions,
		executionWater,
		temperat,
		removeZon)
	executor := execution.NewExecutor(zoneRepository, relayManager, executionLogRepository, notificationSender)
	executorInTime := execution.NewExecutorInTime(executionRepo, executor, st, notificationSender)
	rainRepo := getRainRepository(conf, log)
	rainRead := rain.NewReader(rainRepo)
	weatherStatusSet := weather.NewStatusSetter(st, weatherRepo, rainRead)
	weatherWriteRepo := mysql.NewWeatherRepository(mysqlRepository)
	writerWeath := weather.NewWriter(weatherRepo, weatherWriteRepo)
	d := newDaemon(
		newExecutionDaemon(log, executor),
		newExecutionInTimeDaemon(executorInTime, log),
		newStatusSetterDaemon(log, weatherStatusSet),
		newWeatherDaemon(log, writerWeath),
	)

	duration := 15 * time.Second
	return &Server{
		wait: duration,
		srv: &http.Server{
			Handler:      rout.buildServer(conf.AuthToken),
			Addr:         conf.ServerURL,
			WriteTimeout: duration,
			ReadTimeout:  duration,
		},
		daemon:      d,
		executionCh: executionCh,
		deactivate:  relay2.NewDeactivate(relayManager),
		log:         log,
	}
}

func getWeatherRepository(conf *Config, log logger.Logger) weather.Repository {
	if conf.devMode {
		return temperature.NewInMemoryReader(log)
	}
	return temperature.NewReader()
}

func getRainRepository(conf *Config, log logger.Logger) rain.Repository {
	if conf.devMode {
		return rainsensor.NewInMemoryRepository(log)
	}
	return rainsensor.NewRepository(conf.rainSensorServerUrl)
}

func getLogger(conf *Config) *logger2.Logger {
	log := logger2.New()
	if conf.devMode {
		log.EnableDebug()
	}
	return log
}

func notificationSender(conf *Config, log logger.Logger) execution.NotificationSender {
	if conf.devMode {
		return telegram.NewInMemorySender(log)
	}
	return telegram.NewSender(conf.telegramToken, conf.telegramChatID)
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}()

	go func() {
		if err := s.deactivate.Deactivate(); err != nil {
			s.log.Fatal(err)
		}
		s.log.Debug("relays deactivated")
	}()
	go s.daemon.execution.execute(ctx, s.executionCh)
	go s.daemon.executionInTime.execute(ctx)
	go s.daemon.statusSetter.execute(ctx)
	go s.daemon.weather.execute(ctx)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	cancel()

	// createZone a deadline to wait for.
	ctxShutdown, cancelShoutDown := context.WithTimeout(context.Background(), s.wait)
	defer cancelShoutDown()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = s.srv.Shutdown(ctxShutdown)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}
