package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/a1d1yar/assingment1_Golang/internal/data"
	"github.com/a1d1yar/assingment1_Golang/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	// Fetch user info from the database based on the ID
	userInfo, err := app.db.GetUserInfoByID(id)
	if err != nil {
		// Handle error
		app.serverErrorResponse(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(userInfo)
	if err != nil {

		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) getAllUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	userInfos, err := app.db.GetAllUserInfo()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"userInfos": userInfos}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// editUserInfoHandler handles PUT or Patch requests to edit user info by ID.
func (app *application) editUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	// Decode JSON payload into a new data.UserInfo struct
	var userInfo data.UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		// Handle error
		app.badRequestResponse(w, r, err)
		return
	}

	// Update user info in the database
	err = app.db.UpdateUserInfo(id, &userInfo)
	if err != nil {
		// Handle error
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with success message
	app.writeJSON(w, http.StatusOK, envelope{"message": "User info updated successfully"}, nil)
}
func (app *application) deleteUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	// Delete user info from the database
	err := app.db.DeleteUserInfo(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "User info deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
