package main

import "fmt"

const version = "1.0.0"

func main() {
	options := initOptions()
	optionsValidations := options.validateOptions()
	if !optionsValidations.valid() {
		for key, value := range optionsValidations.Errors {
			fmt.Printf("%s: %s\n", key, value)
		}
		return
	}

	fmt.Println("End")
}
