package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

var errorFormat = color.New(color.FgRed, color.Bold).PrintfFunc()
var warningFormat = color.New(color.FgYellow, color.Bold).PrintfFunc()
var successFormat = color.New(color.FgGreen, color.Bold).PrintfFunc()
var infoFormat = color.New(color.FgHiWhite, color.Bold).PrintfFunc()
var subtextFormat = color.New(color.FgHiBlack).PrintlnFunc()
var highlightFormat = color.New(color.FgCyan, color.Bold).PrintFunc()

// print text with error format
func printError(err string) {
	errorFormat("[x] ")
	infoFormat(err + "\n")
}

// print option error
func (optErr *optionError) printOptionError() {
	errorFormat("[x] ")
	highlightFormat(optErr.option)
	infoFormat(" " + optErr.err.Error() + "\n")
}

// print options validation errors
func printOptionsErrors(err map[string]string) {
	fmt.Println()

	for key, value := range err {
		oe := &optionError{option: key, err: errors.New(value)}
		oe.printOptionError()
	}

	fmt.Println()
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
