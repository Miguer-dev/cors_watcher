package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/Miguer-dev/cors_watcher/internal/validator"
)

var (
	// general errors
	errDefault = errors.New("Program terminated due to a fatal error")

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

// Create log file for errors and setup log format
func initErrorLog() *log.Logger {
	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	errorLog := log.New(errorFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return errorLog
}

// Log error details, print error tittle and exit
func (app *application) errorLogPrintExit(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	printError(errDefault.Error())
	os.Exit(1)
}

// Print several options errors and exit
func optsErrorPrintExit(err []*validator.OptionError) {
	for _, value := range err {
		printOptionError(value)
	}

	os.Exit(1)
}

// Print option error and exit
func optErrorPrintExit(err *validator.OptionError) {
	printOptionError(err)
	os.Exit(1)
}
