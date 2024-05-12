package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodPost, "/v1/moduleinfo", app.createModuleInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/moduleinfo/:id", app.getModuleInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/moduleinfo/:id", app.editModuleInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/moduleinfo/:id", app.deleteModuleInfoHandler)

	// Assignment 2

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getUserInfoHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/users", app.getAllUserInfoHandler)
	// router.HandlerFunc(http.MethodPut, "/v1/users/:id", app.editUserInfoHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserInfoHandler)

	return app.recoverPanic(app.rateLimit(router))
}
