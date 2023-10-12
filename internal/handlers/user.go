package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
			sendHttpError(wr, castHttpError(err))
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
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: uidErr.Error(),
			})
			return
		}
		user, err := sv.Get(uid)

		if err != nil {
			sendHttpError(wr, castHttpError(err))
			return
		}
		wr.WriteHeader(http.StatusFound)
		json.NewEncoder(wr).Encode(user)
	}
	return fn
}

func Post(sv user.UserService) func(wr http.ResponseWriter, r *http.Request) {
	fn := func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Content-Type", "application.json")
		var reqBody structs.UserRequest
		decodedBody := json.NewDecoder(req.Body)
		if err := decodedBody.Decode(&reqBody); err != nil {
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			})
			return
		}

		newUsr, err := sv.Create(reqBody)
		if err != nil {
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			})
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
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: uidErr.Error(),
			})
			return
		}

		err := sv.Delete(uid)

		if err != nil {
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.ExsistingId,
				Status:      http.StatusNotFound,
				Description: err.Error(),
			})
			return
		}

		wr.WriteHeader(http.StatusNoContent)
		json.NewEncoder(wr).Encode(map[string]string{
			"response": fmt.Sprintf("User Deleted ID: %s", id),
		})
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
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.ExsistingId,
				Status:      http.StatusNotFound,
				Description: uidErr.Error(),
			})
			return
		}

		var reqBody structs.UserRequest
		decodedBody := json.NewDecoder(req.Body)
		if err := decodedBody.Decode(&reqBody); err != nil {
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			})
			wr.WriteHeader(http.StatusBadRequest)
			return
		}

		isValid := validator.New()
		wrongBody := isValid.Struct(req)

		if wrongBody != nil {
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.WrongFormat,
				Status:      http.StatusBadRequest,
				Description: wrongBody.Error(),
			})
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
			sendHttpError(wr, structs.HTTPError{
				Code:        structs.NotFound,
				Status:      http.StatusNotFound,
				Description: err.Error(),
			})
			wr.WriteHeader(http.StatusNotFound)
			return
		}

		wr.WriteHeader(http.StatusOK)
		json.NewEncoder(wr).Encode(updatedUser)
	}
	return fn
}
func castHttpError(err error) structs.HTTPError {
	var httpStatus int

	errAssert, ok := err.(structs.ServiceError)
	if !ok {
		httpStatus = http.StatusInternalServerError
	} else {
		switch {
		case errors.Is(err, structs.ServiceError{Code: structs.NotFound}):
			httpStatus = http.StatusNotFound
		case errors.Is(err, structs.ServiceError{Code: structs.ConFailed}):
			httpStatus = http.StatusInternalServerError
		case errors.Is(err, structs.ServiceError{Code: structs.ExsistingId}):
			httpStatus = http.StatusConflict
		default:
			httpStatus = http.StatusInternalServerError
		}
	}

	return structs.HTTPError{
		Code:        errAssert.Code,
		Status:      httpStatus,
		Description: errAssert.Error(),
	}
}

func sendHttpError(w http.ResponseWriter, httpError structs.HTTPError) {
	response := map[string]interface{}{
		"message": httpError.Error(),
		"code":    httpError.Code,
	}

	w.WriteHeader(httpError.Status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
