package services

import (
	"net/http"
)

type Logger interface {
	Debug(Message string, Context ...string)
	Info(Message string, Context ...string)
	Warning(Message string, Context ...string)
	Error(Error error, Context ...string)
	Critical(Error error, Context ...string)
}

type IdCreator interface {
	Create() string
}

type Validator interface {
	ValidateStruct(modelData any) []error
	translateError(errorData *ErrorResponse) error
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Error struct {
	Status  int    `json:"-"`
	Message string `json:"-"`
	Error   string `json:"error"`
}

type BaseService struct {
	Logger Logger
	Error  *Error
	Ulid   IdCreator
}

type Request interface {
}

func (bs *BaseService) BadRequest(_ Request, err error) {
	bs.Error = &Error{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Error:   err.Error(),
	}
}

func (bs *BaseService) ErrorHandler() {
	rec := recover()
	if rec != nil {
		bs.InternalServerError(rec)
	}

}

func (bs *BaseService) InternalServerError(_ any) {
	bs.Error = &Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Error:   "internal_server_error",
	}

}
