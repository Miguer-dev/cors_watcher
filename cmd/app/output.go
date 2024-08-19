package main

import (
	"cors_watcher/internal/validator"
	"fmt"

	"github.com/fatih/color"
)

var errorFormat = color.New(color.FgRed, color.Bold).PrintfFunc()
var warningFormat = color.New(color.FgYellow, color.Bold).PrintfFunc()
var successFormat = color.New(color.FgGreen, color.Bold).PrintfFunc()
var infoFormat = color.New(color.FgHiWhite, color.Bold).PrintfFunc()
var subtextFormat = color.New(color.FgHiBlack).PrintlnFunc()
var highlightFormat = color.New(color.FgCyan, color.Bold).PrintFunc()
var headerFormat = color.New(color.BgBlue, color.FgBlack, color.Bold).PrintFunc()
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

func (transaction *transaction) printTableHeader() {
	fmt.Println()
	headerFormat(" " + transaction.request.URL + " ")
	fmt.Print(" ")
	headerFormat(" " + transaction.request.Method + " ")
	fmt.Println()
	fmt.Println("+------+------+-------------")
	fmt.Println("|STATUS| SIZE |   Origin    ")
	fmt.Println("+------+------+-------------")
}

func (transaction *transaction) printTableTransaction(url string) string {
	if url != transaction.request.URL {
		transaction.printTableHeader()
		url = transaction.request.URL
	}

	status := transaction.response.statusCode
	len := transaction.response.length

	fmt.Printf("| %d%s| %d%s| %s", status, spaces(status, 5), len, spaces(len, 5), transaction.request.Headers["Origin"])

	for _, tag := range transaction.tags {
		fmt.Print(" ")
		tag.print(tag.info)
	}

	fmt.Println()
	return url
}

// strconv.FormatBool(transaction.response.ACDetected), transaction.response.ACAO, transaction.response.ACAC}
//"AC Detected", "ACAO", "ACAC"
