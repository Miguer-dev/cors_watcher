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
	mu       *sync.Mutex
	errorLog *log.Logger
}

func main() {
	printTitle()

	app := application{
		wg:       &sync.WaitGroup{},
		errorLog: initErrorLog(),
		mu:       &sync.Mutex{},
	}

	defer app.recoverPanic()

	go app.captureInterruptSignal()

	options := initOptions()

	transactions := initTransactions(options)

	client := createHttpClient(options)

	printGeneralOptions(options)

	for _, arrayTransactions := range transactions {
		printTableHeader(arrayTransactions[0])

		for _, transaction := range arrayTransactions {
			app.backgroundFuncWithRecover(func() {
				transaction.sendRequest(client)
				transaction.addTags()
				app.printTableTransaction(transaction)
			})

			time.Sleep(time.Duration(options.timedelay * float64(time.Second)))
		}

		app.wg.Wait()
	}

	printFile(options.output, transactions)
}
