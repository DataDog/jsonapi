package jsonapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestDerefValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		do          func() reflect.Value
		expectType  string
		expectKind  string
	}{
		{
			description: "string -> string",
			do: func() reflect.Value {
				var s string
				return derefValue(reflect.ValueOf(s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "*string -> string",
			do: func() reflect.Value {
				var s string
				return derefValue(reflect.ValueOf(&s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "**string -> string",
			do: func() reflect.Value {
				var s *string
				return derefValue(reflect.ValueOf(&s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "Article -> Article",
			do: func() reflect.Value {
				var a Article
				return derefValue(reflect.ValueOf(a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		}, {
			description: "*Article -> Article",
			do: func() reflect.Value {
				var a Article
				return derefValue(reflect.ValueOf(&a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		}, {
			description: "**Article -> Article",
			do: func() reflect.Value {
				var a *Article
				return derefValue(reflect.ValueOf(&a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual := tc.do()
			is.Equal(t, tc.expectKind, actual.Kind().String())
			is.Equal(t, tc.expectType, actual.Type().String())
		})
	}
}

func TestDerefType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		do          func() reflect.Type
		expectType  string
		expectKind  string
	}{
		{
			description: "string -> string",
			do: func() reflect.Type {
				var s string
				return derefType(reflect.TypeOf(s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "*string -> string",
			do: func() reflect.Type {
				var s string
				return derefType(reflect.TypeOf(&s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "**string -> string",
			do: func() reflect.Type {
				var s *string
				return derefType(reflect.TypeOf(&s))
			},
			expectKind: "string",
			expectType: "string",
		}, {
			description: "Article -> Article",
			do: func() reflect.Type {
				var a Article
				return derefType(reflect.TypeOf(a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		}, {
			description: "*Article -> Article",
			do: func() reflect.Type {
				var a Article
				return derefType(reflect.TypeOf(&a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		}, {
			description: "**Article -> Article",
			do: func() reflect.Type {
				var a *Article
				return derefType(reflect.TypeOf(&a))
			},
			expectKind: "struct",
			expectType: "jsonapi.Article",
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual := tc.do()
			is.Equal(t, tc.expectKind, actual.Kind().String())
			is.Equal(t, tc.expectType, actual.String())
		})
	}
}
