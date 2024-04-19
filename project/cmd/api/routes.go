package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodPost, "/v1/moduleinfo", app.createModuleInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/moduleinfo/:id", app.getModuleInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/moduleinfo/:id", app.editModuleInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/moduleinfo/:id", app.deleteModuleInfoHandler)

	return router
}
