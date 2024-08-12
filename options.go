package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
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

	flag.StringVar(&options.url, "u", "", "URL to check it´s CORS policy, it must start with http:// or https://")
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

	flag.Parse()

	return &options
}

// validate options format
func (o *options) validateOptions() *validator {
	validator := validator{}

	validator.check(notBlank(o.url) || notBlank(o.requests_list), "-u,-rl", "You must use one of this commands")

	validator.check(!notBlank(o.url) || maxChars(o.url, 100), "-u", "Cannot be longer than 100 characters")
	validator.check(!notBlank(o.url) || matches(o.url, urlRX), "-u", "Must have a URL format, must start with http:// or https://")

	validator.check(matches(o.method, methodRX), "-m", "Accepted methods GET, POST, PUT, DELETE and PATCH")

	validator.check(!notBlank(o.headers) || maxChars(o.headers, 500), "-h", "Cannot be longer than 500 characters")
	validator.check(!notBlank(o.headers) || matches(o.headers, headerRX), "-h", `Must follow the format "key:value, key:value, ..."`)

	validator.check(!notBlank(o.data) || maxChars(o.data, 500), "-d", "Cannot be longer than 500 characters")

	validator.check(!notBlank(o.origin) || maxChars(o.origin, 100), "-g", "Cannot be longer than 100 characters")
	validator.check(!notBlank(o.origin) || matches(o.origin, urlRX), "-g", "Must have a URL format, must start with http:// or https://")

	validator.check(!notBlank(o.origins_list) || maxChars(o.origins_list, 20), "-gl", "Cannot be longer than 20 characters")
	validator.check(!notBlank(o.origins_list) || matches(o.origins_list, fileRX), "-gl", "A filename cannot contain /")

	validator.check(!notBlank(o.requests_list) || maxChars(o.requests_list, 20), "-rl", "Cannot be longer than 20 characters")
	validator.check(!notBlank(o.requests_list) || matches(o.requests_list, fileRX), "-rl", "A filename cannot contain /")

	validator.check(!notBlank(o.output) || maxChars(o.output, 20), "-o", "Cannot be longer than 20 characters")
	validator.check(!notBlank(o.output) || matches(o.output, fileRX), "-o", "A filename cannot contain /")

	validator.check(minNumber(o.timeout, 0), "-t", "Must be greater that 0")
	validator.check(maxNumber(o.timeout, 10), "-t", "Must be lower that 10")

	validator.check(!notBlank(o.proxy) || matches(o.proxy, proxyRX), "-p", "Must start with http:// or socks5://")

	return &validator
}

// get origin headers from options
func (o *options) getOriginHeaders() ([]string, error) {
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
			return nil, errOpenFile("-gl", o.origins_list)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if url := scanner.Text(); matches(url, urlRX) {
				origins = append(origins, url)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, errReadFile("-gl", o.origins_list)
		}
	}

	return origins, nil
}

// get request from options
func (o *options) getRequests() (*[]request, error) {
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

		headersList := strings.Split(o.headers, ",")
		for _, header := range headersList {
			splitHeader := strings.Split(header, ":")
			headers[splitHeader[0]] = splitHeader[1]
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
			return nil, errOpenFile("-rl", o.requests_list)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var request request

			lineReader := bytes.NewReader(scanner.Bytes())
			err := readJSON(lineReader, &request)
			if err != nil {
				return nil, fmt.Errorf("-rl %s", err.Error())
			}

			requests = append(requests, request.addRequestsByOrigins(origins)...)
		}

		if err := scanner.Err(); err != nil {
			return nil, errReadFile("-rl", o.requests_list)
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

	return &requests, nil
}
