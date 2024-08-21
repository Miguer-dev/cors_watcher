package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type transaction struct {
	request    request
	response   response
	tags       []tag
	err        error
	targetHost string
}

type request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    string            `json:"data,omitempty"`
}

type response struct {
	statusCode int
	length     int64
	ACDetected bool   // Access-Control-* headers detected
	ACAO       string // Access-Control-Allow-Origin header
	ACAC       string // Access-Control-Allow-Credentials header
	vary       bool   // Vary: Origin header detected
}

type tag struct {
	info  string
	print func(a ...interface{})
}

// build transactions with all options
func initTransactions(o *options) [][]*transaction {
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

	var resultTransactions [][]*transaction

	for _, transaction := range transactions {
		resultTransactions = append(resultTransactions, transaction.addtransactionsByOrigins(o))
	}

	return resultTransactions
}

// foreach origin duplicate request
func (t transaction) addtransactionsByOrigins(o *options) []*transaction {
	var transactions []*transaction

	origins := setOrigins(t.request.URL, o)

	t.targetHost = origins[2]

	for _, origin := range origins {
		copyTransaction := t

		copyTransaction.request.Headers = make(map[string]string)
		for key, value := range t.request.Headers {
			copyTransaction.request.Headers[key] = value
		}

		copyTransaction.request.Headers["Origin"] = origin

		transactions = append(transactions, &copyTransaction)
	}

	return transactions
}

// set origins for URL
func setOrigins(u string, o *options) []string {
	originList := []string{}

	if !o.originsFile.onlyOriginsFile {
		originList = append(originList, "https://test.com", "null")

		host, err := url.Parse(u)
		if err == nil {

			hostOriginOption := []string{
				host.Scheme + "://" + host.Host,
				host.Scheme + "://" + "test" + host.Host,
				addPortIfExist(host.Scheme+"://"+host.Hostname()+".test.com", host),
				host.Scheme + "://" + "test." + host.Host,
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"!.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+`".test.com`, host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"$.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"%0b.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"%60.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"_.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"&.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"'.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"(.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+").test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"*.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+",.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+";.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"=.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"^.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"`.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"{.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"|.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"}.test.com", host),
				addPortIfExist(host.Scheme+"://"+"test."+host.Hostname()+"~.test.com", host),
			}

			originList = append(originList, hostOriginOption...)
		}
	}

	if len(o.originsFile.origins) != 0 {
		for _, origin := range o.originsFile.origins {
			originList = append(originList, origin)
		}
	}

	return originList
}

// Send Request
func (t *transaction) sendRequest(client *http.Client) {
	request, err := http.NewRequest(t.request.Method, t.request.URL, bytes.NewBuffer([]byte(t.request.Data)))
	if err != nil {
		t.err = err
		return
	}

	for key, value := range t.request.Headers {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		t.err = err
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.response.length = -1
	} else {
		t.response.length = int64(len(body))
	}

	t.response.statusCode = response.StatusCode

	for key, value := range response.Header {
		if strings.Contains(key, "Access-Control-") {
			t.response.ACDetected = true

			if key == "Access-Control-Allow-Origin" {
				t.response.ACAO = value[0]
			} else if key == "Access-Control-Allow-Credentials" {
				t.response.ACAC = value[0]
			}
		}

		if key == "Vary" {
			for _, vary := range value {
				if vary == "Origin" {
					t.response.vary = true
					break
				}
			}
		}

	}
}

// create tags from response
func (t *transaction) addTags() {
	if t.err != nil {
		t.tags = append(t.tags, tag{info: " Transaction Fail ", print: redBackgroundFormat})
		return
	}

	if t.response.ACDetected {
		t.tags = append(t.tags, tag{info: " AC* ", print: cyanBackgroundFormat})

		switch t.response.ACAO {
		case "":
			break
		case "*":
			t.tags = append(t.tags, tag{info: fmt.Sprintf(" ACAO:%s ", t.response.ACAO), print: greenBackgroundFormat})
		case t.targetHost:
			t.tags = append(t.tags, tag{info: fmt.Sprintf(" ACAO:%s ", t.response.ACAO), print: greenBackgroundFormat})
		case t.request.Headers["Origin"]:
			t.tags = append(t.tags, tag{info: fmt.Sprintf(" ACAO:%s ", t.response.ACAO), print: yellowBackgroundFormat})
		default:
			t.tags = append(t.tags, tag{info: fmt.Sprintf(" ACAO:%s ", t.response.ACAO), print: greenBackgroundFormat})
		}

		switch t.response.ACAC {
		case "true":
			if t.response.ACAO == "*" || t.response.ACAO == "" || t.response.ACAO != t.request.Headers["Origin"] || t.response.ACAO == t.targetHost {
				t.tags = append(t.tags, tag{info: " ACAC:true ", print: greenBackgroundFormat})
			} else {
				t.tags = append(t.tags, tag{info: " ACAC:true ", print: redBackgroundFormat})
			}
		case "false":
			t.tags = append(t.tags, tag{info: " ACAC:false ", print: greenBackgroundFormat})
		}

		if (t.response.ACAO == "*" || t.response.ACAO == t.request.Headers["Origin"]) && t.response.ACAO != t.targetHost {
			if strings.Contains(t.response.ACAO, "http://") {
				t.tags = append(t.tags, tag{info: " HTTP ", print: yellowBackgroundFormat})
			}

			if !t.response.vary {
				t.tags = append(t.tags, tag{info: " Not Vary:Origin ", print: yellowBackgroundFormat})
			}

		}
	}
}
