package server

import (
	"bytes"
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateExecutionHandler_ServeHTTP(t *testing.T) {
	tests := map[string]struct {
		body       ExecutionBody
		statusCode int
		err        error
	}{
		"it should return bad request with empty body": {
			body:       ExecutionBody{},
			statusCode: http.StatusBadRequest,
		},
		"it should return bad request with invalid programs body": {
			body:       ExecutionBody{Daily: &Programs{}, Weekly: &WeeklyPrograms{}, Odd: &Programs{}, Even: &Programs{}},
			statusCode: http.StatusBadRequest,
		},
		"it should return internal server error when repository returns error": {
			body:       createExecutionBody(),
			statusCode: http.StatusInternalServerError,
			err:        errors.New("error"),
		},
		"it should return accepted": {
			body:       createExecutionBody(),
			statusCode: http.StatusAccepted,
		},
	}
	for name, tt := range tests {
		config := getConfig()
		router := getRouter()
		repo := execution.RepositoryMock{}
		log := logger.LoggerMock{}
		router.createExecution = newCreateExecution(execution.NewCreator(&repo, &log), &log)
		server := router.buildServer(config.AuthToken)
		t.Run(name, func(t *testing.T) {
			data, _ := jsoniter.Marshal(tt.body)
			req, err := http.NewRequest(http.MethodPut, "/executions", bytes.NewBuffer(data))
			assert.NoError(t, err)
			req.Header.Add("Authorization", config.AuthToken)

			repo.SaveFunc = func(e execution.Execution) error {
				return tt.err
			}
			log.FatalfFunc = func(format string, v ...interface{}) {
			}
			log.InfofFunc = func(format string, v ...interface{}) {
			}

			writer := httptest.NewRecorder()
			server.ServeHTTP(writer, req)

			assert.Equal(t, tt.statusCode, writer.Code)
		})
	}
}

func createExecutionBody() ExecutionBody {
	stub := execution.NewExecutionStub()
	daily := Programs{}
	weeklyPrgms := WeeklyPrograms{}
	odd := Programs{}
	even := Programs{}
	for _, dailyPrg := range *stub.Daily {
		pr := getProgram(dailyPrg)
		daily.Add(&pr)

	}
	for _, weeklyPrg := range *stub.Weekly {
		prgms := Programs{}
		for _, exec := range *weeklyPrg.Executions {
			pr := getProgram(exec)
			prgms.Add(&pr)
		}
		weeklyPrgms.Add(&Weekly{Weekday: 0, Executions: &prgms})
	}
	for _, odPrg := range *stub.Odd {
		pr := getProgram(odPrg)
		daily.Add(&pr)

	}
	for _, evenPrg := range *stub.Even {
		pr := getProgram(evenPrg)
		daily.Add(&pr)

	}
	body := ExecutionBody{
		Daily:  &daily,
		Weekly: &weeklyPrgms,
		Odd:    &odd,
		Even:   &even,
	}

	return body
}

func getProgram(prg *execution.Program) Program {
	data := ExecutionsData{Hour: prg.Executions.Hour.Format("15:04"), Zones: prg.Executions.Zones}
	pr := Program{uint8(prg.Seconds.Seconds()), &data}
	return pr
}
