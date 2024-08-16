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
	requests []request
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

	requests := options.buildRequests()

	app.requests = requests
}
