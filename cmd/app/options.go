package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"cors_watcher/internal/validator"
)

type options struct {
	url         string
	method      string
	headers     string
	data        string
	origin      string
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
	options := options{}

	flag.StringVar(&options.url, "u", "", "URL to Check it´s CORS policy, it must start with http:// or https://")
	flag.StringVar(&options.method, "m", "GET", "Set request method (GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)")
	flag.StringVar(&options.headers, "e", "", `Set request headers, format "key:value, key:value, ..."`)
	flag.StringVar(&options.data, "d", "", "Set request data")
	flag.StringVar(&options.origin, "g", "", "Set origin header, it must start with http:// or https://")
	flag.StringVar(&options.originsFile.fileName, "gl", "", "Set filename containing the origins list")
	flag.StringVar(&options.requestsFile.fileName, "rl", "", `Set filename containing the requests list, use json format for each row
	{"url": "https://url1.com", "method": "POST", "headers": {"header1": "value1", "header2": "value2"}, "data": "data1"}`)
	flag.StringVar(&options.output, "o", "", "Set filename to save the result")
	flag.IntVar(&options.timeout, "t", 0, "Set requests timeout")
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

	return &options
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

	v.Check(!validator.NotBlank(o.origin) || validator.MaxChars(o.origin, 100), "-g", "Cannot be longer than 100 characters")
	v.Check(!validator.NotBlank(o.origin) || validator.Matches(o.origin, validator.URLRX), "-g", "Must have a URL format, must start with http:// or https://")

	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.MaxChars(o.originsFile.fileName, 20), "-gl", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.originsFile.fileName) || validator.Matches(o.originsFile.fileName, validator.FileRX), "-gl", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.MaxChars(o.requestsFile.fileName, 20), "-rl", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.requestsFile.fileName) || validator.Matches(o.requestsFile.fileName, validator.FileRX), "-rl", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.output) || validator.MaxChars(o.output, 20), "-o", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.output) || validator.Matches(o.output, validator.FileRX), "-o", "A filename cannot contain /")

	v.Check(validator.MinNumber(o.timeout, 0), "-t", "Must be greater that 0")
	v.Check(validator.MaxNumber(o.timeout, 100), "-t", "Must be lower that 100")

	v.Check(!validator.NotBlank(o.proxy) || validator.Matches(o.proxy, validator.ProxyRX), "-p", "Must start with http:// or socks5://")
}

// get and validate origins from originsFile -gl
func (o *options) getOriginsFromFile(v *validator.Validator) {
	file, err := os.Open(o.originsFile.fileName)
	if err != nil {
		v.AddError("-gl", errOpenFile(o.originsFile.fileName).Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		if !validator.NotBlank(url) {
			v.AddError("-gl", `There cannot be an empty row`)
			continue
		}

		if !validator.Matches(url, validator.URLRX) {
			v.AddError("-gl", `Origins must have a URL format, must start with http:// or https://"`)
			continue
		}

		o.originsFile.origins = append(o.originsFile.origins, url)
	}

	if err := scanner.Err(); err != nil {
		v.AddError("-gl", errReadFile(o.originsFile.fileName).Error())
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

// get origin headers from all options
func (o *options) getAllOriginHeaders() []string {
	var origins = []string{"https://test.com", "null"}

	if o.origin != "" {
		origins = append(origins, o.origin)
	}

	if o.url != "" {
		origins = append(origins, o.url)
	}

	if len(o.originsFile.origins) != 0 {
		origins = append(origins, o.originsFile.origins...)
	}

	return origins
}

// build requests with all options
func (o *options) buildRequests() []request {
	var requests []request

	origins := o.getAllOriginHeaders()

	for _, value := range origins {
		fmt.Println(value)
	}

	if o.url != "" {
		var headers = make(map[string]string)

		if o.headers != "" {
			headersList := strings.Split(o.headers, ",")
			for _, header := range headersList {
				splitHeader := strings.Split(header, ":")
				headers[splitHeader[0]] = splitHeader[1]
			}
		}

		request := request{
			URL:     o.url,
			Method:  o.method,
			Headers: headers,
			Data:    o.data,
		}

		requests = append(requests, request.addRequestsByOrigins(origins)...)
	}

	if len(o.requestsFile.requests) != 0 {
		for _, request := range o.requestsFile.requests {
			requests = append(requests, request.addRequestsByOrigins(origins)...)
		}
	}

	for _, value := range requests {
		fmt.Println()
		fmt.Println(value.URL)
		fmt.Println(value.Method)
		fmt.Println(value.Headers)
		fmt.Println(value.Data)
	}

	return requests
}
