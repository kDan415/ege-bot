package main

import "strconv"

type InitializationError struct {
	code   Code
	nested error
}

type Code int8

func (e *InitializationError) Error() string {
	switch e.code {
	case MissingConfigFile:
		return "missing config file"
	case ReadConfigError:
		return "read config error"
	case UnmarshalConfigError:
		return "unmarshal config error"
	case MarshalConfigError:
		return "marshal config error"
	case WriteConfigError:
		return "write config error"
	case ConfigIsNotFiledError:
		return "config is not filed"
	}

	return "unknown error, code: " + strconv.Itoa(int(e.code))
}

const (
	MissingConfigFile     Code = 1
	ReadConfigError       Code = 2
	UnmarshalConfigError  Code = 3
	MarshalConfigError    Code = 4
	WriteConfigError      Code = 5
	ConfigIsNotFiledError Code = 6
)

func NewInitializationError(code Code, nested error) error {
	return &InitializationError{code: code, nested: nested}
}

func (e *InitializationError) CompareCode(c Code) bool {
	return e.code == c
}

func (e *InitializationError) GetNestedError() error {
	return e.nested
}
