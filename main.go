package main

import "fmt"

func main() {
	options := initOptions()
	optionsValidations := options.validateOptions()
	if !optionsValidations.Valid() {
		for key, value := range optionsValidations.Errors {
			fmt.Printf("%s: %s\n", key, value)
		}
		return
	}

	fmt.Println("End")
}
