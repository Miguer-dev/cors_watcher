package main

import (
	"flag"
)

type Options struct {
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
func initOptions() *Options {
	options := Options{}

	flag.StringVar(&options.url, "u", "", "URL to check it´s CORS policy, it must start with http:// or https://")
	flag.StringVar(&options.method, "m", "GET", "Set request method (GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)")
	flag.StringVar(&options.headers, "e", "", `Set request headers, format "key: value, key:value, ..."`)
	flag.StringVar(&options.data, "d", "", "Set request data")
	flag.StringVar(&options.origin, "g", "", "Set origin header, it must start with http:// or https://")
	flag.StringVar(&options.origins_list, "gl", "", "Set filename containing the origins list")
	flag.StringVar(&options.requests_list, "rl", "", "Set filename containing the requests list")
	flag.StringVar(&options.output, "o", "", "Set filename to save the result")
	flag.IntVar(&options.timeout, "t", 0, "Set requests timeout")
	flag.StringVar(&options.proxy, "p", "", "Set proxy (http or socks5)")

	flag.Parse()

	return &options
}

// validate options format
func (o *Options) validateOptions() *Validator {
	validator := Validator{}

	validator.Check(NotBlank(o.url), "-u", "You must use the -u command to set the target url")
	validator.Check(MaxChars(o.url, 100), "-u", "Cannot be longer than 100 characters")
	validator.Check(Matches(o.url, URLRX), "-u", "Must have a URL format, must start with http:// or https://")

	validator.Check(Matches(o.method, MethodRX), "-m", "Accepted methods GET, POST, PUT, DELETE, HEAD, OPTIONS and PATCH")

	validator.Check(!NotBlank(o.headers) || MaxChars(o.headers, 500), "-h", "Cannot be longer than 500 characters")
	validator.Check(!NotBlank(o.headers) || Matches(o.headers, HeaderRX), "-h", `Must follow the format "key: value, key:value, ..."`)

	validator.Check(!NotBlank(o.data) || MaxChars(o.data, 500), "-d", "Cannot be longer than 500 characters")

	validator.Check(!NotBlank(o.origin) || MaxChars(o.origin, 100), "-g", "Cannot be longer than 100 characters")
	validator.Check(!NotBlank(o.origin) || Matches(o.origin, URLRX), "-g", "Must have a URL format, must start with http:// or https://")

	validator.Check(!NotBlank(o.origins_list) || MaxChars(o.origins_list, 20), "-gl", "Cannot be longer than 20 characters")
	validator.Check(!NotBlank(o.origins_list) || Matches(o.origins_list, FileRX), "-gl", "A filename cannot contain /")

	validator.Check(!NotBlank(o.requests_list) || MaxChars(o.requests_list, 20), "-rl", "Cannot be longer than 20 characters")
	validator.Check(!NotBlank(o.requests_list) || Matches(o.requests_list, FileRX), "-rl", "A filename cannot contain /")

	validator.Check(!NotBlank(o.output) || MaxChars(o.output, 20), "-o", "Cannot be longer than 20 characters")
	validator.Check(!NotBlank(o.output) || Matches(o.output, FileRX), "-o", "A filename cannot contain /")

	validator.Check(MinNumber(o.timeout, 0), "-t", "Must be greater that 0")
	validator.Check(MaxNumber(o.timeout, 10), "-t", "Must be lower that 10")

	validator.Check(!NotBlank(o.proxy) || Matches(o.proxy, ProxyRX), "-p", "Must start with http:// or socks5://")

	return &validator
}
