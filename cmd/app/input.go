package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"

	"cors_watcher/internal/validator"
)

type options struct {
	url         string
	method      string
	headers     string
	data        string
	originsFile struct {
		fileName string
		origins  []string
	}
	requestsFile struct {
		fileName string
		requests []request
	}
	output  string
	timeout int
	proxy   string
}

// init options intance with command options values
func initOptions() *options {
	options := &options{}

	flag.StringVar(&options.url, "u", "", "URL to Check it´s CORS policy, it must start with http:// or https://")
	flag.StringVar(&options.method, "m", "GET", "Set request method (GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)")
	flag.StringVar(&options.headers, "e", "", `Set request headers, format "key:value, key:value, ..."`)
	flag.StringVar(&options.data, "d", "", "Set request data")
	flag.StringVar(&options.originsFile.fileName, "ol", "", "Set filename containing the origins list")
	flag.StringVar(&options.requestsFile.fileName, "rl", "", `Set filename containing the requests list, use json format for each row
	{"url": "https://url1.com", "method": "POST", "headers": {"header1": "value1", "header2": "value2"}, "data": "data1"}`)
	flag.StringVar(&options.output, "o", "", "Set filename to save the result")
	flag.IntVar(&options.timeout, "t", 10, "Set requests timeout, default 10 seconds")
	flag.StringVar(&options.proxy, "p", "", "Set proxy (http or socks5)")

	displayVersion := flag.Bool("v", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	v := &validator.Validator{}

	options.validateOptions(v)

	if !v.Valid() {
		optsErrorPrintExit(v.Errors)
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

	v.Check(validator.NotBlank(o.url) || validator.NotBlank(o.requestsFile.fileName), "-u,-rl", "You must use one of this commands")

	v.Check(!validator.NotBlank(o.url) || validator.MaxChars(o.url, 100), "-u", "Cannot be longer than 100 characters")
	v.Check(!validator.NotBlank(o.url) || validator.Matches(o.url, validator.URLRX), "-u", "Must have a URL format, must start with http:// or https://")

	v.Check(validator.Matches(o.method, validator.MethodRX), "-m", "Accepted methods GET, POST, PUT, DELETE and PATCH")

	v.Check(!validator.NotBlank(o.headers) || validator.MaxChars(o.headers, 500), "-e", "Cannot be longer than 500 characters")
	v.Check(!validator.NotBlank(o.headers) || validator.Matches(o.headers, validator.HeaderRX), "-e", `Must follow the format "key:value, key:value, ..."`)

	v.Check(!validator.NotBlank(o.data) || validator.MaxChars(o.data, 500), "-d", "Cannot be longer than 500 characters")

	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.MaxChars(o.originsFile.fileName, 20), "-ol", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.Matches(o.originsFile.fileName, validator.FileRX), "-ol", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.MaxChars(o.requestsFile.fileName, 20), "-rl", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.Matches(o.requestsFile.fileName, validator.FileRX), "-rl", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.output) || validator.MaxChars(o.output, 20), "-o", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.output) || validator.Matches(o.output, validator.FileRX), "-o", "A filename cannot contain /")

	v.Check(validator.MinNumber(o.timeout, 0), "-t", "Must be greater that 0")
	v.Check(validator.MaxNumber(o.timeout, 10), "-t", "Must be lower that 10")

	v.Check(!validator.NotBlank(o.proxy) || validator.Matches(o.proxy, validator.ProxyRX), "-p", "Must start with http:// or socks5://")
}

// get and validate origins from originsFile -ol
func (o *options) getOriginsFromFile(v *validator.Validator) {
	file, err := os.Open(o.originsFile.fileName)
	if err != nil {
		v.AddError("-ol", errOpenFile(o.originsFile.fileName).Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		if !validator.NotBlank(url) {
			v.AddError("-ol", `There cannot be an empty row`)
			continue
		}

		if !validator.Matches(url, validator.URLRX) {
			v.AddError("-ol", `Origins must have a URL format, must start with http:// or https://`)
			continue
		}

		o.originsFile.origins = append(o.originsFile.origins, url)
	}

	if err := scanner.Err(); err != nil {
		v.AddError("-ol", errReadFile(o.originsFile.fileName).Error())
		return
	}
}

// get and validate requests from requestsFile -rl
func (o *options) getRequestsFromFile(v *validator.Validator) {
	file, err := os.Open(o.requestsFile.fileName)
	if err != nil {
		v.AddError("-rl", errOpenFile(file.Name()).Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var request request

		lineReader := bytes.NewReader(scanner.Bytes())

		err := readJSON(lineReader, &request)
		if err != nil {
			v.AddError("-rl", err.Error())
			continue
		}

		if !validator.NotBlank(request.URL) {
			v.AddError("-rl", `Each request must contain the “url” key`)
			continue
		}

		if !validator.Matches(request.URL, validator.URLRX) {
			v.AddError("-rl", `The “url” key must have a URL format, must start with http:// or https://`)
			continue
		}

		if !validator.NotBlank(request.Method) {
			v.AddError("-rl", `Each request must contain the “method” key`)
			continue
		}

		if !validator.Matches(request.Method, validator.MethodRX) {
			v.AddError("-rl", `The “method” key accepted values: GET, POST, PUT, DELETE and PATCH`)
			continue
		}

		o.requestsFile.requests = append(o.requestsFile.requests, request)
	}

	if err := scanner.Err(); err != nil {
		v.AddError("-rl", errReadFile(o.requestsFile.fileName).Error())
		return
	}
}
