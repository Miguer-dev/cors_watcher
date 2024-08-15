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
	url           string
	method        string
	headers       string
	data          string
	origin        string
	origins_list  string
	requests_list string
	output        string
	timeout       int
	proxy         string
}

// init options intance with command options values
func initOptions() *options {
	options := options{}

	flag.StringVar(&options.url, "u", "", "URL to Check it´s CORS policy, it must start with http:// or https://")
	flag.StringVar(&options.method, "m", "GET", "Set request method (GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)")
	flag.StringVar(&options.headers, "e", "", `Set request headers, format "key:value, key:value, ..."`)
	flag.StringVar(&options.data, "d", "", "Set request data")
	flag.StringVar(&options.origin, "g", "", "Set origin header, it must start with http:// or https://")
	flag.StringVar(&options.origins_list, "gl", "", "Set filename containing the origins list")
	flag.StringVar(&options.requests_list, "rl", "", `Set filename containing the requests list, use json format for each row
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

	return &options
}

// validate options format
func (o *options) validateOptions() *validator.Validator {
	v := validator.Validator{}

	v.Check(validator.NotBlank(o.url) || validator.NotBlank(o.requests_list), "-u,-rl", "You must use one of this commands")

	v.Check(!validator.NotBlank(o.url) || validator.MaxChars(o.url, 100), "-u", "Cannot be longer than 100 characters")
	v.Check(!validator.NotBlank(o.url) || validator.Matches(o.url, validator.URLRX), "-u", "Must have a URL format, must start with http:// or https://")

	v.Check(validator.Matches(o.method, validator.MethodRX), "-m", "Accepted methods GET, POST, PUT, DELETE and PATCH")

	v.Check(!validator.NotBlank(o.headers) || validator.MaxChars(o.headers, 500), "-e", "Cannot be longer than 500 characters")
	v.Check(!validator.NotBlank(o.headers) || validator.Matches(o.headers, validator.HeaderRX), "-e", `Must follow the format "key:value, key:value, ..."`)

	v.Check(!validator.NotBlank(o.data) || validator.MaxChars(o.data, 500), "-d", "Cannot be longer than 500 characters")

	v.Check(!validator.NotBlank(o.origin) || validator.MaxChars(o.origin, 100), "-g", "Cannot be longer than 100 characters")
	v.Check(!validator.NotBlank(o.origin) || validator.Matches(o.origin, validator.URLRX), "-g", "Must have a URL format, must start with http:// or https://")

	v.Check(!validator.NotBlank(o.origins_list) || validator.MaxChars(o.origins_list, 20), "-gl", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.origins_list) || validator.Matches(o.origins_list, validator.FileRX), "-gl", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.requests_list) || validator.MaxChars(o.requests_list, 20), "-rl", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.requests_list) || validator.Matches(o.requests_list, validator.FileRX), "-rl", "A filename cannot contain /")

	v.Check(!validator.NotBlank(o.output) || validator.MaxChars(o.output, 20), "-o", "Cannot be longer than 20 characters")
	v.Check(!validator.NotBlank(o.output) || validator.Matches(o.output, validator.FileRX), "-o", "A filename cannot contain /")

	v.Check(validator.MinNumber(o.timeout, 0), "-t", "Must be greater that 0")
	v.Check(validator.MaxNumber(o.timeout, 100), "-t", "Must be lower that 100")

	v.Check(!validator.NotBlank(o.proxy) || validator.Matches(o.proxy, validator.ProxyRX), "-p", "Must start with http:// or socks5://")

	return &v
}

// validate request format from requestFile -rl
func validateRequestList(url string, method string) *validator.Validator {
	v := validator.Validator{}

	v.Check(validator.NotBlank(url), "-rl", `Must contain key "url"`)
	v.Check(validator.Matches(url, validator.URLRX), "-rl", `Key "url" must have a URL format, must start with http:// or https://"`)

	v.Check(validator.NotBlank(method), "-rl", `Must contain key "method"`)
	v.Check(validator.Matches(method, validator.MethodRX), "-rl", `Key "method" accepted methods GET, POST, PUT, DELETE and PATCH`)

	return &v
}

// validate origins url format from originFile -gl
func validateOriginList(origin string) *validator.Validator {
	v := validator.Validator{}

	v.Check(validator.NotBlank(origin), "-gl", `origin file can´t be empty`)
	v.Check(validator.Matches(origin, validator.URLRX), "-gl", `origins must have a URL format, must start with http:// or https://"`)

	return &v
}

// get origin headers from options
func (o *options) getOriginHeaders() ([]string, *optionError) {
	var origins = []string{"https://test.com", "null"}

	if o.origin != "" {
		origins = append(origins, o.origin)
	}

	if o.url != "" {
		origins = append(origins, o.url)
	}

	if o.origins_list != "" {
		file, err := os.Open(o.origins_list)
		if err != nil {
			return nil, &optionError{option: "-gl", err: errOpenFile(o.origins_list)}
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()

			if originListValidations := validateOriginList(url); !originListValidations.Valid() {
				optsErrorPrintExit(originListValidations.Errors)
			}

			origins = append(origins, url)
		}

		if err := scanner.Err(); err != nil {
			return nil, &optionError{option: "-gl", err: errReadFile(o.origins_list)}
		}
	}

	return origins, nil
}

// get request from options
func (o *options) getRequests() ([]request, *optionError) {
	var requests []request

	origins, err := o.getOriginHeaders()
	if err != nil {
		return nil, err
	}

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

		originRequests := request.addRequestsByOrigins(origins)

		requests = append(requests, originRequests...)

	}

	if o.requests_list != "" {
		file, err := os.Open(o.requests_list)
		if err != nil {
			return nil, &optionError{option: "-rl", err: errOpenFile(o.requests_list)}
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var request request

			lineReader := bytes.NewReader(scanner.Bytes())
			err := readJSON(lineReader, &request)
			if err != nil {
				return nil, &optionError{option: "-rl", err: err}
			}

			if requestListValidations := validateRequestList(request.URL, request.Method); !requestListValidations.Valid() {
				optsErrorPrintExit(requestListValidations.Errors)
			}

			requests = append(requests, request.addRequestsByOrigins(origins)...)
		}

		if err := scanner.Err(); err != nil {
			return nil, &optionError{option: "-rl", err: errReadFile(o.requests_list)}
		}

	}

	for _, value := range requests {
		fmt.Println()
		fmt.Println(value.URL)
		fmt.Println(value.Method)
		fmt.Println(value.Headers)
		fmt.Println(value.Data)
		fmt.Println()
	}

	return requests, nil
}
