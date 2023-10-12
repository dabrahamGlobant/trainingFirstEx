package structs

import "errors"

type StorageError struct {
	Code        string `json:"code"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ServiceError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type HTTPError struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func (e StorageError) Error() string {
	return e.Description
}

func (e ServiceError) Error() string {
	return e.Description
}

func (e HTTPError) Error() string {
	return e.Description
}

func (s ServiceError) Is(e error) bool {
	err, ok := e.(ServiceError)
	if !ok {
		return false
	}
	return err.Code == s.Code
}

// Err Codes
const (
	Internal    = "InternalError"
	NotFound    = "NotFound"
	ExsistingId = "ExistingId"
	ConFailed   = "ConnectionFailed"
	WrongFormat = "WrongRequestFormat"
)

// Err Descriptions

var (
	ErrNotFoundErr    = errors.New("not found")
	ErrExistingIdErr  = errors.New("that element already exists for : ")
	ErrWrongFormatErr = errors.New("wrong body format for : ")
	ErrConFailedErr   = errors.New("the connection failed")
	ErrJsonParse      = errors.New("json parse failed at : ")
)

// Test Err Codes

const (
	NotFoundTest = "not found"
)
