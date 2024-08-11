package main

import "fmt"

// execute function on the background with recover on Panic
func (app *application) backgroundFuncWithRecover(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				printError(fmt.Sprint(err))
			}
		}()

		fn()
	}()
}
