package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// file errors
	errOpenFile = func(option string, filename string) error {
		return fmt.Errorf(`%s Unable to open "%s" file`, option, filename)
	}
	errReadFile = func(option string, filename string) error {
		return fmt.Errorf(`%s Unable to read "%s" file`, option, filename)
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
