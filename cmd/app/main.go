package main

import (
	"log"
	"sync"
	"time"

	"cors_watcher/internal/vcs"
)

var (
	version = vcs.Version()
)

type application struct {
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

	transactions := initTransactions(options)

	client := createHttpClient(options)

	var url string
	for _, transaction := range transactions {
		transaction.sendRequest(client)
		transaction.addTags()

		url = transaction.printTableTransaction(url)

		time.Sleep(time.Duration(options.timedelay * float64(time.Second)))
	}
}
