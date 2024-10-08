package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/Miguer-dev/cors_watcher/internal/assert"
	"github.com/Miguer-dev/cors_watcher/internal/validator"
)

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name           string
		options        options
		expectedErrors []*validator.OptionError
	}{
		{
			name:    "no options",
			options: options{method: "GET"}, // method is default
			expectedErrors: []*validator.OptionError{
				{Option: "-url,-requests-file",
					Err: "You must use one of these options",
				},
			},
		},
		{
			name: "wrong url format",
			options: options{method: "GET",
				url: "test.com"},
			expectedErrors: []*validator.OptionError{
				{Option: "-url",
					Err: "Must be a valid URL, starting with http:// or https://",
				},
			},
		},
		{
			name: "good url format",
			options: options{method: "GET",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "wrong method",
			options: options{method: "TEST",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{
				{Option: "-method",
					Err: "Accepted methods are GET, POST, PUT, DELETE, and PATCH",
				},
			},
		},
		{
			name: "accepted method POST",
			options: options{method: "POST",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "accepted method PUT",
			options: options{method: "PUT",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "accepted method DELETE",
			options: options{method: "DELETE",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "accepted method PATCH",
			options: options{method: "PATCH",
				url: "http://test.com"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "accepted header format",
			options: options{method: "GET",
				url:     "http://test.com",
				headers: "h1:value1, h2:value2, h3:value3"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "wrong header format 1",
			options: options{method: "GET",
				url:     "http://test.com",
				headers: "h1: value1, h2: value2, h3: value3"},
			expectedErrors: []*validator.OptionError{
				{Option: "-headers",
					Err: `Must follow the format "key:value, key:value, ..."`,
				},
			},
		},
		{
			name: "wrong header format 2",
			options: options{method: "GET",
				url:     "http://test.com",
				headers: "header1, header2,header3"},
			expectedErrors: []*validator.OptionError{
				{Option: "-headers",
					Err: `Must follow the format "key:value, key:value, ..."`,
				},
			},
		},
		{
			name: "accepted filenames",
			options: options{method: "GET",
				url: "http://test.com",
				originsFile: struct {
					fileName        string
					origins         []string
					onlyOriginsFile bool
				}{
					fileName: "origins"},
				requestsFile: struct {
					fileName string
					requests []request
				}{
					fileName: "requests1.txt"},
				output:     "output. file",
				outputJSON: "file.json",
				outputCSV:  "csvFile",
				outputYAML: "result"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "wrong filenames",
			options: options{method: "GET",
				url: "http://test.com",
				originsFile: struct {
					fileName        string
					origins         []string
					onlyOriginsFile bool
				}{
					fileName: "origi/ns"},
				requestsFile: struct {
					fileName string
					requests []request
				}{
					fileName: "requests/1.txt"},
				output:     "output. /file",
				outputJSON: "file01234567890123456789.json",
				outputCSV:  "/csvFile",
				outputYAML: "res/ult"},
			expectedErrors: []*validator.OptionError{
				{Option: "-origins-file",
					Err: "A filename cannot contain '/'",
				},
				{Option: "-requests-file",
					Err: "A filename cannot contain '/'",
				},
				{Option: "-output",
					Err: "A filename cannot contain '/'",
				},
				{Option: "-output-json",
					Err: "Cannot exceed 20 characters",
				},
				{Option: "-output-csv",
					Err: "A filename cannot contain '/'",
				},
				{Option: "-output-yaml",
					Err: "A filename cannot contain '/'",
				},
			},
		},
		{
			name: "negative timeout",
			options: options{method: "GET",
				url:     "http://test.com",
				timeout: -1},
			expectedErrors: []*validator.OptionError{
				{Option: "-timeout",
					Err: "Must be greater than 0",
				},
			},
		},
		{
			name: "max 10 timeout",
			options: options{method: "GET",
				url:     "http://test.com",
				timeout: 11},
			expectedErrors: []*validator.OptionError{
				{Option: "-timeout",
					Err: "Must be less than 10",
				},
			},
		},
		{
			name: "negative timedelay",
			options: options{method: "GET",
				url:       "http://test.com",
				timedelay: -1},
			expectedErrors: []*validator.OptionError{
				{Option: "-delay",
					Err: "Must be greater than 0",
				},
			},
		},
		{
			name: "max 5 timedelay",
			options: options{method: "GET",
				url:       "http://test.com",
				timedelay: 6},
			expectedErrors: []*validator.OptionError{
				{Option: "-delay",
					Err: "Must be less than 5",
				},
			},
		},
		{
			name: "wrong proxy",
			options: options{method: "GET",
				url:   "http://test.com",
				proxy: "localhost:8081"},
			expectedErrors: []*validator.OptionError{
				{Option: "-proxy",
					Err: "Must start with http:// or socks5://",
				},
			},
		},
		{
			name: "accepted proxy http://",
			options: options{method: "GET",
				url:   "http://test.com",
				proxy: "http://localhost:8081"},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "accepted proxy socks5://",
			options: options{method: "GET",
				url:   "http://test.com",
				proxy: "socks5://localhost:8081"},
			expectedErrors: []*validator.OptionError{},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			v := &validator.Validator{}

			test.options.validateOptions(v)

			for index, _ := range test.expectedErrors {
				assert.EqualStruct(t, v.Errors[index], test.expectedErrors[index])
			}
		})
	}
}

func TestGetOriginsFromFile(t *testing.T) {
	tests := []struct {
		name            string
		options         options
		expectedOrigins []string
		expectedErrors  []*validator.OptionError
	}{
		{
			name: "originsFile good format",
			options: options{
				originsFile: struct {
					fileName        string
					origins         []string
					onlyOriginsFile bool
				}{
					fileName: createTempFile(
						t,
						"originFile",
						fmt.Sprintf("%s\n%s", "http://example.com", "https://golang.org")).Name(),
				},
			},
			expectedOrigins: []string{
				"http://example.com",
				"https://golang.org",
			},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "originsFile empty",
			options: options{
				originsFile: struct {
					fileName        string
					origins         []string
					onlyOriginsFile bool
				}{
					fileName: createTempFile(
						t,
						"originFile",
						"").Name(),
				},
			},
			expectedOrigins: []string{},
			expectedErrors:  []*validator.OptionError{},
		},
		{
			name: "originsFile bad format",
			options: options{
				originsFile: struct {
					fileName        string
					origins         []string
					onlyOriginsFile bool
				}{
					fileName: createTempFile(
						t,
						"originFile",
						fmt.Sprintf("%s\n%s\n%s", "http://example.com", "", "invalid")).Name(),
				},
			},
			expectedOrigins: []string{
				"http://example.com",
			},
			expectedErrors: []*validator.OptionError{
				{Option: "-origins-file",
					Err: "There cannot be an empty row",
				},
				{Option: "-origins-file",
					Err: `Each origin must be a valid URL, starting with http:// or https://`,
				},
			},
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			defer os.Remove(test.options.originsFile.fileName)

			v := &validator.Validator{}
			test.options.getOriginsFromFile(v)

			for index, _ := range test.expectedOrigins {
				assert.EqualStruct(t, test.options.originsFile.origins[index], test.expectedOrigins[index])
			}

			for index, _ := range test.expectedErrors {
				assert.EqualStruct(t, v.Errors[index], test.expectedErrors[index])
			}
		})
	}
}

func TestGetRequestsFromFile(t *testing.T) {
	tests := []struct {
		name             string
		options          options
		expectedRequests []request
		expectedErrors   []*validator.OptionError
	}{
		{
			name: "requestsFile good format",
			options: options{
				requestsFile: struct {
					fileName string
					requests []request
				}{
					fileName: createTempFile(t, "requestsFile", fmt.Sprintf("%s\n%s",
						`{"url": "https://url1.com", "method": "POST", "headers": {"header1": "value1", "header2": "value2"}, "data": "post data"}`,
						`{"url": "https://url2.com", "method": "GET"}`)).Name(),
				},
			},
			expectedRequests: []request{
				{
					URL:     "https://url1.com",
					Method:  "POST",
					Headers: map[string]string{"header1": "value1", "header2": "value2"},
					Data:    "post data",
				},
				{
					URL:    "https://url2.com",
					Method: "GET",
				},
			},
			expectedErrors: []*validator.OptionError{},
		},
		{
			name: "requestsFile empty",
			options: options{
				requestsFile: struct {
					fileName string
					requests []request
				}{
					fileName: createTempFile(t, "requestsFile", "").Name(),
				},
			},
			expectedRequests: []request{},
			expectedErrors:   []*validator.OptionError{},
		},
		{
			name: "requestsFile bad format",
			options: options{
				requestsFile: struct {
					fileName string
					requests []request
				}{
					fileName: createTempFile(t, "requestsFile", fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
						`{"url": "https://url1.com", "method": "POST", "headers": {"header1": "value1", "header2": "value2"}, "data": "post data"}`,
						"",
						"invalid",
						`{"url": "https:/rl2.com", "method": "GET"}`,
						`{"headers": {"header1": "value1", "header2": "value2"}, "data": "post data"}`,
						`{"ul": "https://url4.com", "method": "GET"}`)).Name(),
				},
			},
			expectedRequests: []request{
				{
					URL:     "https://url1.com",
					Method:  "POST",
					Headers: map[string]string{"header1": "value1", "header2": "value2"},
					Data:    "post data",
				},
			},
			expectedErrors: []*validator.OptionError{
				{Option: "-requests-file",
					Err: "body must not be empty",
				},
				{Option: "-requests-file",
					Err: "body contains badly-formed JSON (at character 1)",
				},
				{Option: "-requests-file",
					Err: "The 'url' key must be a valid URL, starting with http:// or https://",
				},
				{Option: "-requests-file",
					Err: "Each request must contain the 'url' key",
				},
				{Option: "-requests-file",
					Err: `body contains unknown key "ul"`,
				},
			},
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			defer os.Remove(test.options.requestsFile.fileName)

			v := &validator.Validator{}
			test.options.getRequestsFromFile(v)

			for index, _ := range test.expectedRequests {
				assert.EqualStruct(t, test.options.requestsFile.requests[index], test.expectedRequests[index])
			}

			for index, _ := range test.expectedErrors {
				assert.EqualStruct(t, v.Errors[index], test.expectedErrors[index])
			}
		})
	}
}
