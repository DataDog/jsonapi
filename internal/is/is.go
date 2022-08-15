// Package is provides test assertion utilities.
package is

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func errorf(t *testing.T, expected, actual any) {
	t.Helper()
	t.Errorf("\nexpected: %+v\n  actual: %+v", expected, actual)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}

	return false
}

// Nil asserts that the given value is nil.
func Nil(t *testing.T, actual any) bool {
	t.Helper()

	if !isNil(actual) {
		errorf(t, nil, actual)
		return false
	}

	return true
}

func equal(expected, actual any) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}

	return bytes.Equal(exp, act)
}

// Equal asserts that the given values are equal.
func Equal(t *testing.T, expected, actual any) bool {
	t.Helper()

	if !equal(expected, actual) {
		errorf(t, expected, actual)
		return false
	}

	return true
}

// MustEqual requires that the given values are equal.
func MustEqual(t *testing.T, expected, actual any) {
	t.Helper()

	if !Equal(t, expected, actual) {
		t.Fatal()
	}
}

// EqualError asserts that the given error values are equal.
func EqualError(t *testing.T, expected, actual error) bool {
	t.Helper()

	if expected == nil || actual == nil {
		return Equal(t, expected, actual)
	}

	return Equal(t, expected.Error(), actual.Error())
}

// MustEqualError requires that the given error values are equal.
func MustEqualError(t *testing.T, expected, actual error) {
	t.Helper()

	if expected == nil || actual == nil {
		MustEqual(t, expected, actual)
		return
	}

	MustEqual(t, expected.Error(), actual.Error())
}

// NoError asserts that the given err is nil.
func NoError(t *testing.T, err error) bool {
	t.Helper()

	if err != nil {
		errorf(t, nil, err)
		return false
	}

	return true
}

// MustNoError requires that the given err is nil.
func MustNoError(t *testing.T, err error) {
	t.Helper()

	if !NoError(t, err) {
		t.Fatal()
	}
}

// Error asserts that the given err is not nil.
func Error(t *testing.T, err error) bool {
	t.Helper()

	if err == nil {
		errorf(t, "an error", err)
		return false
	}

	return true
}

// MustError requires that the given err is not nil.
func MustError(t *testing.T, err error) {
	t.Helper()

	if !Error(t, err) {
		t.Fatal()
	}
}

// EqualJSON asserts that the given strings are equal after unmarshaling as json.
func EqualJSON(t *testing.T, expected, actual string) bool {
	t.Helper()

	var em map[string]any
	err := json.Unmarshal([]byte(expected), &em)
	MustNoError(t, err)

	var am map[string]any
	err = json.Unmarshal([]byte(actual), &am)
	MustNoError(t, err)

	if !reflect.DeepEqual(em, am) {
		errorf(t, expected, actual)
		return false
	}

	return true
}
