package server

type daemon struct {
	execution       *executionDaemon
	executionInTime *executionInTimeDaemon
	statusSetter    *statusSetterDaemon
	weather         *weatherDaemon
}

func newDaemon(execution *executionDaemon,
	executionInTime *executionInTimeDaemon,
	statusSetter *statusSetterDaemon,
	weather *weatherDaemon) *daemon {
	return &daemon{execution: execution, executionInTime: executionInTime, statusSetter: statusSetter, weather: weather}
}
