package assert

import (
	"reflect"
	"strings"
	"testing"
)

// Equal fail test if "actual" is different to "expected"
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// EqualStruct fail test if "actual" is different to "expected"
func EqualStruct(t *testing.T, actual, expected any) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// StringContains fail test if "actual" not contain "expectedSubstring"
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// NilError fail if actual error is not nil
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
