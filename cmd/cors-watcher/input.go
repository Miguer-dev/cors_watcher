package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Miguer-dev/cors_watcher/internal/validator"
)

type options struct {
	url         string
	method      string
	headers     string
	data        string
	originsFile struct {
		fileName        string
		origins         []string
		onlyOriginsFile bool
	}
	requestsFile struct {
		fileName string
		requests []request
	}
	output     string
	outputJSON string
	outputCSV  string
	outputYAML string
	timeout    int64
	timedelay  float64
	proxy      string
}

// init options intance with command options values
func initOptions() *options {
	options := &options{}

	flag.StringVar(&options.url, "url", "", "URL to check its CORS policy. It must start with http:// or https://.")
	flag.StringVar(&options.method, "method", "GET", "Set the request method (GET, POST, PUT, DELETE, PATCH).")
	flag.StringVar(&options.headers, "headers", "", `Set request headers in the format "key:value, key:value, ...".`)
	flag.StringVar(&options.data, "data", "", "Set request data.")
	flag.StringVar(&options.originsFile.fileName, "origins-file", "", "Specify the filename containing the list of origins.")
	flag.BoolVar(&options.originsFile.onlyOriginsFile, "only-origins", false, "Use only the origins from the specified origins list file.")
	flag.StringVar(&options.requestsFile.fileName, "requests-file", "", `Specify the filename containing the list of requests, using JSON format for each entry:
	{"url": "https://url1.com", "method": "POST", "headers": {"header1": "value1", "header2": "value2"}, "data": "data1"}`)
	flag.StringVar(&options.output, "output", "", "Specify the filename to save the results in a readable format.")
	flag.StringVar(&options.outputJSON, "output-json", "", "Specify the filename to save the results in json format.")
	flag.StringVar(&options.outputCSV, "output-csv", "", "Specify the filename to save the results in csv format.")
	flag.StringVar(&options.outputYAML, "output-yaml", "", "Specify the filename to save the results in yaml format.")
	flag.Int64Var(&options.timeout, "timeout", 10, "Set the request timeout (in seconds).")
	flag.Float64Var(&options.timedelay, "delay", 0, "Set the delay between requests (in seconds) (default 0).")
	flag.StringVar(&options.proxy, "proxy", "", "Set the proxy (HTTP or SOCKS5).")

	displayVersion := flag.Bool("version", false, "Display version")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	v := &validator.Validator{}

	options.validateOptions(v)

	if !v.Valid() {
		optsErrorPrintExit(v.Errors)
	}

	if options.outputJSON != "" && !strings.Contains(options.outputJSON, ".json") {
		options.outputJSON += ".json"
	}

	if options.outputCSV != "" && !strings.Contains(options.outputCSV, ".csv") {
		options.outputCSV += ".csv"
	}

	if options.originsFile.fileName != "" {
		options.getOriginsFromFile(v)

		if !v.Valid() {
			optsErrorPrintExit(v.Errors)
		}
	}

	if options.requestsFile.fileName != "" {
		options.getRequestsFromFile(v)

		if !v.Valid() {
			optsErrorPrintExit(v.Errors)
		}
	}

	return options
}

