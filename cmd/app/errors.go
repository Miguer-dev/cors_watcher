package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type optionError struct {
	option string
	err    error
}

var (
	// file errors
	errOpenFile = func(filename string) error {
		return fmt.Errorf(`Unable to open "%s" file`, filename)
	}
	errReadFile = func(filename string) error {
		return fmt.Errorf(`Unable to read "%s" file`, filename)
	}

	// json errors
	errJsonSyntax = func(syntaxError *json.SyntaxError) error {
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	}
	errJsonUnexpectedEOF = errors.New("body contains badly-formed JSON")
	errJsonUnmarshalType = func(unmarshalTypeError *json.UnmarshalTypeError) error {
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	}
	errJsonUnmarshalTypeField = func(unmarshalTypeError *json.UnmarshalTypeError) error {
		return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
	}
	errJsonEOF          = errors.New("body must not be empty")
	errJsonUnknownField = func(fieldName string) error {
		return fmt.Errorf("body contains unknown key %s", fieldName)
	}
	errJsonSingleValue = errors.New("body must only contain a single JSON value")
)
