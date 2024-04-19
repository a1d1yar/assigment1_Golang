package main

import (
	"net/http"

	"github.com/a1d1yar/assingment1_Golang/internal/data"
)

func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	var moduleInfo data.ModuleInfo
	if err := app.readJSON(w, r, &moduleInfo); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.database.Insert(&moduleInfo); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"module_info": moduleInfo}, nil)
}

func (app *application) getModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	moduleInfo, err := app.database.Retrieve(int(id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"module_info": moduleInfo}, nil)
}

func (app *application) editModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var moduleInfo data.ModuleInfo
	if err := app.readJSON(w, r, &moduleInfo); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	moduleInfo.ID = int(id)

	if err := app.database.Update(&moduleInfo); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"module_info": moduleInfo}, nil)
}

func (app *application) deleteModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.database.Delete(int(id)); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "ModuleInfo deleted successfully"}, nil)
}
