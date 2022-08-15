package jsonapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestParseJSONTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		fieldName   string
		given       any
		expect      string
		expectOK    bool
		expectOmit  bool
	}{
		{
			description: "valid",
			fieldName:   "Foo",
			given: struct {
				Foo string `json:"foo"`
			}{},
			expect:     "foo",
			expectOK:   true,
			expectOmit: false,
		}, {
			description: "valid multiple values",
			fieldName:   "Foo",
			given: struct {
				Foo string `json:"foo,omitempty"`
			}{},
			expect:     "foo",
			expectOK:   true,
			expectOmit: true,
		}, {
			description: "no tag, uses field name",
			fieldName:   "Foo",
			given: struct {
				Foo string
			}{},
			expect:     "Foo",
			expectOK:   true,
			expectOmit: false,
		}, {
			description: "unexported field",
			fieldName:   "foo",
			given: struct {
				foo string
			}{},
			expect:     "",
			expectOK:   false,
			expectOmit: false,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			sf, ok := reflect.TypeOf(tc.given).FieldByName(tc.fieldName)
			is.Equal(t, true, ok)

			tag, ok, omit := parseJSONTag(sf)
			is.Equal(t, tc.expectOK, ok)
			is.Equal(t, tc.expect, tag)
			is.Equal(t, tc.expectOmit, omit)
		})
	}
}

func TestParseJSONAPITag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		given       any
		expect      *tag
		expectError error
	}{
		{
			description: "valid jsonapi, attribute",
			given: struct {
				Foo string `jsonapi:"attribute"`
			}{},
			expect: &tag{directive: attribute},
		}, {
			description: "valid jsonapi, relationship",
			given: struct {
				Foo string `jsonapi:"relationship"`
			}{},
			expect: &tag{directive: relationship},
		}, {
			description: "valid jsonapi, primary",
			given: struct {
				Foo string `jsonapi:"primary,foo"`
			}{},
			expect: &tag{directive: primary, resourceType: "foo"},
		}, {
			description: "no struct tags",
			given:       struct{ Foo string }{},
			expect:      nil,
		}, {
			description: "invalid jsonapi tag (missing value)",
			given: struct {
				Foo string `jsonapi:"primary"`
			}{},
			expect: nil,
			expectError: &TagError{
				TagName: "jsonapi",
				Field:   "Foo",
				Reason:  "missing type in primary directive",
			},
		}, {
			description: "invalid jsonapi tag (invalid directive)",
			given: struct {
				Foo string `jsonapi:"invalid,foo"`
			}{},
			expect: nil,
			expectError: &TagError{
				TagName: "jsonapi",
				Field:   "Foo",
				Reason:  "invalid directive",
			},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			sf, ok := reflect.TypeOf(tc.given).FieldByName("Foo")
			is.Equal(t, true, ok)

			tag, err := parseJSONAPITag(sf)
			is.MustEqualError(t, tc.expectError, err)
			is.Equal(t, tc.expect, tag)
		})
	}
}
