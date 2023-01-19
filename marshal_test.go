package jsonapi

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	articleAPtr := &articleA
	articleBPtr := &articleB
	articlePtrSlicePtr := &[]*Article{articleAPtr, articleBPtr}

	tests := []struct {
		description string
		given       any
		expect      string
		expectError error
	}{
		{
			description: "nil",
			given:       nil,
			expect:      `{"data":null}`,
			expectError: nil,
		}, {
			description: "Article (empty)",
			given:       Article{},
			expect:      `{"data":null}`,
			expectError: nil,
		}, {
			description: "*Article (empty)",
			given:       new(Article),
			expect:      "",
			expectError: ErrEmptyPrimaryField,
		}, {
			description: "[]*Article (empty)",
			given:       make([]*Article, 0),
			expect:      emptyBody,
			expectError: nil,
		}, {
			description: "Article (missing ID)",
			given:       Article{Title: "A"},
			expect:      "",
			expectError: ErrEmptyPrimaryField,
		}, {
			description: "Article",
			given:       articleA,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "*Article",
			given:       &articleA,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "*ArticleComplete",
			given:       &articleComplete,
			expect:      articleCompleteBody,
			expectError: nil,
		}, {
			description: "**Article",
			given:       &articleA,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "[]Article",
			given:       articlesAB,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "*[]Article",
			given:       &articlesAB,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "[]*Article",
			given:       articlesABPtr,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "*[]*Article",
			given:       articlePtrSlicePtr,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "**[]*Article",
			given:       &articlePtrSlicePtr,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "*ArticleLinked",
			given:       &articleALinked,
			expect:      articleALinkedBody,
			expectError: nil,
		}, {
			description: "*ArticleLinkedOnlySelf",
			given:       &articleLinkedOnlySelf,
			expect:      articleLinkedOnlySelfBody,
			expectError: nil,
		}, {
			description: "*ArticleOmitTitle (full)",
			given:       &articleOmitTitleFull,
			expect:      articleOmitTitleFullBody,
			expectError: nil,
		}, {
			description: "*ArticleOmitTitle (partial)",
			given:       &articleOmitTitlePartial,
			expect:      articleOmitTitlePartialBody,
			expectError: nil,
		}, {
			description: "invalid Link.Self",
			given:       &articleLinkedInvalidSelf,
			expect:      "",
			expectError: &TypeError{Actual: "int", Expected: []string{"*LinkObject", "string"}},
		}, {
			description: "invalid Link.Related",
			given:       &articleLinkedInvalidRelated,
			expect:      "",
			expectError: &TypeError{Actual: "int", Expected: []string{"*LinkObject", "string"}},
		}, {
			description: "invalid Link (nil Link.Self and Link.Related)",
			given:       &articleLinkedInvalidMissingFields,
			expect:      "",
			expectError: ErrMissingLinkFields,
		}, {
			description: "invalid Link (empty Link.Self and Link.Related)",
			given:       &articleLinkedInvalidMissingFieldsEmptyValues,
			expect:      "",
			expectError: ErrMissingLinkFields,
		}, {
			description: "string",
			given:       "a",
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "slice"}},
		}, {
			description: "[]string",
			given:       []string{"a", "b"},
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct"}},
		}, {
			description: "*ArticleIntID (MarshalIdentifier)",
			given:       &articleAIntID,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "[]*ArticleIntID (MarshalIdentifier)",
			given:       articlesIntIDABPtr,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "*ArticleIntIDID (fmt.Stringer)",
			given:       &articleAIntIDID,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "[]*ArticleIntIDID (fmt.Stringer)",
			given:       &articlesIntIDIDABPtr,
			expect:      articlesABBody,
			expectError: nil,
		}, {
			description: "non-string id",
			given: &struct {
				ID int `jsonapi:"primary,test"`
			}{ID: 1},
			expect:      "",
			expectError: ErrMarshalInvalidPrimaryField,
		}, {
			description: "missing primary tag",
			given: &struct {
				ID string `jsonapi:"attr,id"`
			}{ID: "1"},
			expect:      "",
			expectError: ErrMissingPrimaryField,
		}, {
			description: "ArticleWithResourceObjectMeta",
			given:       &articleWithResourceObjectMeta,
			expect:      articleWithResourceObjectMetaBody,
			expectError: nil,
		}, {
			description: "ArticleWithoutResourceObjectMeta",
			given:       &articleWithoutResourceObjectMeta,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "Error simple",
			given:       errorsSimpleStruct,
			expect:      errorsSimpleStructBody,
			expectError: nil,
		}, {
			description: "*Error simple",
			given:       &errorsSimpleStruct,
			expect:      errorsSimpleStructBody,
			expectError: nil,
		}, {
			description: "[]Error simple single",
			given:       errorsSimpleSliceSingle,
			expect:      errorsSimpleStructBody,
			expectError: nil,
		}, {
			description: "[]*Error simple single",
			given:       errorsSimpleSliceSinglePtr,
			expect:      errorsSimpleStructBody,
			expectError: nil,
		}, {
			description: "Error complex",
			given:       errorsComplexStruct,
			expect:      errorsComplexStructBody,
			expectError: nil,
		}, {
			description: "[]Error complex many",
			given:       errorsComplexSliceMany,
			expect:      errorsComplexSliceManyBody,
			expectError: nil,
		}, {
			description: "[]*Error complex many",
			given:       errorsComplexSliceManyPtr,
			expect:      errorsComplexSliceManyBody,
			expectError: nil,
		}, {
			description: "Error with link object",
			given:       errorsWithLinkObject,
			expect:      errorsWithLinkObjectBody,
			expectError: nil,
		}, {
			description: "Error with invalid link",
			given:       errorsWithInvalidLink,
			expect:      "",
			expectError: &TypeError{Actual: "int", Expected: []string{"*LinkObject", "string"}},
		}, {
			description: "Error empty",
			given:       Error{},
			expect:      `{"errors":[{}]}`,
			expectError: nil,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given)
			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Nil(t, actual)
				return
			}
			is.MustNoError(t, err)
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalMeta(t *testing.T) {
	t.Parallel()

	articleAMetaBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"meta":{"foo":"bar"}}`
	errorsObjectMetaBody := `{"meta":{"foo":"bar"},"errors":[{"title":"T"}]}`

	tests := []struct {
		description string
		given       any
		givenMeta   any
		expect      string
		expectError error
	}{
		{
			description: "map[string]any",
			given:       &articleA,
			givenMeta:   map[string]any{"foo": "bar"},
			expect:      articleAMetaBody,
			expectError: nil,
		}, {
			description: "struct",
			given:       &articleA,
			givenMeta: &struct {
				Foo string `json:"foo"`
			}{Foo: "bar"},
			expect:      articleAMetaBody,
			expectError: nil,
		}, {
			description: "non-object type",
			given:       &articleA,
			givenMeta:   "foo",
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
		}, {
			description: "map[string]any errors object",
			given:       errorsSimpleSliceSinglePtr,
			givenMeta:   map[string]any{"foo": "bar"},
			expect:      errorsObjectMetaBody,
			expectError: nil,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given, MarshalMeta(tc.givenMeta))
			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Nil(t, actual)
				return
			}
			is.MustNoError(t, err)
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalJSONAPI(t *testing.T) {
	t.Parallel()

	articleAJSONAPIBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"jsonapi":{"version":"1.0"}}`
	articleAJSONAPIMetaBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"jsonapi":{"version":"1.0","meta":{"foo":"bar"}}}`
	errorsObjectMetaBody := `{"jsonapi":{"version":"1.0","meta":{"foo":"bar"}},"errors":[{"title":"T"}]}`

	tests := []struct {
		description  string
		given        any
		givenJSONAPI any
		expect       string
		expectError  error
	}{
		{
			description:  "include jsonapi",
			given:        &articleA,
			givenJSONAPI: nil,
			expect:       articleAJSONAPIBody,
			expectError:  nil,
		}, {
			description:  "include jsonapi and meta (map)",
			given:        &articleA,
			givenJSONAPI: map[string]any{"foo": "bar"},
			expect:       articleAJSONAPIMetaBody,
			expectError:  nil,
		}, {
			description: "include jsonapi and meta (struct)",
			given:       &articleA,
			givenJSONAPI: &struct {
				Foo string `json:"foo"`
			}{Foo: "bar"},
			expect:      articleAJSONAPIMetaBody,
			expectError: nil,
		}, {
			description:  "include jsonapi and meta (non-object type)",
			given:        &articleA,
			givenJSONAPI: "foo",
			expect:       "",
			expectError:  &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
		}, {
			description:  "include jsonapi and meta (map) errors object",
			given:        errorsSimpleSliceSinglePtr,
			givenJSONAPI: map[string]any{"foo": "bar"},
			expect:       errorsObjectMetaBody,
			expectError:  nil,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given, MarshalJSONAPI(tc.givenJSONAPI))
			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Nil(t, actual)
				return
			}
			is.MustNoError(t, err)
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalLinks(t *testing.T) {
	t.Parallel()

	docLinksArticleABody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"links":{"self":"https://example.com/articles/1"}}`

	tests := []struct {
		description string
		given       any
		givenLink   *Link
		expect      string
	}{
		{
			description: "with link",
			given:       &articleA,
			givenLink:   &Link{Self: fmt.Sprintf("https://example.com/articles/%s", articleA.ID)},
			expect:      docLinksArticleABody,
		}, {
			description: "without link",
			given:       &articleA,
			givenLink:   nil,
			expect:      articleABody,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given, MarshalLinks(tc.givenLink))
			is.MustNoError(t, err) // resource object errors covered in TestMarshal
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalFields(t *testing.T) {
	t.Parallel()

	articleCompleteNoInfoBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A","subtitle":"AA"}}}`

	tests := []struct {
		description string
		setQuery    func(q url.Values)
		expect      string
	}{
		{
			description: "no fields filter",
			setQuery:    func(q url.Values) {},
			expect:      articleCompleteBody,
		}, {
			description: "single fields filter (invalid field)",
			setQuery: func(q url.Values) {
				q.Set("fields[articles]", "not-a-field-in-articles")
			},
			expect: articleOmitTitleFullBody,
		}, {
			description: "single field filter",
			setQuery: func(q url.Values) {
				q.Set("fields[articles]", "title")
			},
			expect: articleABody,
		}, {
			description: "multiple fields filter",
			setQuery: func(q url.Values) {
				q.Set("fields[articles]", "title,subtitle")
			},
			expect: articleCompleteNoInfoBody,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			u, err := url.Parse("http://example.com")
			is.MustNoError(t, err)

			q := u.Query()
			tc.setQuery(q)
			u.RawQuery = q.Encode()

			actual, err := Marshal(&articleComplete, MarshalFields(u.Query()))
			is.MustNoError(t, err)

			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalRelationships(t *testing.T) {
	t.Parallel()

	// filter relationship include data
	commentArchivedQuery := url.Values{}
	commentArchivedQuery.Set("fields[articles]", "title,comments")
	commentArchivedQuery.Set("fields[comments]", "archived")
	articleCommentsArchivedFilterBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},"included":[{"id":"1","type":"comments","attributes":{"archived":true}}]}`

	// articles filter is missing the relationship name
	noCommentsQuery := url.Values{}
	noCommentsQuery.Set("fields[articles]", "title")
	noCommentsQuery.Set("fields[comments]", "archived")
	articleNoCommentsFilterBody := `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"included":[{"id":"1","type":"comments","attributes":{"archived":true}}]}`

	tests := []struct {
		description    string
		given          any
		marshalOptions []MarshalOption
		expect         string
		expectError    error
	}{
		{
			description:    "empty relationships with omitempty",
			given:          &articleRelated,
			marshalOptions: nil,
			expect:         articleABody,
		}, {
			description:    "empty relationships without omitempty",
			given:          &articleRelatedNoOmitEmpty,
			marshalOptions: nil,
			expect:         articleRelatedNoOmitEmptyBody,
			expectError:    nil,
		}, {
			description:    "with related author",
			given:          &articleRelatedAuthor,
			marshalOptions: nil,
			expect:         articleRelatedAuthorBody,
			expectError:    nil,
		}, {
			description:    "with related author (+meta)",
			given:          &articleRelatedAuthorWithMeta,
			marshalOptions: nil,
			expect:         articleRelatedAuthorWithMetaBody,
			expectError:    nil,
		}, {
			description:    "with related comments",
			given:          &articleRelatedComments,
			marshalOptions: nil,
			expect:         articleRelatedCommentsBody,
			expectError:    nil,
		}, {
			description:    "with related author and comments",
			given:          &articleRelatedComplete,
			marshalOptions: nil,
			expect:         articleRelatedCompleteBody,
		}, {
			description:    "with related comments and included comment",
			given:          &articleRelatedComments,
			marshalOptions: []MarshalOption{MarshalInclude(&commentAWithAuthor)}, // TODO: just commentA?
			expect:         articleRelatedCommentsWithIncludeBody,
			expectError:    nil,
		}, {
			description:    "with related comments, included comment, and included author related to comment",
			given:          &articleRelatedComments,
			marshalOptions: []MarshalOption{MarshalInclude(&commentAWithAuthor, &authorA)},
			expect:         articleRelatedCommentsNestedWithIncludeBody,
			expectError:    nil,
		}, {
			description:    "with related comments, included comment, and included fields archive/comment",
			given:          &articleRelatedCommentsArchived,
			marshalOptions: []MarshalOption{MarshalInclude(&commentArchived), MarshalFields(commentArchivedQuery)},
			expect:         articleCommentsArchivedFilterBody,
			expectError:    nil,
		}, {
			description:    "with related comments, included comment, and included filed archive",
			given:          &articleRelatedCommentsArchived,
			marshalOptions: []MarshalOption{MarshalInclude(&commentArchived), MarshalFields(noCommentsQuery)},
			expect:         articleNoCommentsFilterBody,
			expectError:    nil,
		}, {
			description:    "with included comment, included author, and no relationship (partial linkage)",
			given:          &articleA,
			marshalOptions: []MarshalOption{MarshalInclude(&commentAWithAuthor, &authorA)},
			expect:         "",
			expectError:    &PartialLinkageError{[]string{"{Type: comments, ID: 1}", "{Type: author, ID: 1}"}},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given, tc.marshalOptions...)

			if tc.expectError != nil {
				is.EqualError(t, tc.expectError, err)
				is.Nil(t, actual)
				return
			}
			is.MustNoError(t, err) // resource objects covered in TestMarshal
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}

func TestMarshalClientMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		given       any
		expect      string
	}{
		{
			description: "empty primary field",
			given:       &Article{ID: "", Title: "A"},
			expect:      articleANoIDBody,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			actual, err := Marshal(tc.given, MarshalClientMode())
			is.MustNoError(t, err)
			is.EqualJSON(t, tc.expect, string(actual))
		})
	}
}
