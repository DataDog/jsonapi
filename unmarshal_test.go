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
			description: "*Article (no id)",
			given:       articleANoIDBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleANoID,
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
			description: "[]Article (nonunique data)",
			given:       articlesABNonuniqueData,
			do: func(body []byte) (any, error) {
				var a []Article
				err := Unmarshal(body, &a, UnmarshalCheckUniqueness())
				return a, err
			},
			expect:      ([]Article)(nil),
			expectError: ErrNonuniqueResource,
		}, {
			description: "[]Article (empty)",
			given:       emptyManyBody,
			do: func(body []byte) (any, error) {
				var a []Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []Article{},
			expectError: nil,
		}, {
			description: "[]*Article (empty)",
			given:       emptyManyBody,
			do: func(body []byte) (any, error) {
				var a []*Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []*Article{},
			expectError: nil,
		}, {
			description: "Article (empty)",
			given:       emptySingleBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      Article{},
			expectError: ErrEmptyDataObject,
		}, {
			description: "*Article (empty)",
			given:       emptySingleBody,
			do: func(body []byte) (any, error) {
				var a *Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      (*Article)(nil),
			expectError: ErrEmptyDataObject,
		}, {
			description: "Article null data",
			given:       nullDataBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      Article{},
			expectError: nil,
		}, {
			description: "*Article null data",
			given:       nullDataBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &Article{},
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
			description: "*ArticleUnmarshalType",
			given:       articleACustomTypeBody,
			do: func(body []byte) (any, error) {
				var a ArticleTypeOverride
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleACustomType,
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
			description: "*ArticleEncodingIntID",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a ArticleEncodingIntID
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleAEncodingIntID,
			expectError: nil,
		}, {
			description: "[]*ArticleEncodingIntID",
			given:       articlesABBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleEncodingIntID
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      []*ArticleEncodingIntID{&articleAEncodingIntID, &articleBEncodingIntID},
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
			description: "ArticleEmbedded",
			given:       articleEmbeddedBody,
			do: func(body []byte) (any, error) {
				var a ArticleEmbedded
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleEmbedded,
			expectError: nil,
		}, {
			description: "ArticleEmbeddedPointer",
			given:       articleEmbeddedBody,
			do: func(body []byte) (any, error) {
				var a ArticleEmbeddedPointer
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleEmbeddedPointer,
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
			expectError: ErrDocumentMissingRequiredMembers,
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
		}, {
			description: "*ArticleDoubleID invalid",
			given:       articleABody,
			do: func(body []byte) (any, error) {
				var a ArticleDoubleID
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleDoubleID{ID: "1"},
			expectError: ErrUnmarshalDuplicatePrimaryField,
		}, {
			description: "*Article with included author (not linked)",
			given:       articleWithIncludeOnlyBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      new(Article),
			expectError: &PartialLinkageError{[]string{"{Type: author, ID: 1}"}},
		}, {
			description: "*ArticleRelated empty relationships object",
			given:       articleRelatedInvalidEmptyRelationshipBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{},
			expectError: ErrRelationshipMissingRequiredMembers,
		}, {
			// verifies for empty relationship data: null -> nil and [] -> []Type{}
			description: "*ArticleRelated empty relationships data (valid)",
			given:       articleRelatedNoOmitEmptyBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{ID: "1", Title: "A", Author: nil, Comments: []*Comment{}},
			expectError: nil,
		}, {
			// this test verifies that empty relationship data objects do not unmarshal
			description: "*ArticleRelated empty relationships data (invalid)",
			given:       articleRelatedInvalidEmptyDataBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{},
			expectError: ErrEmptyDataObject,
		}, {
			description: "*ArticleRelated with nil relationships data",
			given:       articleRelatedLiteralNullData,
			do: func(body []byte) (any, error) {
				a := ArticleRelated{
					Author:   &Author{ID: "100"},
					Comments: []*Comment{{ID: "200"}},
				}
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{ID: "1", Title: "A", Author: nil, Comments: nil},
			expectError: nil,
		}, {
			description: "*ArticleRelated.Author",
			given:       articleRelatedAuthorBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect: &ArticleRelated{
				ID:     "1",
				Title:  "A",
				Author: &Author{ID: "1"},
			},
			expectError: nil,
		}, {
			description: "*ArticleRelated.Author (links only)",
			given:       articleRelatedAuthorLinksOnlyBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{ID: "1", Title: "A"},
			expectError: nil,
		}, {
			description: "*ArticleRelated.Author (meta only)",
			given:       articleRelatedAuthorMetaOnlyBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{ID: "1", Title: "A"},
			expectError: nil,
		}, {
			description: "[]*ArticleRelated.Author twice",
			given:       articleRelatedAuthorTwiceBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect: &[]*ArticleRelated{
				{ID: "1", Title: "A", Author: &Author{ID: "1"}},
				{ID: "2", Title: "B", Author: &Author{ID: "1"}},
			},
			expectError: nil,
		}, {
			description: "*ArticleRelated Complete",
			given:       articleRelatedCompleteBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect: &ArticleRelated{
				ID:       "1",
				Title:    "A",
				Author:   &Author{ID: "1"},
				Comments: []*Comment{{ID: "1"}, {ID: "2"}},
			},
			expectError: nil,
		}, {
			description: "[]*ArticleRelated.Author twice with include",
			given:       articleRelatedAuthorTwiceWithIncludeBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect: &[]*ArticleRelated{
				{ID: "1", Title: "A", Author: &authorA},
				{ID: "2", Title: "B", Author: &authorA},
			},
			expectError: nil,
		}, {
			description: "[]*ArticleRelated complete with include",
			given:       articleRelatedCompleteWithIncludeBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &ArticleRelated{ID: "1", Title: "A", Author: &authorA, Comments: commentsAB},
			expectError: nil,
		}, {
			description: "*ArticleRelated.Comments.Author with include",
			given:       articleRelatedCommentsNestedWithIncludeBody,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleRelatedCommentsNested,
			expectError: nil,
		}, {
			description: "ArticleRelated (nonunique linkage)",
			given:       articleRelatedNonuniqueLinkage,
			do: func(body []byte) (any, error) {
				var a ArticleRelated
				err := Unmarshal(body, &a, UnmarshalCheckUniqueness())
				return a, err
			},
			expect:      ArticleRelated{},
			expectError: ErrNonuniqueResource,
		},
		{
			description: "ArticleRelatedWithContact (nonunique linkage)",
			given:       articleRelatedDifferentRefToSameObject,
			do: func(body []byte) (any, error) {
				var a ArticleRelatedWithContact
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect: &ArticleRelatedWithContact{
				ID:       "1",
				Title:    "A",
				Author:   &Author{ID: "1"},
				Comments: []*Comment{{ID: "1"}},
				Contact:  &Author{ID: "1"},
			},
			expectError: nil,
		},
		{
			description: "links member only",
			given:       `{"links":null}`,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &Article{},
			expectError: ErrDocumentMissingRequiredMembers,
		}, {
			description: "[]*ArticleRelated complex relationships with include",
			given:       articlesRelatedComplexBody,
			do: func(body []byte) (any, error) {
				var a []*ArticleRelated
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articlesRelatedComplex,
			expectError: nil,
		}, {
			// TODO(#36): Add unmarshaling of errors as a supported use-case
			description: "Errors don't unmarshal",
			given:       errorsSimpleStructBody,
			do: func(body []byte) (any, error) {
				var a Article
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &Article{},
			expectError: ErrErrorUnmarshalingNotImplemented,
		}, {
			description: "CommentEmbedded",
			given:       commentEmbeddedBody,
			do: func(body []byte) (any, error) {
				var a CommentEmbedded
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &commentEmbedded,
			expectError: nil,
		},
		{
			description: "CommentEmbeddedPointer",
			given:       commentEmbeddedBody,
			do: func(body []byte) (any, error) {
				var a CommentEmbeddedPointer
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &commentEmbeddedPointer,
			expectError: nil,
		},
		{
			description: "ArticleLinkedOnlySelf",
			given:       articleLinkedOnlySelfBody,
			do: func(body []byte) (any, error) {
				var a ArticleLinkedOnlySelf
				err := Unmarshal(body, &a)
				return &a, err
			},
			expect:      &articleLinkedOnlySelf,
			expectError: nil,
		},
		{
			description: "*ArticleLinkedOnlySelf",
			given:       articleLinkedOnlySelfBody,
			do: func(body []byte) (any, error) {
				var a *ArticleLinkedOnlySelf
				err := Unmarshal(body, &a)
				return a, err
			},
			expect:      &articleLinkedOnlySelf,
			expectError: nil,
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

	articlesABToplevelMetaBody := `{"data":[{"type":"articles","id":"1","attributes":{"title":"A"}},{"type":"articles","id":"2","attributes":{"title":"B"}}],"meta":{"foo":"bar"}}`
	articleAInvalidToplevelMetaBody := `{"data":{"id":"1","type":"articles"},"meta":"foo"}`

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
				err := Unmarshal([]byte(articleAToplevelMetaBody), &a, UnmarshalMeta(&m))
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
				err := Unmarshal([]byte(articleAToplevelMetaBody), &a, UnmarshalMeta(&m))
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
				err := Unmarshal([]byte(articlesABToplevelMetaBody), &a, UnmarshalMeta(&m))
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
				err := Unmarshal([]byte(articleAToplevelMetaBody), &a, UnmarshalMeta(m))
				return &m, err
			},
			expect:      new(meta),
			expectError: &json.InvalidUnmarshalError{Type: reflect.TypeOf(meta{})},
		}, {
			description: "invalid meta type",
			do: func() (any, error) {
				var (
					a Article
					m string
				)
				err := Unmarshal([]byte(articleAInvalidToplevelMetaBody), &a, UnmarshalMeta(m))
				return &m, err
			},
			expect:      nil,
			expectError: &json.InvalidUnmarshalError{Type: reflect.TypeOf("")},
		}, {
			description: "meta (empty Article)",
			do: func() (any, error) {
				var (
					a Article
					m meta
				)
				err := Unmarshal([]byte(articleNullWithToplevelMetaBody), &a, UnmarshalMeta(&m))
				return &m, err
			},
			expect:      &meta{Foo: "bar"},
			expectError: nil,
		}, {
			description: "meta (empty []*Article)",
			do: func() (any, error) {
				var (
					a []*Article
					m meta
				)
				err := Unmarshal([]byte(articleEmptyArrayWithToplevelMetaBody), &a, UnmarshalMeta(&m))
				return &m, err
			},
			expect:      &meta{Foo: "bar"},
			expectError: nil,
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

				if tc.expect != nil {
					is.Equal(t, tc.expect, actual)
				}
				return
			}
			is.MustNoError(t, err)
			is.Equal(t, tc.expect, actual)
		})
	}
}

func TestUnmarshalLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		do          func() (*Link, error)
		expected    *Link
	}{
		{
			description: "null",
			do: func() (*Link, error) {
				var (
					a []Article
					l Link
				)
				err := Unmarshal([]byte(articlesABBody), &a, UnmarshalLinks(&l))
				return &l, err
			},
			expected: &Link{},
		},
		{
			description: "*Link",
			do: func() (*Link, error) {
				var (
					a Article
					l Link
				)
				err := Unmarshal([]byte(articleAWithTopLevelLink), &a, UnmarshalLinks(&l))
				return &l, err
			},
			expected: &Link{
				Self: "http://example.com/article/1",
			},
		},
		{
			description: "*Link with related(any) and meta(any)",
			do: func() (*Link, error) {
				var (
					a Article
					l Link
				)
				err := Unmarshal([]byte(articleAWithTopLevelLinksRelated), &a, UnmarshalLinks(&l))
				return &l, err
			},
			expected: &Link{
				Related: map[string]interface{}{
					"href": "http://example.com/article/1/comments",
					"meta": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},
	}

	for i, tc := range tests {
		tc := tc

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := tc.do()
			is.NoError(t, err)
			is.Equal(t, tc.expected, actual)
		})
	}
}

func TestUnmarshalNumbers(t *testing.T) {
	t.Parallel()

	type nested struct {
		Max int64 `json:"max"`
	}
	type numbers struct {
		ID     string         `jsonapi:"primary,numbers"`
		Int    int64          `jsonapi:"attribute" json:"int,omitempty"`
		Uint   uint64         `jsonapi:"attribute" json:"uint,omitempty"`
		Float  float64        `jsonapi:"attribute" json:"float,omitempty"`
		Nested *nested        `jsonapi:"attribute" json:"nested,omitempty"`
		List   []int64        `jsonapi:"attribute" json:"list,omitempty"`
		Values map[string]any `jsonapi:"attribute" json:"values,omitempty"`
		Int32  int32          `jsonapi:"attribute" json:"int32,omitempty"`
	}

	body := func(attributes string) []byte {
		return []byte(`{"data":{"type":"numbers","id":"1","attributes":` + attributes + `}}`)
	}

	tests := []struct {
		description string
		given       []byte
		opts        []UnmarshalOption
		expect      any
	}{
		{
			description: "int64 above 2^53",
			given:       body(`{"int":9007199254740993}`),
			expect:      &numbers{ID: "1", Int: 9007199254740993},
		}, {
			description: "uint64 max",
			given:       body(`{"uint":18446744073709551615}`),
			expect:      &numbers{ID: "1", Uint: 18446744073709551615},
		}, {
			description: "int64 min and max in nested struct and slice",
			given:       body(`{"nested":{"max":9223372036854775807},"list":[-9223372036854775808,9007199254740993]}`),
			expect:      &numbers{ID: "1", Nested: &nested{Max: 9223372036854775807}, List: []int64{-9223372036854775808, 9007199254740993}},
		}, {
			description: "non-integer literals with integer values decode into int64",
			given:       body(`{"int":100.0,"list":[1e2]}`),
			expect:      &numbers{ID: "1", Int: 100, List: []int64{100}},
		}, {
			description: "non-integer literals with integer values decode into int64 next to an integer above 2^53",
			given:       body(`{"int":100.0,"uint":9007199254740993}`),
			expect:      &numbers{ID: "1", Int: 100, Uint: 9007199254740993},
		}, {
			description: "integer above 2^53 into float64 rounds",
			given:       body(`{"float":9007199254740993}`),
			expect:      &numbers{ID: "1", Float: 9007199254740992},
		}, {
			description: "map[string]any values are float64",
			given:       body(`{"values":{"int":9007199254740993,"float":1.5}}`),
			expect:      &numbers{ID: "1", Values: map[string]any{"int": float64(9007199254740992), "float": 1.5}},
		}, {
			description: "map[string]any values are json.Number with UnmarshalUseNumber",
			given:       body(`{"values":{"int":9007199254740993,"float":1.5}}`),
			opts:        []UnmarshalOption{UnmarshalUseNumber()},
			expect:      &numbers{ID: "1", Values: map[string]any{"int": json.Number("9007199254740993"), "float": json.Number("1.5")}},
		}, {
			description: "non-integer literals with integer values decode into int64 with UnmarshalUseNumber",
			given:       body(`{"int":100.0}`),
			opts:        []UnmarshalOption{UnmarshalUseNumber()},
			expect:      &numbers{ID: "1", Int: 100},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			var actual numbers
			err := Unmarshal(tc.given, &actual, tc.opts...)
			is.MustNoError(t, err)
			is.Equal(t, tc.expect, &actual)
		})
	}

	t.Run("int32 out of range", func(t *testing.T) {
		t.Parallel()

		var actual numbers
		err := Unmarshal(body(`{"int32":3000000000}`), &actual)
		is.Error(t, err)
	})
}

// TestUnmarshalMemberNameValidation collects tests which verify that invalid member names are
// caught during unmarshaling, no matter where they're placed. This test does not exhaustively test
// every possible invalid name.
func TestUnmarshalMemberNameValidation(t *testing.T) {
	t.Parallel()

	articleWithInvalidToplevelMetaMemberNameBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"meta":{"foo%":2}}`
	articleWithInvalidJSONAPIMetaMemberNameBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"jsonapi":{"version":"1.0","meta":{"foo%":1}}}`
	articleWithInvalidRelationshipAttributeNameNotIncludedBody := `{"data":{"id":"1","type":"articles","relationships":{"author":{"data":{"id":"1","type":"author"}}}}}`
	articlesWithOneInvalidResourceMetaMemberName := `{"data":[{"id":"1","type":"articles"},{"id":"2","type":"articles","meta":{"foo%":1}}]}`

	tests := []struct {
		description string
		given       string
		do          func(body []byte, opts ...UnmarshalOption) error
		expectError error
	}{
		{
			description: "Article with valid member names",
			given:       articleABody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a Article
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: nil,
		}, {
			description: "Author with invalid type name",
			given:       authorWithInvalidTypeNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a AuthorWithInvalidTypeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Author with invalid attribute member name",
			given:       authorWithInvalidAttributeNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a AuthorWithInvalidAttributeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"na%me"},
		}, {
			description: "Article with invalid resource meta member name",
			given:       articleWithInvalidResourceMetaMemberNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithGenericMeta
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Article with invalid top-level meta member name",
			given:       articleWithInvalidToplevelMetaMemberNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var (
					a Article
					m map[string]any
				)
				opts = append(opts, UnmarshalMeta(&m))
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Article with invalid link meta member name",
			given:       articleWithInvalidLinkMetaMemberNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithInvalidLinkMetaMemberName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Article with invalid jsonapi meta member name",
			given:       articleWithInvalidJSONAPIMetaMemberNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a Article
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Article with invalid relationship name",
			given:       articleWithInvalidRelationshipNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithInvalidRelationshipName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Article with invalid relationship type name body",
			given:       articleWithInvalidRelationshipTypeNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithInvalidRelationshipTypeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Article with invalid relationship attribute member names not included",
			given:       articleWithInvalidRelationshipAttributeNameNotIncludedBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithInvalidRelationshipAttributeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: nil,
		}, {
			description: "Article with invalid relationship attribute member names included",
			given:       articleWithInvalidRelationshipAttributeNameIncludedBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a ArticleWithInvalidRelationshipAttributeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"na%me"},
		}, {
			description: "[]*Article with one invalid resource meta member name",
			given:       articlesWithOneInvalidResourceMetaMemberName,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a []*ArticleWithGenericMeta
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Website with invalid nested relationship type member name",
			given:       websiteWithInvalidNestedRelationshipTypeNameBody,
			do: func(body []byte, opts ...UnmarshalOption) error {
				var a WebsiteWithInvalidNestedRelationshipTypeName
				err := Unmarshal(body, &a, opts...)
				return err
			},
			expectError: &MemberNameValidationError{"aut%hor"},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			err := tc.do([]byte(tc.given))
			is.EqualError(t, tc.expectError, err)

			err = tc.do([]byte(tc.given), UnmarshalSetNameValidation(DisableValidation))
			is.MustNoError(t, err)
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	benchmarks := []struct {
		name   string
		data   string
		target any
		opts   []UnmarshalOption
	}{
		{
			name:   "ArticleSimple",
			data:   articleABody,
			target: Article{},
			opts:   nil,
		}, {
			name:   "ArticleSimpleWithToplevelMeta",
			data:   articleAToplevelMetaBody,
			target: Article{},
			opts:   []UnmarshalOption{UnmarshalMeta(map[string]any{"foo": "bar"})},
		}, {
			name:   "ArticleComplex",
			data:   articleRelatedCommentsNestedWithIncludeBody,
			target: ArticleRelated{},
			opts:   nil,
		}, {
			name:   "ArticleComplexDisableNameValidation",
			data:   articleRelatedCommentsNestedWithIncludeBody,
			target: ArticleRelated{},
			opts:   []UnmarshalOption{UnmarshalSetNameValidation(DisableValidation)},
		}, {
			name:   "ArticlesComplex",
			data:   articlesRelatedComplexBody,
			target: []*ArticleRelated{},
			opts:   nil,
		}, {
			name:   "ArticlesComplexDisableNameValidation",
			data:   articlesRelatedComplexBody,
			target: []*ArticleRelated{},
			opts:   []UnmarshalOption{UnmarshalSetNameValidation(DisableValidation)},
		},
	}

	for _, bm := range benchmarks {
		bm := bm
		data := []byte(bm.data)
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = Unmarshal(data, &bm.target, bm.opts...)
			}
		})
	}
}