// validate options format
func (o *options) validateOptions(v *validator.Validator) {

	v.Check(validator.NotBlank(o.url) || validator.NotBlank(o.requestsFile.fileName), "-url,-requests-file", "You must use one of these options")

	v.Check(!validator.NotBlank(o.url) || validator.MaxChars(o.url, 100), "-url", "Cannot exceed 100 characters")
	v.Check(!validator.NotBlank(o.url) || validator.Matches(o.url, validator.URLRX), "-url", "Must be a valid URL, starting with http:// or https://")

	v.Check(validator.Matches(o.method, validator.MethodRX), "-method", "Accepted methods are GET, POST, PUT, DELETE, and PATCH")

	v.Check(!validator.NotBlank(o.headers) || validator.MaxChars(o.headers, 500), "-headers", "Cannot exceed 500 characters")
	v.Check(!validator.NotBlank(o.headers) || validator.Matches(o.headers, validator.HeaderRX), "-headers", `Must follow the format "key:value, key:value, ..."`)

	v.Check(!validator.NotBlank(o.data) || validator.MaxChars(o.data, 500), "-data", "Cannot exceed 500 characters")

	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.MaxChars(o.originsFile.fileName, 20), "-origins-file", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.Matches(o.originsFile.fileName, validator.FileRX), "-origins-file", "A filename cannot contain '/'")

	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.MaxChars(o.requestsFile.fileName, 20), "-requests-file", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.Matches(o.requestsFile.fileName, validator.FileRX), "-requests-file", "A filename cannot contain '/'")

	v.Check(!validator.NotBlank(o.output) || validator.MaxChars(o.output, 20), "-output", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.output) || validator.Matches(o.output, validator.FileRX), "-output", "A filename cannot contain '/'")

	v.Check(!validator.NotBlank(o.outputJSON) || validator.MaxChars(o.outputJSON, 20), "-output-json", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.outputJSON) || validator.Matches(o.outputJSON, validator.FileRX), "-output-json", "A filename cannot contain '/'")

	v.Check(!validator.NotBlank(o.outputCSV) || validator.MaxChars(o.outputCSV, 20), "-output-csv", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.outputCSV) || validator.Matches(o.outputCSV, validator.FileRX), "-output-csv", "A filename cannot contain '/'")

	v.Check(!validator.NotBlank(o.outputYAML) || validator.MaxChars(o.outputYAML, 20), "-output-yaml", "Cannot exceed 20 characters")
	v.Check(!validator.NotBlank(o.outputYAML) || validator.Matches(o.outputYAML, validator.FileRX), "-output-yaml", "A filename cannot contain '/'")

	v.Check(validator.MinNumber(o.timeout, 0), "-timeout", "Must be greater than 0")
	v.Check(validator.MaxNumber(o.timeout, 10), "-timeout", "Must be less than 10")

	v.Check(validator.MinNumber(o.timedelay, 0), "-delay", "Must be greater than 0")
	v.Check(validator.MaxNumber(o.timedelay, 5), "-delay", "Must be less than 5")

	v.Check(!validator.NotBlank(o.proxy) || validator.Matches(o.proxy, validator.ProxyRX), "-proxy", "Must start with http:// or socks5://")
}

// get and validate origins from originsFile -origins-file
func (o *options) getOriginsFromFile(v *validator.Validator) {
	file, err := os.Open(o.originsFile.fileName)
	if err != nil {
		v.AddError("-origins-file", errOpenFile(o.originsFile.fileName).Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		if !validator.NotBlank(url) {
			v.AddError("-origins-file", `There cannot be an empty row`)
			continue
		}

		if !validator.Matches(url, validator.URLRX) {
			v.AddError("-origins-file", `Each origin must be a valid URL, starting with http:// or https://`)
			continue
		}

		o.originsFile.origins = append(o.originsFile.origins, url)
	}

	if err := scanner.Err(); err != nil {
		v.AddError("-origins-file", errReadFile(o.originsFile.fileName).Error())
		return
	}
}

// get and validate requests from requestsFile -requests-file
func (o *options) getRequestsFromFile(v *validator.Validator) {
	file, err := os.Open(o.requestsFile.fileName)
	if err != nil {
		v.AddError("-requests-file", errOpenFile(o.requestsFile.fileName).Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var request request

		lineReader := bytes.NewReader(scanner.Bytes())

		err := readJSON(lineReader, &request)
		if err != nil {
			v.AddError("-requests-file", err.Error())
			continue
		}

		if !validator.NotBlank(request.URL) {
			v.AddError("-requests-file", `Each request must contain the 'url' key`)
			continue
		}

		if !validator.Matches(request.URL, validator.URLRX) {
			v.AddError("-requests-file", `The 'url' key must be a valid URL, starting with http:// or https://`)
			continue
		}

		if !validator.NotBlank(request.Method) {
			v.AddError("-requests-file", `Each request must contain the 'method' key`)
			continue
		}

		if !validator.Matches(request.Method, validator.MethodRX) {
			v.AddError("-requests-file", `The 'method' key must be one of: GET, POST, PUT, DELETE, PATCH`)
			continue
		}

		o.requestsFile.requests = append(o.requestsFile.requests, request)
	}

	if err := scanner.Err(); err != nil {
		v.AddError("-requests-file", errReadFile(o.requestsFile.fileName).Error())
		return
	}
}
