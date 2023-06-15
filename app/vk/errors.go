package vk

import "strconv"

type Error struct {
	code   Code
	nested error
}

type Code int8

func (e *Error) Error() string {
	errorString := "unknown error, code: " + strconv.Itoa(int(e.code))

	switch e.code {
	case UnhandledError:
		errorString = "Необратываемая ошибка"
	case APIVKError:
		errorString = "vk api error"
	}
	if e.nested != nil {
		errorString += ", " + e.nested.Error()
	}

	return errorString
}

const (
	UnhandledError Code = 0
	APIVKError     Code = 1
)

func NewError(code Code, nested error) error {
	return &Error{code: code, nested: nested}
}

func (e *Error) CompareCode(c Code) bool {
	return e.code == c
}

func (e *Error) GetNested() error {
	return e.nested
}
