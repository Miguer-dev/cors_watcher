package main

import (
	"os"
	"sync"

	"cors_watcher/internal/vcs"
)

var (
	version = vcs.Version()
)

type application struct {
	options  *options
	requests *[]request
	wg       *sync.WaitGroup
}

func main() {
	defer recoverPanic()

	options := initOptions()

	printTitle()

	optionsValidations := options.validateOptions()
	if !optionsValidations.Valid() {
		printOptionsErrors(optionsValidations.Errors)
		os.Exit(1)
	}

	requests, err := options.getRequests()
	if err != nil {
		err.printOptionError()
		os.Exit(1)
	}

	app := application{
		options:  options,
		requests: requests,
		wg:       &sync.WaitGroup{},
	}

	go app.captureInterruptSignal()

}
