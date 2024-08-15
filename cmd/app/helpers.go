package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Print err after Panic
func (app *application) recoverPanic() {
	if err := recover(); err != nil {
		app.errorLogPrintExit(errors.New(fmt.Sprint(err)))
	}
}

// execute function on the background with recover on Panic
func (app *application) backgroundFuncWithRecover(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()
		defer app.recoverPanic()

		fn()
	}()
}

// capture interrupt signal and exit waiting for goroutines finish
func (app *application) captureInterruptSignal() {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit

	fmt.Println()
	printWarning(fmt.Sprintf("Signal: %s", s.String()))
	printWarning("Leaving ...")

	app.wg.Wait()

	os.Exit(0)
}

// read json and validate json fields
func readJSON(body io.Reader, dst any) error {

	dec := json.NewDecoder(body)

	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return errJsonSyntax(syntaxError)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errJsonUnexpectedEOF

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errJsonUnmarshalTypeField(unmarshalTypeError)
			}
			return errJsonUnmarshalType(unmarshalTypeError)

		case errors.Is(err, io.EOF):
			return errJsonEOF

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errJsonUnknownField(fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// if the r.body has more info return an error
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errJsonSingleValue
	}

	return nil
}
