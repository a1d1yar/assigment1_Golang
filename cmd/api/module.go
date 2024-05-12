package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a1d1yar/assingment1_Golang/internal/data"
	"github.com/a1d1yar/assingment1_Golang/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	var moduleInfo data.ModuleInfo
	if err := app.readJSON(w, r, &moduleInfo); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.db.Insert(&moduleInfo); err != nil {
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

	moduleInfo, err := app.db.Retrieve(int(id))
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

	if err := app.db.Update(&moduleInfo); err != nil {
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

	if err := app.db.Delete(int(id)); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "ModuleInfo deleted successfully"}, nil)
}

// Defence

func (app *application) createDepartmentInfoHandler(w http.ResponseWriter, r *http.Request) {
	var departmentInfo data.DepartmentInfo
	if err := app.readJSON(w, r, &departmentInfo); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.db.InsertDepartmentInfo(&departmentInfo); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"department_info": departmentInfo}, nil)
}

func (app *application) getDepartmentInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	departmentInfo, err := app.db.RetrieveDepartmentInfo(int(id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"department_info": departmentInfo}, nil)
}

// 2 assingment

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
	// token generation to activate account
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {

		//
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		fmt.Print(token.Plaintext)

		// sending context data to template page
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.serverErrorResponse(w, r, err)
			return
			// app.logger.PrintError(err, nil)
		}
	})
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	// StatusAccepted - request accepted for processing but not completed yet
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
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
