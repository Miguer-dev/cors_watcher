package main

import (
	"fmt"
	"log"
	"sync"

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

	for _, transaction := range transactions {

		fmt.Println(transaction.name)
		fmt.Println(transaction.tags)
		fmt.Println(transaction.request.URL)
		fmt.Println(transaction.request.Method)
		fmt.Println(transaction.request.Headers)
		fmt.Println(transaction.request.Data)
		fmt.Println()

	}
}
