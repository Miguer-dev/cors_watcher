package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const version = "1.0.0"

type application struct {
	options  *options
	requests *[]request
	wg       *sync.WaitGroup
}

func main() {
	printTitle()

	options := initOptions()

	optionsValidations := options.validateOptions()
	if !optionsValidations.valid() {
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

	time.Sleep(10 * time.Second)
}

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
