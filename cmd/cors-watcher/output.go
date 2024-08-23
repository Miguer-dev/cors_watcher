package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	statusSpaces := spaces(int64(status), 5)
	len := transaction.response.length
	lenSpaces := spaces(len, 5)

	lenString := strconv.FormatInt(len, 10)
	if lenString == "-1" {
		lenString = "unk"
		lenSpaces = "  "
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	fmt.Printf("| %d%s| %s%s| %s", status, statusSpaces, lenString, lenSpaces, transaction.request.Headers["Origin"])

	for _, tag := range transaction.tags {
		fmt.Print(" ")
		tag.print(tag.Info)
	}

	fmt.Println()
}

// save output in file
func printFile(filename string, transactions [][]*transaction) {
	if filename != "" {

		file, err := os.Create(filename)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output", Err: errCreateFile(filename).Error()})
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
				statusSpaces := spaces(int64(status), 5)
				len := transaction.response.length
				lenSpaces := spaces(len, 5)

				lenString := strconv.FormatInt(len, 10)
				if lenString == "-1" {
					lenString = "unk"
					lenSpaces = "  "
				}

				text += fmt.Sprintf("| %d%s| %s%s| %s", status, statusSpaces, lenString, lenSpaces, transaction.request.Headers["Origin"])

				for _, tag := range transaction.tags {
					text += fmt.Sprint(" ")
					text += tag.Info
				}

				text += "\n"
			}
			text += "\n"
		}

		_, err = file.WriteString(text)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output", Err: errWriteFile(filename).Error()})
		}

		fmt.Println()
		printWarning("Saving output ...")
		printInfo(fmt.Sprintf(`Output successfully saved in “%s”`, filename))
	}
}

type transactionsOutput struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	Data      string            `json:"data,omitempty"`
	Responses []responsesOutput `json:"responses"`
}

type responsesOutput struct {
	StatusCode int    `json:"status_code"`
	Length     int64  `json:"size"`
	Origin     string `json:"origin"`
	Tags       []tag  `json:"tags,omitempty"`
}

// save output in json format
func printJsonFile(filename string, transactions [][]*transaction) {
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output-json", Err: errCreateFile(filename).Error()})
		}
		defer file.Close()

		arrayTransactionsOutput := []transactionsOutput{}
		for _, arrayTransactions := range transactions {
			transactionOutput := transactionsOutput{
				URL:     arrayTransactions[0].request.URL,
				Method:  arrayTransactions[0].request.Method,
				Data:    arrayTransactions[0].request.Data,
				Headers: map[string]string{}}

			if len(arrayTransactions[0].request.Headers) > 1 {

				for key, value := range arrayTransactions[0].request.Headers {
					if key != "Origin" {
						transactionOutput.Headers[key] = value
					}
				}
			}

			for _, transaction := range arrayTransactions {
				response := responsesOutput{
					StatusCode: transaction.response.statusCode,
					Length:     transaction.response.length,
					Origin:     transaction.request.Headers["Origin"],
					Tags:       transaction.tags,
				}

				transactionOutput.Responses = append(transactionOutput.Responses, response)
			}

			arrayTransactionsOutput = append(arrayTransactionsOutput, transactionOutput)
		}

		json, err := writeJSON(map[string]any{"requests": arrayTransactionsOutput})
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output-json", Err: errWriteFile(filename).Error()})
		}
		_, err = file.Write(json)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output-json", Err: errWriteFile(filename).Error()})
		}

		fmt.Println()
		printWarning("Saving output in JSON format ...")
		printInfo(fmt.Sprintf(`Output successfully saved in “%s”`, filename))
	}
}

// save output in csv format
func printCsvFile(filename string, transactions [][]*transaction) {

	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output-csv", Err: errCreateFile(filename).Error()})
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Url", "Method", "Headers", "Data", "Status", "Size", "Origin", "Tags"}
		if err := writer.Write(header); err != nil {
			optErrorPrintExit(&validator.OptionError{Option: "-output-csv", Err: errWriteFile(filename).Error()})
		}

		data := [][]string{}

		for _, arrayTransactions := range transactions {
			for _, transaction := range arrayTransactions {
				headers := map[string]string{}

				if len(transaction.request.Headers) > 1 {

					for key, value := range transaction.request.Headers {
						if key != "Origin" {
							headers[key] = value
						}
					}
				}

				headersJSON, err := json.Marshal(headers)
				if err != nil {
					optErrorPrintExit(&validator.OptionError{Option: "-output-csv", Err: errWriteFile(filename).Error()})
				}

				tags := []string{}
				for _, tag := range transaction.tags {
					tags = append(tags, tag.Info)
				}
				tagsString := strings.Join(tags, ",")

				origin := transaction.request.Headers["Origin"]
				if strings.Contains(origin, ";") {
					origin = `"` + origin + `"`
				}

				row := []string{
					transaction.request.URL,
					transaction.request.Method,
					string(headersJSON),
					transaction.request.Data,
					strconv.Itoa(transaction.response.statusCode),
					strconv.FormatInt(transaction.response.length, 10),
					origin,
					tagsString,
				}

				data = append(data, row)
			}
		}

		for _, row := range data {
			if err := writer.Write(row); err != nil {
				optErrorPrintExit(&validator.OptionError{Option: "-output-csv", Err: errWriteFile(filename).Error()})
			}
		}

		fmt.Println()
		printWarning("Saving output in CSV format ...")
		printInfo(fmt.Sprintf(`Output successfully saved in “%s”`, filename))
	}
}
