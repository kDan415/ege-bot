package ege

import (
	"net/http"
	"strconv"
)

type Error struct {
	code     Code
	httpCode int
	nested   error
}

type Code int8

func (e *Error) Error() string {
	errorString := "unknown error, code: " + strconv.Itoa(int(e.code))

	switch e.code {
	case UnhandledError:
		errorString = "unhandled error"
	case ConstructorError:
		errorString = "constructor error"
	case ServerError:
		errorString = "server error"
	case JsonDecodeError:
		errorString = "json decode error"
	case CreateRequestError:
		errorString = "create request error"
	case HttpError:
		errorString = "http error"
	case NotificatorIsAlreadyOn:
		return "Менеджер уведомлений уже включен"
	case NotificatorIsAlreadyOff:
		return "Менеджер уведомлений уже выключен"
	}

	if e.httpCode != 0 {
		errorString += ": " + http.StatusText(e.httpCode)
	}

	if e.nested != nil {
		errorString += ", " + e.nested.Error()
	}

	return errorString
}

const (
	UnhandledError          Code = 0
	ConstructorError        Code = 1
	ServerError             Code = 2
	CreateRequestError      Code = 3
	JsonDecodeError         Code = 4
	HttpError               Code = 5
	NotificatorIsAlreadyOn  Code = 6
	NotificatorIsAlreadyOff Code = 7
)

func NewError(code Code, httpCode int, nested error) error {
	return &Error{code: code, httpCode: httpCode, nested: nested}
}

func (e *Error) CompareCode(c Code) bool {
	return e.code == c
}

func (e *Error) GetNested() error {
	return e.nested
}
