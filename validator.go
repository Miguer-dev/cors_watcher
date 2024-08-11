package main

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

type validator struct {
	Errors map[string]string
}

// Valid returns true if the validator struct doesn't contain any errors.
func (v *validator) valid() bool {
	return len(v.Errors) == 0
}

// AddFieldError adds an error message to the Errors map
func (v *validator) addError(key, message string) {

	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// CheckField adds an error message to the Errors map only if a validation check is not 'ok'.
func (v *validator) check(ok bool, key, message string) {
	if !ok {
		v.addError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string.
func notBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func maxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars returns true if a value contains more than n characters.
func minChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

type number interface {
	constraints.Integer | constraints.Float
}

// NotCero return true if a value is not 0
func notCero[T number](value T) bool {
	return value != 0
}

// MaxNumber returns true if a value is minor that n
func maxNumber[T number](value T, n T) bool {
	return value <= n
}

// MinNumber returns true if a value is greater than n
func minNumber[T number](value T, n T) bool {
	return value >= n
}

// EqualValue returns true if the values are equals
func equalValue[T comparable](value1 T, value2 T) bool {
	return value1 == value2
}

// PermittedValue returns true if a value is in a list of permitted.
func permittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

var urlRX = regexp.MustCompile("((((https?|ftps?|gopher|telnet|nntp)://)|(mailto:|news:))([-%()_.!~*';/?:@&=+$,A-Za-z0-9])+)")
var proxyRX = regexp.MustCompile("^(http://|socks5://)")
var methodRX = regexp.MustCompile("^(GET|POST|PUT|DELETE|PATCH)$")
var headerRX = regexp.MustCompile("^([a-zA-Z0-9!#$%&'*+.^_`" + `|~-]+):[\S]+(, [a-zA-Z0-9!#$%&'*+.^_` + "`|~-]+:[" + `\S]+)*$`)
var fileRX = regexp.MustCompile(`^[^/]*$`)

// Matches return true if match with the pattern
func matches(value string, rxp *regexp.Regexp) bool {
	return rxp.MatchString(value)
}

// Unique return true if all values in a slice are unique.
func unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
