package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	user "first-ex/internal/repos"
	"first-ex/internal/structs"
)

func GetAll(sv user.UserService) func(wr http.ResponseWriter, r *http.Request) {
	fn := func(wr http.ResponseWriter, r *http.Request) {
		users, err := sv.GetAll()

		wr.Header().Set("Content-Type", "application.json")

		if err != nil {
			//sendError(wr, http.StatusNotFound, err)
			return
		}
		wr.WriteHeader(http.StatusFound)
		json.NewEncoder(wr).Encode(users)
	}
	return fn
}

func Get(sv user.UserService) func(wr http.ResponseWriter, r *http.Request) {
	fn := func(wr http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id := params["id"]

		wr.Header().Set("Content-Type", "application.json")

		uid, uidErr := uuid.Parse(id)

		if uidErr != nil {
			// Error parsing the id in the correct way
			return
		}
		user, err := sv.Get(uid)

		if err != nil {
			//sendError(wr, http.StatusNotFound, err)
			return
		}
		wr.WriteHeader(http.StatusFound)
		json.NewEncoder(wr).Encode(user)
	}
	return fn
}

func Post(sv user.UserService) func(wr http.ResponseWriter, r *http.Request) {
	fn := func(wr http.ResponseWriter, req *http.Request) {

		// Possible Improvements: To make the validation better, is possible to implement third party validators like validator V10

		wr.Header().Set("Content-Type", "application.json")
		var reqBody structs.UserRequest
		decodedBody := json.NewDecoder(req.Body)
		if err := decodedBody.Decode(&reqBody); err != nil {
			//Print an error
			return
		}

		newUsr, err := sv.Create(reqBody)
		if err != nil {
			// Print error
			return
		}
		wr.WriteHeader(http.StatusCreated)
		json.NewEncoder(wr).Encode(newUsr)

	}
	return fn
}

func Delete(sv user.UserService) func(wr http.ResponseWriter, r *http.Request) {
	fn := func(wr http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id := params["id"]

		wr.Header().Set("Content-Type", "application/json")

		uid, uidErr := uuid.Parse(id)

		if uidErr != nil {
			// Error parsing the ID in the correct way, return an error response
			wr.WriteHeader(http.StatusBadRequest)
			return
		}

		err := sv.Delete(uid)

		if err != nil {
			// Handle the error, return an error response
			wr.WriteHeader(http.StatusNotFound)
			return
		}

		wr.WriteHeader(http.StatusNoContent)
	}
	return fn

}
func Put(sv user.UserService) func(wr http.ResponseWriter, req *http.Request) {
	fn := func(wr http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id := params["id"]

		wr.Header().Set("Content-Type", "application/json")

		uid, uidErr := uuid.Parse(id)

		if uidErr != nil {
			// Error parsing the ID in the correct way, return an error response
			wr.WriteHeader(http.StatusBadRequest)
			return
		}

		var reqBody structs.UserRequest
		decodedBody := json.NewDecoder(req.Body)
		if err := decodedBody.Decode(&reqBody); err != nil {
			// Handle the error (bad request), return an error response
			wr.WriteHeader(http.StatusBadRequest)
			return
		}

		isValid := validator.New()
		wrongBody := isValid.Struct(req)

		if wrongBody != nil {
			// Add error handler in case of body input error
			return
		}

		updUser := structs.User{
			ID:       uid,
			Name:     reqBody.Name,
			LastName: reqBody.LastName,
			Email:    reqBody.Email,
			Active:   reqBody.Active,
			Address:  structs.Address(reqBody.Address),
		}

		updatedUser, err := sv.Update(uid, updUser)

		if err != nil {
			// Handle the error (user not found or other), return an error response
			wr.WriteHeader(http.StatusNotFound)
			return
		}

		wr.WriteHeader(http.StatusOK)
		json.NewEncoder(wr).Encode(updatedUser)
	}
	return fn
}
