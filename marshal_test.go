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
			expect:      nullDataBody,
			expectError: nil,
		}, {
			description: "Article (empty)",
			given:       Article{},
			expect:      nullDataBody,
			expectError: nil,
		}, {
			description: "*Article (empty)",
			given:       new(Article),
			expect:      "",
			expectError: ErrEmptyPrimaryField,
		}, {
			description: "[]*Article (nil)",
			given:       []*Article(nil),
			expect:      emptyManyBody,
			expectError: nil,
		}, {
			description: "[]*Article (empty)",
			given:       make([]*Article, 0),
			expect:      emptyManyBody,
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
			description: "invalid Link.Self.Meta",
			given:       &articleLinkedInvalidSelfMeta,
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
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
			description: "*ArticleEncodingIntID (encoding.TextMarshaler)",
			given:       &articleAEncodingIntID,
			expect:      articleABody,
			expectError: nil,
		}, {
			description: "[]*ArticleEncodinfIntID (encoding.TextMarshaler)",
			given:       &articlesEncodingIntIDABPtr,
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
			description: "ArticleEmbedded",
			given:       &articleEmbedded,
			expect:      articleEmbeddedBody,
			expectError: nil,
		}, {
			description: "ArticleEmbeddedPointer",
			given:       &articleEmbeddedPointer,
			expect:      articleEmbeddedBody,
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
			description: "Error with invalid meta",
			given:       errorsWithInvalidMeta,
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
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
			description: "Error with invalid Links.About.Meta",
			given:       errorsWithInvalidLinkMeta,
			expect:      "",
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
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

	errorsObjectToplevelMetaBody := `{"meta":{"foo":"bar"},"errors":[{"title":"T"}]}`

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
			expect:      articleAToplevelMetaBody,
			expectError: nil,
		}, {
			description: "struct",
			given:       &articleA,
			givenMeta: &struct {
				Foo string `json:"foo"`
			}{Foo: "bar"},
			expect:      articleAToplevelMetaBody,
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
			expect:      errorsObjectToplevelMetaBody,
			expectError: nil,
		}, {
			description: "map[string]any with nil body",
			given:       nil,
			givenMeta:   map[string]any{"foo": "bar"},
			expect:      articleNullWithToplevelMetaBody,
			expectError: nil,
		}, {
			description: "map[string]any with Article (empty)",
			given:       Article{},
			givenMeta:   map[string]any{"foo": "bar"},
			expect:      articleNullWithToplevelMetaBody,
			expectError: nil,
		}, {
			description: "map[string]any with body []*Article (empty)",
			given:       []*Article(nil),
			givenMeta:   map[string]any{"foo": "bar"},
			expect:      articleEmptyArrayWithToplevelMetaBody,
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
			setQuery:    func(_ url.Values) {},
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
		}, {
			description:    "multiple complex with included authors and comments",
			given:          &articlesRelatedComplex,
			marshalOptions: articlesRelatedComplexMarshalOptions,
			expect:         articlesRelatedComplexBody,
			expectError:    nil,
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

	articleNoIDRelatedComplexBody := `{"data":{"type":"articles","attributes":{"title":"Bazel 101"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/articles//relationships/author","related":"http://example.com/articles//author"}},"comments":{"data":[{"type":"comments"},{"type":"comments"},{"type":"comments"}],"links":{"self":"http://example.com/articles//relationships/comments","related":"http://example.com/articles//comments"}}}}}`

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
		{
			description: "empty primary field with relationships",
			given: &ArticleRelated{
				Title:  "Bazel 101",
				Author: &authorA,
				Comments: []*Comment{
					{Body: "Why is Bazel so slow on my computerr?", Archived: true, Author: &authorBWithMeta},
					{Body: "Why is Bazel so slow on my computer?", Author: &authorBWithMeta},
					{Body: "Just use an Apple M1", Author: &authorA},
				},
			},
			expect: articleNoIDRelatedComplexBody,
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

// TestMarshalMemberNameValidation collects tests which verify that invalid member names are caught
// during marshaling, no matter where they're placed. This test does not exhaustively test every
// possible invalid name.
func TestMarshalMemberNameValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description       string
		given             any
		expectError       error
		additionalOptions []MarshalOption
	}{
		{
			description: "Article with valid member names",
			given:       &articleA,
			expectError: nil,
		}, {
			description: "Author with invalid type name",
			given:       &authorWithInvalidTypeName,
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Author with invalid attribute name",
			given:       &authorWithInvalidAttributeName,
			expectError: &MemberNameValidationError{"na%me"},
		}, {
			description: "Article with invalid resource meta member name",
			given:       &articleWithInvalidResourceMetaMemberName,
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description:       "Article with invalid top-level meta member name",
			given:             &articleA,
			expectError:       &MemberNameValidationError{"foo%"},
			additionalOptions: []MarshalOption{MarshalMeta(map[string]any{"foo%": 2})},
		}, {
			description:       "Article with invalid jsonapi meta member name",
			given:             &articleA,
			expectError:       &MemberNameValidationError{"foo%"},
			additionalOptions: []MarshalOption{MarshalJSONAPI(map[string]any{"foo%": 1})},
		}, {
			description: "Article with invalid link meta member name",
			given:       &articleWithInvalidLinkMetaMemberName,
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Article with invalid relationship name",
			given:       &articleWithInvalidRelationshipName,
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Article with invalid relationship type name",
			given:       &articleWithInvalidRelationshipTypeName,
			expectError: &MemberNameValidationError{"aut%hor"},
		}, {
			description: "Article with invalid relationship attribute name not included",
			given:       &articleWithInvalidRelationshipAttributeName,
			expectError: nil,
		}, {
			description:       "Article with invalid relationship attribute name included",
			given:             &articleWithInvalidRelationshipAttributeName,
			expectError:       &MemberNameValidationError{"na%me"},
			additionalOptions: []MarshalOption{MarshalInclude(&authorWithInvalidAttributeName)},
		}, {
			description: "Articles with one invalid resource meta member name",
			given: []*ArticleWithGenericMeta{
				{ID: "1"}, {ID: "1", Meta: map[string]any{"foo%": 1}},
			},
			expectError: &MemberNameValidationError{"foo%"},
		}, {
			description: "Website with invalid nested relationship type name",
			given:       &websiteWithInvalidNestedRelationshipTypeName,
			expectError: &MemberNameValidationError{"aut%hor"},
			additionalOptions: []MarshalOption{
				MarshalInclude(
					websiteWithInvalidNestedRelationshipTypeName.Articles[0],
					websiteWithInvalidNestedRelationshipTypeName.Articles[1],
					websiteWithInvalidNestedRelationshipTypeName.Articles[1].Author,
				),
			},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			opts := tc.additionalOptions
			_, err := Marshal(tc.given, opts...)
			is.EqualError(t, tc.expectError, err)

			opts = append(opts, MarshalSetNameValidation(DisableValidation))
			_, err = Marshal(tc.given, opts...)
			is.MustNoError(t, err)
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	benchmarks := []struct {
		name  string
		given any
		opts  []MarshalOption
	}{
		{
			name:  "ArticleSimple",
			given: articleA,
			opts:  nil,
		}, {
			name:  "ArticleSimpleWithToplevelMeta",
			given: articleA,
			opts: []MarshalOption{
				MarshalMeta(map[string]any{"foo": "bar"}),
			},
		}, {
			name:  "ArticleComplex",
			given: articleRelatedComments,
			opts: []MarshalOption{
				MarshalInclude(&commentAWithAuthor, &authorA),
			},
		}, {
			name:  "ArticleComplexDisableNameValidation",
			given: articleRelatedComments,
			opts: []MarshalOption{
				MarshalInclude(&commentAWithAuthor, &authorA),
				MarshalSetNameValidation(DisableValidation),
			},
		}, {
			name:  "ArticlesComplex",
			given: articlesRelatedComplex,
			opts:  articlesRelatedComplexMarshalOptions,
		}, {
			name:  "ArticlesComplexDisableNameValidation",
			given: articlesRelatedComplex,
			opts: append(
				[]MarshalOption{MarshalSetNameValidation(DisableValidation)},
				articlesRelatedComplexMarshalOptions...,
			),
		},
	}

	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_, _ = Marshal(bm.given, bm.opts...)
			}
		})
	}
}
