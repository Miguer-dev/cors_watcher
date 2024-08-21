package main

import (
	"fmt"
	"os"

	"github.com/Miguer-dev/cors_watcher/internal/validator"

	"github.com/fatih/color"
)

var errorFormat = color.New(color.FgRed, color.Bold).PrintfFunc()
var warningFormat = color.New(color.FgYellow, color.Bold).PrintfFunc()
var successFormat = color.New(color.FgGreen, color.Bold).PrintfFunc()
var infoFormat = color.New(color.FgHiWhite, color.Bold).PrintfFunc()
var subtextFormat = color.New(color.FgHiBlack).PrintlnFunc()
var highlightFormat = color.New(color.FgCyan, color.Bold).PrintFunc()
var headerFormat = color.New(color.BgWhite, color.FgBlack, color.Bold).PrintFunc()
var redBackgroundFormat = color.New(color.BgRed, color.FgBlack, color.Bold).PrintFunc()
var yellowBackgroundFormat = color.New(color.BgYellow, color.FgBlack, color.Bold).PrintFunc()
var greenBackgroundFormat = color.New(color.BgGreen, color.FgBlack, color.Bold).PrintFunc()
var cyanBackgroundFormat = color.New(color.BgCyan, color.FgBlack, color.Bold).PrintFunc()

// print text with error format
func printError(err string) {
	errorFormat("[x] ")
	infoFormat(err + "\n")
}

// print option error
func printOptionError(optErr *validator.OptionError) {
	errorFormat("[x] ")
	highlightFormat(optErr.Option)
	infoFormat(" " + optErr.Err + "\n")
}

// print text with info format
func printInfo(text string) {
	successFormat("[+] ")
	infoFormat(text + "\n")
}

// print text with warning format
func printWarning(text string) {
	warningFormat("[!] ")
	infoFormat(text + "\n")
}

// print interrupt signal
func printInterrupt(s os.Signal) {
	fmt.Println()
	fmt.Println()
	printWarning(fmt.Sprintf("Signal: %s", s.String()))
	printWarning("Leaving ...")
}

// print app title
func printTitle() {
	fmt.Println()
	fmt.Println()
	fmt.Println("█▀▀ █▀█ █▀█ █▀▀  █░█░█ ░▄▄ █▄ ▄▄▄ █▄▄ ▄▄░ ▄▄")
	fmt.Println("█▄▄ █▄█ █▀▄ ▄▄█  █▄█▄█ ▀▄█ █▄ █▄▄ █░█ ██▄ █░")
	fmt.Println("░░░ ░░░ ░░░ ░░░  ░░░░░ ░░░ ░░ ░░░ ░░░ ░░░ ░░")
	subtextFormat("                               by Miguer-dev")
	fmt.Println()
}

// print options common to all request
func printGeneralOptions(options *options) {
	printInfo(fmt.Sprintf("Timeout: %d", options.timeout))
	printInfo(fmt.Sprintf("Delay: %.1f", options.timedelay))

	if options.proxy != "" {
		printInfo("Proxy: " + options.proxy)
	}
}

// print table hader and request info
func printTableHeader(transaction *transaction) {
	fmt.Println()
	printInfo("URL: " + transaction.request.URL + " ")
	printInfo("Method: " + transaction.request.Method + " ")

	if len(transaction.request.Headers) > 1 {

		headers := "Headers: {"
		for key, value := range transaction.request.Headers {
			if key != "Origin" {
				headers += fmt.Sprintf("%s: %s, ", key, value)
			}
		}
		headers = headers[:len(headers)-2]
		headers += "}"

		printInfo(headers)
	}

	if transaction.request.Data != "" && transaction.request.Method != "GET" {
		printInfo("Data: " + transaction.request.Data + " ")
	}

	fmt.Println("+------+------+-------------")
	fmt.Println("|STATUS| SIZE |   ORIGIN    ")
	fmt.Println("+------+------+-------------")
}

// print transaction has table row
func (app *application) printTableTransaction(transaction *transaction) {
	status := transaction.response.statusCode
	len := transaction.response.length

	app.mu.Lock()
	defer app.mu.Unlock()

	fmt.Printf("| %d%s| %d%s| %s", status, spaces(status, 5), len, spaces(len, 5), transaction.request.Headers["Origin"])

	for _, tag := range transaction.tags {
		fmt.Print(" ")
		tag.print(tag.info)
	}

	fmt.Println()
}

// save output in file
func printFile(filename string, transactions [][]*transaction) {
	if filename != "" {

		file, err := os.Create(filename)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output", Err: err.Error()})
		}
		defer file.Close()

		var text string

		for _, arrayTransactions := range transactions {

			text += "[+] URL: " + arrayTransactions[0].request.URL + "\n"
			text += "[+] Method: " + arrayTransactions[0].request.Method + "\n"

			if len(arrayTransactions[0].request.Headers) > 1 {

				headers := "[+] Headers: {"
				for key, value := range arrayTransactions[0].request.Headers {
					if key != "Origin" {
						headers += fmt.Sprintf("%s: %s, ", key, value)
					}
				}
				headers = headers[:len(headers)-2]
				headers += "}\n"

				text += headers
			}

			if arrayTransactions[0].request.Data != "" && arrayTransactions[0].request.Method != "GET" {
				text += "[+] Data: " + arrayTransactions[0].request.Data + "\n"
			}

			text += "+------+------+-------------\n"
			text += "|STATUS| SIZE |   ORIGIN    \n"
			text += "+------+------+-------------\n"

			for _, transaction := range arrayTransactions {
				status := transaction.response.statusCode
				len := transaction.response.length

				text += fmt.Sprintf("| %d%s| %d%s| %s", status, spaces(status, 5), len, spaces(len, 5), transaction.request.Headers["Origin"])

				for _, tag := range transaction.tags {
					text += fmt.Sprint(" ")
					text += tag.info
				}

				text += "\n"
			}
			text += "\n"
		}

		_, err = file.WriteString(text)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output", Err: err.Error()})
		}

		fmt.Println()
		printWarning("Saving output ...")
		printInfo(fmt.Sprintf(`Output successfully saved in “%s”`, filename))
	}
}
