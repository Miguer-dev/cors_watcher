package main

import (
	"cors_watcher/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
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

	printInterrupt(s)

	app.wg.Wait()

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
				optErrorPrintExit(&validator.OptionError{Option: "-p", Err: err.Error()})

			} else {
				transport := &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				}

				client.Transport = transport
			}
		} else if strings.Contains(options.proxy, "socks5://") {
			dialer, err := proxy.SOCKS5("tcp", options.proxy, nil, proxy.Direct)
			if err != nil {
				optErrorPrintExit(&validator.OptionError{Option: "-p", Err: err.Error()})

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

// add port to text if url has a port
func addPortIfExist(text string, url *url.URL) string {
	if url.Port() != "" {
		return text + ":" + url.Port()
	}

	return text
}

// create string with only spaces
func spaces(num int, max int) string {
	if len(strconv.Itoa(num)) > max {
		return ""
	} else {
		return strings.Repeat(" ", max-len(strconv.Itoa(num)))
	}
}
