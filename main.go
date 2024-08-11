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
	options *options
	wg      *sync.WaitGroup
}

func main() {
	title()

	options := initOptions()

	app := application{
		options: options,
		wg:      &sync.WaitGroup{},
	}

	optionsValidations := options.validateOptions()
	if !optionsValidations.valid() {
		printOptionsErrors(optionsValidations.Errors)
		os.Exit(1)
	}

	go func() {
		quit := make(chan os.Signal, 1)
		defer close(quit)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		fmt.Println()
		printWarning(fmt.Sprintf("Signal: %s", s.String()))
		printWarning("Leaving ...")

		app.wg.Wait()

		os.Exit(0)
	}()

	time.Sleep(10 * time.Second)
}
