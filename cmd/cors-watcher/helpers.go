package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Miguer-dev/cors_watcher/internal/validator"

	"golang.org/x/net/proxy"
)

// Print err after Panic
func (app *application) recoverPanic() {
	if err := recover(); err != nil {
		switch err.(type) {
		case string:
			errorPrintExit(errors.New(err.(string)))
		case error:
			errorPrintExit(err.(error))
		}
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

// capture interrupt signal and exit
func (app *application) captureInterruptSignal() {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit

	printInterrupt(s)

	os.Exit(0)
}

// create http client to send requests
func createHttpClient(options *options) *http.Client {

	client := &http.Client{
		Timeout: time.Duration(options.timeout) * time.Second,
	}

	if options.proxy != "" {
		if strings.Contains(options.proxy, "http://") {
			proxyURL, err := url.Parse(options.proxy)
			if err != nil {
				optErrorPrintExit(&validator.OptionError{Option: "-proxy", Err: err.Error()})

			} else {
				transport := &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				}

				client.Transport = transport
			}
		} else if strings.Contains(options.proxy, "socks5://") {
			dialer, err := proxy.SOCKS5("tcp", options.proxy, nil, proxy.Direct)
			if err != nil {
				optErrorPrintExit(&validator.OptionError{Option: "-proxy", Err: err.Error()})

			} else {
				transport := &http.Transport{
					Dial: dialer.Dial,
				}

				client.Transport = transport
			}
		}
	}

	return client
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

// Build json
func writeJSON(data any) ([]byte, error) {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	js = append(js, '\n')

	return js, nil
}

// add port to text if url has a port
func addPortIfExist(text string, url *url.URL) string {
	if url.Port() != "" {
		return text + ":" + url.Port()
	}

	return text
}

// create string with only spaces
func spaces(num int64, max int) string {
	if len(strconv.FormatInt(num, 10)) > max {
		return " "
	} else {
		return strings.Repeat(" ", max-len(strconv.FormatInt(num, 10)))
	}
}
