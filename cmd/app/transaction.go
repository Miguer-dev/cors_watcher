package main

import (
	"strings"
)

type transaction struct {
	name     string
	request  request
	response response
	tags     []tag
}

type request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    string            `json:"data,omitempty"`
}

type response struct {
	statusCode int
	length     int
	ACDetected bool   // Access-Control-* headers detected
	ACAO       string // Access-Control-Allow-Origin value
	ACAC       string // Access-Control-Allow-Credentials value
}

type tag struct {
	exploit string
	threat  string
}

var (
	xsrf = tag{
		exploit: "Origin impersonation + XSRF",
		threat:  "hight",
	}

	browser = tag{
		exploit: "Browser dependent",
		threat:  "medium",
	}

	xss = tag{
		exploit: "Trusted Subdomains + XSS",
		threat:  "medium",
	}

	http = tag{
		exploit: "Http + ManInTheMiddle",
		threat:  "low",
	}

	cache = tag{
		exploit: "Cache Poisoning",
		threat:  "low",
	}
)

// build transactions with all options
func initTransactions(o *options) []*transaction {
	var transactions []*transaction

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

		transactions = append(transactions, &transaction{request: request})
	}

	if len(o.requestsFile.requests) != 0 {
		for _, request := range o.requestsFile.requests {
			transactions = append(transactions, &transaction{request: request})
		}
	}

	var resultTransactions []*transaction

	for _, transaction := range transactions {
		resultTransactions = append(resultTransactions, transaction.addtransactionsByOrigins(o)...)
	}

	return resultTransactions
}

// foreach origin duplicate request
func (t transaction) addtransactionsByOrigins(o *options) []*transaction {
	var transactions []*transaction

	origins := setOrigins(t.request.URL, o)

	for _, origin := range origins {
		copyTransaction := t

		copyTransaction.request.Headers = make(map[string]string)
		for key, value := range t.request.Headers {
			copyTransaction.request.Headers[key] = value
		}

		copyTransaction.request.Headers["Origin"] = origin.origin
		copyTransaction.name = origin.name
		copyTransaction.tags = origin.tags

		transactions = append(transactions, &copyTransaction)
	}

	return transactions
}

type originSearch struct {
	name   string
	origin string
	tags   []tag
}

// set origins fota URL
func setOrigins(url string, o *options) []*originSearch {
	originDefaults := []*originSearch{
		{
			name:   "Origin reflected",
			origin: "https://test.com",
			tags:   []tag{xsrf},
		},
		{
			name:   "Origin null",
			origin: "null",
			tags:   []tag{xsrf},
		},
	}

	if url != "" {
		splitOrigin := splitURL(url)

		if splitOrigin[0] != "" {
			hostOriginOption := []*originSearch{
				{
					name:   "Origin domain",
					origin: o.url,
				},
				{
					name:   "Origin suffix",
					origin: splitOrigin[0] + "test" + splitOrigin[1],
					tags:   []tag{xsrf},
				},
				{
					name:   "Origin prefix",
					origin: o.url + ".test.com",
					tags:   []tag{xsrf},
				},
				{
					name:   "Origin subdomain",
					origin: splitOrigin[0] + "test." + splitOrigin[1],
					tags:   []tag{xss},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "!.test.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + `".evil.com`,
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "$.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "%0b.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "%60.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "&.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "'.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "(.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + ").evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "*.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + ",.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + ";.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "=.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "^.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "`.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "{.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "|.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "}.evil.com",
					tags:   []tag{browser, xsrf},
				},
				{
					name:   "Origin subdomain special characters",
					origin: splitOrigin[0] + "test." + splitOrigin[1] + "~.evil.com",
					tags:   []tag{browser, xsrf},
				},
			}

			originDefaults = append(originDefaults, hostOriginOption...)
		}

	}

	if len(o.originsFile.origins) != 0 {
		for _, origin := range o.originsFile.origins {
			originDefaults = append(originDefaults, &originSearch{name: "Origin List", origin: origin})
		}

	}

	return originDefaults
}
