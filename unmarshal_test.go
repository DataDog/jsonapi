package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		given       string
		do          func(body []byte) (any, error)
		expect      any
		expectError error
	}{
		{
			description: "*ArticleComplete",
			given:       articleCompleteBody,
			do: func(body []byte) (any, error) {
				var a ArticleComplete
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleComplete,
			expectError: nil,
		}, {
			description: "*Article",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleA,
			expectError: nil,
		}, {
			description: "**Article",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a *Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      &articleA,
			expectError: nil,
		}, {
			description: "[]Article",
			given:       articlesABBody,
			do: func(body []byte) (any, error) {
				var a []Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []Article{articleA, articleB},
			expectError: nil,
		}, {
			description: "[]*Article",
			given:       articlesABBody,
			do: func(body []byte) (any, error) {
				var a []*Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []*Article{&articleA, &articleB},
			expectError: nil,
		}, {
			description: "*ArticleIntID",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a ArticleIntID
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleAIntID,
			expectError: nil,
		}, {
			description: "[]*ArticleIntID",
			given:       articlesABBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleIntID
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []*ArticleIntID{&articleAIntID, &articleBIntID},
			expectError: nil,
		}, {
			description: "*ArticleIntIDID",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a ArticleIntIDID
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleAIntIDID,
			expectError: nil,
		}, {
			description: "[]*ArticleIntIDID",
			given:       articlesABBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleIntIDID
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []*ArticleIntIDID{&articleAIntIDID, &articleBIntIDID},
			expectError: nil,
		}, {
			description: "*ArticleWithMeta",
			given:       articleAWithMetaBody,
			do: func(body []byte) (any, error) {
				var a ArticleWithMeta
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleAWithMeta,
			expectError: nil,
		}, {
			description: "nil",
			given:       "",
			do: func(body []byte) (any, error) {
				err := Unmarshal(body, nil)
				return nil, err
			},
			expect:      nil,
			expectError: &TypeError{Actual: "invalid", Expected: []string{"non-nil pointer"}},
		}, {
			description: "not a pointer",
			given:       "",
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, a)
				return &a, err
			},
			expect:      new(Article),
			expectError: &TypeError{Actual: "struct", Expected: []string{"non-nil pointer"}},
		}, {
			description: "empty json body",
			given:       "{}",
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      new(Article),
			expectError: &RequestBodyError{t: new(Article)},
		}, {
			description: "*Article (invalid type)",
			given:       articleAInvalidTypeBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      new(Article),
			expectError: &TypeError{Actual: "not-articles", Expected: []string{"articles"}},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := tc.do([]byte(tc.given))
			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Equal(t, tc.expect, actual)
				return
			}
			is.MustNoError(t, err)
			is.Equal(t, tc.expect, actual)
		})
	}
}

func TestUnmarshalMeta(t *testing.T) {
	t.Parallel()

	articleAMetaBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"meta":{"foo":"bar"}}`
	articlesABMetaBody := `{"data":[{"type":"articles","id":"1","attributes":{"title":"A"}},{"type":"articles","id":"2","attributes":{"title":"B"}}],"meta":{"foo":"bar"}}`

	type meta struct {
		Foo string `json:"foo"`
	}

	tests := []struct {
		description string
		do          func() (any, error)
		expect      any
		expectError error
	}{
		{
			description: "map[string]any",
			do: func() (any, error) {
				var (
					a Article
					m map[string]any
				)
				err := Unmarshal([]byte(articleAMetaBody), &a, UnmarshalMeta(&m))
				return m, err
			},
			expect:      map[string]any{"foo": "bar"},
			expectError: nil,
		}, {
			description: "*meta (*Article)",
			do: func() (any, error) {
				var (
					a Article
					m meta
				)
				err := Unmarshal([]byte(articleAMetaBody), &a, UnmarshalMeta(&m))
				return &m, err
			},
			expect:      &meta{Foo: "bar"},
			expectError: nil,
		}, {
			description: "*meta ([]*Article)",
			do: func() (any, error) {
				var (
					a []*Article
					m meta
				)
				err := Unmarshal([]byte(articlesABMetaBody), &a, UnmarshalMeta(&m))
				return &m, err
			},
			expect:      &meta{Foo: "bar"},
			expectError: nil,
		}, {
			description: "meta",
			do: func() (any, error) {
				var (
					a Article
					m meta
				)
				err := Unmarshal([]byte(articleAMetaBody), &a, UnmarshalMeta(m))
				return &m, err
			},
			expect:      new(meta),
			expectError: &json.InvalidUnmarshalError{Type: reflect.TypeOf(meta{})},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := tc.do()
			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Equal(t, tc.expect, actual)
				return
			}
			is.MustNoError(t, err)
			is.Equal(t, tc.expect, actual)
		})
	}
}
