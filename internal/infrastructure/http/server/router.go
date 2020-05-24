package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type router struct {
	homepage,
	createZone,
	getZones,
	getExecutionLogs,
	createExecution,
	getExecutions,
	executionWater,
	temperature,
	removeZone http.Handler
}

func newRouter(homepage,
	createZone,
	getZones,
	getExecutionLogs,
	createExecution,
	getExecutions,
	executionWater,
	temperature,
	removeZone http.Handler) *router {
	return &router{homepage: homepage,
		createZone:       createZone,
		getZones:         getZones,
		getExecutionLogs: getExecutionLogs,
		createExecution:  createExecution,
		getExecutions:    getExecutions,
		executionWater:   executionWater,
		temperature:      temperature,
		removeZone:       removeZone,
	}
}

func (r *router) buildServer(authToken string) *mux.Router {
	rout := mux.NewRouter()
	rout.HandleFunc("/", r.homepage.ServeHTTP).Methods(http.MethodGet)
	rout.HandleFunc("/zones", r.createZone.ServeHTTP).Methods(http.MethodPut)
	rout.HandleFunc("/zones", r.getZones.ServeHTTP).Methods(http.MethodGet)
	rout.HandleFunc("/zones/{zone_id}", r.removeZone.ServeHTTP).Methods(http.MethodDelete)
	rout.HandleFunc("/executions/logs", r.getExecutionLogs.ServeHTTP).Methods(http.MethodGet)
	rout.HandleFunc("/executions", r.createExecution.ServeHTTP).Methods(http.MethodPut)
	rout.HandleFunc("/executions", r.getExecutions.ServeHTTP).Methods(http.MethodGet)
	rout.HandleFunc("/water", r.executionWater.ServeHTTP).Methods(http.MethodPost)
	rout.HandleFunc("/temperature", r.temperature.ServeHTTP).Methods(http.MethodGet)

	middleware := newAuthMiddleware(authToken)

	rout.Use(middleware.middleware)
	return rout
}
