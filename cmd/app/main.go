package main

import (
	"log"
	"sync"

	"cors_watcher/internal/vcs"
)

var (
	version = vcs.Version()
)

type application struct {
	requests *[]request
	wg       *sync.WaitGroup
	errorLog *log.Logger
}

func main() {
	printTitle()

	app := application{
		wg:       &sync.WaitGroup{},
		errorLog: initErrorLog(),
	}

	defer app.recoverPanic()

	go app.captureInterruptSignal()

	options := initOptions()
	optionsValidations := options.validateOptions()
	if !optionsValidations.Valid() {
		optsErrorPrintExit(optionsValidations.Errors)
	}

	requests, err := options.getRequests()
	if err != nil {
		optErrorPrintExit(err)
	}

	app.requests = requests

}
