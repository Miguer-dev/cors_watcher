package main

import (
	"fmt"

	"github.com/fatih/color"
)

var errorFormat = color.New(color.FgRed, color.Bold).PrintfFunc()
var warningFormat = color.New(color.FgYellow, color.Bold).PrintfFunc()
var successFormat = color.New(color.FgGreen, color.Bold).PrintfFunc()
var infoFormat = color.New(color.FgHiWhite, color.Bold).PrintfFunc()
var subtextFormat = color.New(color.FgHiBlack).PrintlnFunc()

// print text with error format
func printError(err string) {
	errorFormat("[x] ")
	infoFormat(err + "\n")
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

// print options validation errors
func printOptionsErrors(err map[string]string) {
	fmt.Println()

	for key, value := range err {
		printError(key + " " + value)
	}

	fmt.Println()
}

// print app title
func title() {
	fmt.Println()
	fmt.Println("█▀▀ █▀█ █▀█ █▀▀  █░█░█ ░▄▄ █▄ ▄▄▄ █▄▄ ▄▄░ ▄▄")
	fmt.Println("█▄▄ █▄█ █▀▄ ▄▄█  █▄█▄█ ▀▄█ █▄ █▄▄ █░█ ██▄ █░")
	fmt.Println("░░░ ░░░ ░░░ ░░░  ░░░░░ ░░░ ░░ ░░░ ░░░ ░░░ ░░")
	subtextFormat("                               by Miguer-dev")
}
