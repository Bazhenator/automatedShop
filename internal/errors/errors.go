package errors

import (
	"errors"
	"net/http"
)

// TODO: rewrite http errors to sql errors

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrTooSmallPwdLen = errors.New("password must contain at least 5 symbols")
)

var ErrHttpInternal = errors.New("some internal error happened")
var ErrHttpConflict = errors.New("server state conflict")
var ErrHttpTimeout = errors.New("request timeout")
var ErrHttpBadGateway = errors.New("request timeout")
var ErrHttpGatewayTimeout = errors.New("gateway timeout")
var ErrHttpBadRequest = errors.New("probably bad request")
var ErrHttpInvalidRequestData = errors.New("some data is invalid")

type JSONError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	} `json:"error"`
}

func GetJSONError(msg string, err error) *JSONError {
	var jsonErr = JSONError{}
	jsonErr.Error.Message = msg
	switch {
	case errors.Is(err, ErrHttpInternal):
		jsonErr.Error.Code = http.StatusInternalServerError
	case errors.Is(err, ErrHttpBadRequest):
		jsonErr.Error.Code = http.StatusBadRequest
	case errors.Is(err, ErrHttpInvalidRequestData):
		jsonErr.Error.Code = http.StatusUnprocessableEntity
	case errors.Is(err, ErrHttpConflict):
		jsonErr.Error.Code = http.StatusConflict
	case errors.Is(err, ErrHttpTimeout):
		jsonErr.Error.Code = http.StatusRequestTimeout
	case errors.Is(err, ErrHttpGatewayTimeout):
		jsonErr.Error.Code = http.StatusGatewayTimeout
	case errors.Is(err, ErrHttpBadGateway):
		jsonErr.Error.Code = http.StatusBadGateway
	default:
		jsonErr.Error.Code = http.StatusInternalServerError
	}

	return &jsonErr
}
