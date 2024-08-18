package main

import (
	"fmt"
	"log"
	"net/http"
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

	client := &http.Client{
		Timeout: time.Duration(options.timeout) * time.Second,
	}

	for _, transaction := range transactions {
		transaction.sendRequest(client)
		time.Sleep(time.Duration(options.timedelay * float64(time.Second)))

		fmt.Println(transaction.name)
		fmt.Println(transaction.tags)
		fmt.Println(transaction.request.URL)
		fmt.Println(transaction.request.Method)
		fmt.Println(transaction.request.Headers)
		fmt.Println(transaction.request.Data)
		fmt.Println(transaction.response.length)
		fmt.Println(transaction.response.statusCode)
		fmt.Println(transaction.response.ACDetected)
		fmt.Println(transaction.response.ACAO)
		fmt.Println(transaction.response.ACAC)
		fmt.Println(transaction.err)
		fmt.Println()
	}
}
