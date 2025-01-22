package jsonapi

import (
	"encoding"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	// authors
	authorA         = Author{ID: "1", Name: "A"}
	authorAWithMeta = Author{ID: "1", Name: "A", Meta: map[string]any{"count": 10}}
	authorBWithMeta = Author{ID: "2", Name: "B", Meta: map[string]any{"count": 10.0}}

	// comments
	commentA              = Comment{ID: "1", Body: "A"}
	commentB              = Comment{ID: "2", Body: "B"}
	commentAWithAuthor    = Comment{ID: "1", Body: "A", Author: &authorA}
	commentArchived       = Comment{ID: "1", Body: "A", Archived: true}
	commentsAB            = []*Comment{&commentA, &commentB}
	commentEmbeddedFields = CommentFields{Body: "A", Author: Author{ID: "1"}}
	commentEmbedded       = CommentEmbedded{ID: "1", CommentFields: commentEmbeddedFields}

	commentEmbeddedFieldsPointer = CommentFieldsPointer{Body: "A", Author: &Author{ID: "1"}}
	commentEmbeddedPointer       = CommentEmbeddedPointer{ID: "1", CommentFieldsPointer: commentEmbeddedFieldsPointer}

	// articles
	articleA        = Article{ID: "1", Title: "A"}
	articleANoID    = Article{Title: "A"}
	articleB        = Article{ID: "2", Title: "B"}
	articlesAB      = []Article{articleA, articleB}
	articlesAA      = []Article{articleA, articleA}
	articlesABPtr   = []*Article{&articleA, &articleB}
	articleComplete = ArticleComplete{
		ID:       "1",
		Title:    "A",
		SubTitle: "AA",
		Info: &ArticleInfo{
			PublishDate: time.Date(1989, 06, 15, 0, 0, 0, 0, time.UTC),
			Tags:        []string{"a", "b"},
			IsPublic:    true,
			Metrics: &ArticleMetrics{
				Views: 10,
				Reads: 4,
			},
		},
	}
	articleALinked                               = ArticleLinked{ID: "1", Title: "A"}
	articleLinkedOnlySelf                        = ArticleLinkedOnlySelf{ID: "1"}
	articleLinkedInvalidSelf                     = ArticleLinkedInvalidSelf{ID: "1"}
	articleLinkedInvalidRelated                  = ArticleLinkedInvalidRelated{ID: "1"}
	articleLinkedInvalidMissingFields            = ArticleLinkedInvalidMissingFields{ID: "1"}
	articleLinkedInvalidMissingFieldsEmptyValues = ArticleLinkedInvalidMissingFieldsEmptyValues{ID: "1"}
	articleLinkedInvalidSelfMeta                 = ArticleLinkedInvalidSelfMeta{ID: "1"}
	articleOmitTitleFull                         = ArticleOmitTitle{ID: "1"}
	articleOmitTitlePartial                      = ArticleOmitTitle{ID: "1", Subtitle: "A"}
	articleAIntID                                = ArticleIntID{ID: 1, Title: "A"}
	articleBIntID                                = ArticleIntID{ID: 2, Title: "B"}
	articlesIntIDABPtr                           = []*ArticleIntID{&articleAIntID, &articleBIntID}
	articleAIntIDID                              = ArticleIntIDID{ID: IntID(1), Title: "A"}
	articleBIntIDID                              = ArticleIntIDID{ID: IntID(2), Title: "B"}
	articlesIntIDIDABPtr                         = []*ArticleIntIDID{&articleAIntIDID, &articleBIntIDID}
	articleAEncodingIntID                        = ArticleEncodingIntID{ID: EncodingIntID(1), Title: "A"}
	articleBEncodingIntID                        = ArticleEncodingIntID{ID: EncodingIntID(2), Title: "B"}
	articlesEncodingIntIDABPtr                   = []*ArticleEncodingIntID{&articleAEncodingIntID, &articleBEncodingIntID}
	articleEmbedded                              = ArticleEmbedded{ID: "1", Title: "A", Metadata: Metadata{LastModified: time.Date(1989, 06, 15, 0, 0, 0, 0, time.UTC)}}
	articleEmbeddedPointer                       = ArticleEmbeddedPointer{ID: "1", Title: "A", Metadata: &Metadata{LastModified: time.Date(1989, 06, 15, 0, 0, 0, 0, time.UTC)}}

	// articles with optional meta
	articleAWithMeta              = ArticleWithMeta{ID: "1", Title: "A", Meta: &ArticleMetrics{Views: 10, Reads: 4}}
	articleWithResourceObjectMeta = ArticleWithResourceObjectMeta{
		ID:    "1",
		Title: "A",
		Meta:  map[string]any{"count": 10},
	}
	articleWithoutResourceObjectMeta = ArticleWithResourceObjectMeta{ID: "1", Title: "A"}

	// articles with relationships
	articleRelated                  = ArticleRelated{ID: "1", Title: "A"}
	articleRelatedNoOmitEmpty       = ArticleRelatedNoOmitEmpty{ID: "1", Title: "A"}
	articleRelatedAuthor            = ArticleRelated{ID: "1", Title: "A", Author: &authorA}
	articleRelatedAuthorWithMeta    = ArticleRelated{ID: "1", Title: "A", Author: &authorAWithMeta}
	articleRelatedComments          = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentA}}
	articleRelatedNonuniqueComments = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentA, &commentA}}
	articleRelatedCommentsArchived  = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentArchived}}
	articleRelatedCommentsNested    = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentAWithAuthor}}
	articleRelatedComplete          = ArticleRelated{ID: "1", Title: "A", Author: &authorAWithMeta, Comments: commentsAB}
	articlesRelatedComplex          = []*ArticleRelated{
		{
			ID:     "1",
			Title:  "Bazel 101",
			Author: &authorA,
			Comments: []*Comment{
				{ID: "11", Body: "Why is Bazel so slow on my computerr?", Archived: true, Author: &authorBWithMeta},
				{ID: "12", Body: "Why is Bazel so slow on my computer?", Author: &authorBWithMeta},
				{ID: "13", Body: "Just use an Apple M1", Author: &authorA},
			},
		}, {
			ID:     "2",
			Title:  "Why Should I Use JSON:API?",
			Author: &authorBWithMeta,
			Comments: []*Comment{
				{ID: "21", Body: "I wish they changed the name...", Author: &authorBWithMeta},
			},
		}, {
			ID:    "3",
			Title: "Internal Test Article Created In Production For Some Reason",
			Comments: []*Comment{
				{ID: "31", Body: "test1"},
				{ID: "32", Body: "test2"},
			},
		}, {
			ID:     "4",
			Title:  "How to Rewrite Everything in Rust",
			Author: &authorA,
		},
	}
	articlesRelatedComplexMarshalOptions = []MarshalOption{
		MarshalJSONAPI(map[string]any{"meta_kind": "jsonapi meta"}),
		MarshalMeta(map[string]any{"meta_kind": "document-level meta"}),
		MarshalInclude(
			authorA,
			authorBWithMeta,
			articlesRelatedComplex[0].Comments[0],
			articlesRelatedComplex[0].Comments[1],
			articlesRelatedComplex[0].Comments[2],
			articlesRelatedComplex[1].Comments[0],
			articlesRelatedComplex[2].Comments[0],
			articlesRelatedComplex[2].Comments[1],
		),
	}

	// articles with non-spec-conforming member names
	authorWithInvalidTypeName                    = AuthorWithInvalidTypeName{ID: "1"}
	authorWithInvalidAttributeName               = AuthorWithInvalidAttributeName{ID: "1", Name: "A"}
	articleWithInvalidResourceMetaMemberName     = ArticleWithGenericMeta{ID: "1", Meta: map[string]any{"foo%": 2}}
	articleWithInvalidLinkMetaMemberName         = ArticleWithInvalidLinkMetaMemberName{ID: "1"}
	articleWithInvalidRelationshipName           = ArticleWithInvalidRelationshipName{ID: "1", Author: &authorA}
	articleWithInvalidRelationshipTypeName       = ArticleWithInvalidRelationshipTypeName{ID: "1", Author: &authorWithInvalidTypeName}
	articleWithInvalidRelationshipAttributeName  = ArticleWithInvalidRelationshipAttributeName{ID: "1", Author: &authorWithInvalidAttributeName}
	websiteWithInvalidNestedRelationshipTypeName = WebsiteWithInvalidNestedRelationshipTypeName{
		ID:       "1",
		Articles: []*ArticleWithInvalidRelationshipTypeName{{ID: "2"}, &articleWithInvalidRelationshipTypeName},
	}

	// article bodies
	nullDataBody                          = `{"data":null}`
	emptySingleBody                       = `{"data":{}}`
	emptyManyBody                         = `{"data":[]}`
	articleABody                          = `{"data":{"type":"articles","id":"1","attributes":{"title":"A"}}}`
	articleANoIDBody                      = `{"data":{"type":"articles","attributes":{"title":"A"}}}`
	articleAInvalidTypeBody               = `{"data":{"type":"not-articles","id":"1","attributes":{"title":"A"}}}`
	articleOmitTitleFullBody              = `{"data":{"type":"articles","id":"1"}}`
	articleOmitTitlePartialBody           = `{"data":{"type":"articles","id":"1","attributes":{"subtitle":"A"}}}`
	articlesABBody                        = `{"data":[{"type":"articles","id":"1","attributes":{"title":"A"}},{"type":"articles","id":"2","attributes":{"title":"B"}}]}`
	articlesABNonuniqueData               = `{"data":[{"type":"articles","id":"1","attributes":{"title":"A"}},{"type":"articles","id":"1","attributes":{"title":"B"}}]}`
	articleCompleteBody                   = `{"data":{"id":"1","type":"articles","attributes":{"info":{"publishDate":"1989-06-15T00:00:00Z","tags":["a","b"],"isPublic":true,"metrics":{"views":10,"reads":4}},"title":"A","subtitle":"AA"}}}`
	articleALinkedBody                    = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"links":{"self":"https://example.com/articles/1","related":{"href":"https://example.com/articles/1/comments","meta":{"count":10}}}}}`
	articleLinkedOnlySelfBody             = `{"data":{"id":"1","type":"articles","links":{"self":"https://example.com/articles/1"}}}`
	articleWithResourceObjectMetaBody     = `{"data":{"type":"articles","id":"1","attributes":{"title":"A"},"meta":{"count":10}}}`
	articleAToplevelMetaBody              = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"meta":{"foo":"bar"}}`
	articleAWithMetaBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"meta":{"views":10,"reads":4}}}`
	articleNullWithToplevelMetaBody       = `{"data":null,"meta":{"foo":"bar"}}`
	articleEmptyArrayWithToplevelMetaBody = `{"data":[],"meta":{"foo":"bar"}}`
	articleEmbeddedBody                   = `{"data":{"type":"articles","id":"1","attributes":{"title":"A","lastModified":"1989-06-15T00:00:00Z"}}}`
	commentEmbeddedBody                   = `{"data":{"id":"1","type":"comments","attributes":{"body":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}}}}}`
	articleAWithTopLevelLink              = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"links":{"self":"http://example.com/article/1"}}`
	articleAWithTopLevelLinksRelated      = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"links":{"related":{"href":"http://example.com/article/1/comments","meta":{"foo":"bar"}}}}`

	// articles with relationships bodies
	articleRelatedInvalidEmptyRelationshipBody  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{}}}}`
	articleRelatedInvalidEmptyDataBody          = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{}}}}}`
	articleRelatedNoOmitEmptyBody               = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":null},"comments":{"data":[]}}}}`
	articleRelatedLiteralNullData               = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":null},"comments":{"data":null}}}}`
	articleRelatedAuthorBody                    = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}}}}}`
	articleRelatedAuthorTwiceBody               = `{"data":[{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}}}},{"id":"2","type":"articles","attributes":{"title":"B"},"relationships":{"author":{"data":{"id":"1","type":"author"}}}}]}`
	articleRelatedAuthorTwiceWithIncludeBody    = `{"data":[{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}}}},{"id":"2","type":"articles","attributes":{"title":"B"},"relationships":{"author":{"data":{"id":"1","type":"author"}}}}],"included":[{"id":"1","type":"author","attributes":{"name":"A"}}]}`
	articleRelatedAuthorWithMetaBody            = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}}}}}`
	articleRelatedCommentsBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}}}`
	articleRelatedCommentsWithIncludeBody       = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},"included":[{"id":"1","type":"comments","attributes":{"body":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/comments/1/relationships/author","related":"http://example.com/comments/1/author"}}}}]}`
	articleRelatedCompleteBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}},"comments":{"data":[{"id":"1","type":"comments"},{"id":"2","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}}}`
	articleRelatedCompleteWithIncludeBody       = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}},"comments":{"data":[{"id":"1","type":"comments"},{"id":"2","type":"comments"}]}}},"included":[{"id":"1","type":"author","attributes":{"name":"A"}},{"id":"1","type":"comments","attributes":{"body":"A"}},{"id":"2","type":"comments","attributes":{"body":"B"}}]}`
	articleRelatedNonuniqueLinkage              = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}},"comments":{"data":[{"id":"1","type":"comments"},{"id":"1","type":"comments"}]}}}}`
	articleRelatedDifferentRefToSameObject      = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"}},"contact":{"data":{"id":"1","type":"author"}},"comments":{"data":[{"id":"1","type":"comments"}]}}}}`
	articleRelatedCommentsNestedWithIncludeBody = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},"included":[{"id":"1","type":"comments","attributes":{"body":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/comments/1/relationships/author","related":"http://example.com/comments/1/author"}}}},{"id":"1","type":"author","attributes":{"name":"A"}}]}`
	articleWithIncludeOnlyBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"included":[{"id":"1","type":"author","attributes":{"name":"A"}}]}`
	articleRelatedAuthorLinksOnlyBody           = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}}}}}`
	articleRelatedAuthorMetaOnlyBody            = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"meta":{"foo":"bar"}}}}}`
	articlesRelatedComplexBody                  = `{"data":[{"id":"1","type":"articles","attributes":{"title":"Bazel 101"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}},"comments":{"data":[{"id":"11","type":"comments"},{"id":"12","type":"comments"},{"id":"13","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},{"id":"2","type":"articles","attributes":{"title":"Why Should I Use JSON:API?"},"relationships":{"author":{"data":{"id":"2","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/articles/2/relationships/author","related":"http://example.com/articles/2/author"}},"comments":{"data":[{"id":"21","type":"comments"}],"links":{"self":"http://example.com/articles/2/relationships/comments","related":"http://example.com/articles/2/comments"}}}},{"id":"3","type":"articles","attributes":{"title":"Internal Test Article Created In Production For Some Reason"},"relationships":{"comments":{"data":[{"id":"31","type":"comments"},{"id":"32","type":"comments"}],"links":{"self":"http://example.com/articles/3/relationships/comments","related":"http://example.com/articles/3/comments"}}}},{"id":"4","type":"articles","attributes":{"title":"How to Rewrite Everything in Rust"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/articles/4/relationships/author","related":"http://example.com/articles/4/author"}}}}],"meta":{"meta_kind":"document-level meta"},"jsonapi":{"version":"1.0","meta":{"meta_kind":"jsonapi meta"}},"included":[{"id":"1","type":"author","attributes":{"name":"A"}},{"id":"2","type":"author","attributes":{"name":"B"},"meta":{"count":10}},{"id":"11","type":"comments","attributes":{"archived":true,"body":"Why is Bazel so slow on my computerr?"},"relationships":{"author":{"data":{"id":"2","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/comments/11/relationships/author","related":"http://example.com/comments/11/author"}}}},{"id":"12","type":"comments","attributes":{"body":"Why is Bazel so slow on my computer?"},"relationships":{"author":{"data":{"id":"2","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/comments/12/relationships/author","related":"http://example.com/comments/12/author"}}}},{"id":"13","type":"comments","attributes":{"body":"Just use an Apple M1"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/comments/13/relationships/author","related":"http://example.com/comments/13/author"}}}},{"id":"21","type":"comments","attributes":{"body":"I wish they changed the name..."},"relationships":{"author":{"data":{"id":"2","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/comments/21/relationships/author","related":"http://example.com/comments/21/author"}}}},{"id":"31","type":"comments","attributes":{"body":"test1"}},{"id":"32","type":"comments","attributes":{"body":"test2"}}]}`

	// articles with non-conforming member name bodies
	authorWithInvalidTypeNameBody                           = `{"data":{"id":"1","type":"aut%hor"}}`
	authorWithInvalidAttributeNameBody                      = `{"data":{"id":"1","type":"author","attributes":{"na%me":"A"}}}`
	articleWithInvalidResourceMetaMemberNameBody            = `{"data":{"id":"1","type":"articles","meta":{"foo%":2}}}`
	articleWithInvalidLinkMetaMemberNameBody                = `{"data":{"id":"1","type":"articles","links":{"self":{"href":"foo","meta":{"foo%":1}}}}}`
	articleWithInvalidRelationshipNameBody                  = `{"data":{"id":"1","type":"articles","relationships":{"aut%hor":{"data":{"id":"1","type":"author"}}}}}`
	articleWithInvalidRelationshipTypeNameBody              = `{"data":{"id":"1","type":"articles","relationships":{"author":{"data":{"id":"1","type":"aut%hor"}}}}}`
	articleWithInvalidRelationshipAttributeNameIncludedBody = `{"data":{"id":"1","type":"articles","relationships":{"author":{"data":{"id":"1","type":"author"}}}},"included":[{"id":"1","type":"author","attributes":{"na%me":"A"}}]}`
	websiteWithInvalidNestedRelationshipTypeNameBody        = `{"data":{"id":"1","type":"website","relationships":{"articles":{"data":[{"id":"2","type":"articles"},{"id":"1","type":"articles"}]}}},"included":[{"id":"2","type":"articles","relationships":{"author":{"data":null}}},{"id":"1","type":"articles","relationships":{"author":{"data":{"id":"1","type":"aut%hor"}}}},{"id":"1","type":"aut%hor"}]}`

	// error structs
	errorsSimpleStruct         = Error{Title: "T"}           //nolint: errname
	errorsSimpleSliceSingle    = []Error{errorsSimpleStruct} //nolint: errname
	errorsSimpleSliceSinglePtr = []*Error{&errorsSimpleStruct}
	errorsComplexStruct        = Error{ //nolint: errname
		ID:     "1",
		Links:  &ErrorLink{About: "A", Type: "TY"},
		Status: Status(http.StatusInternalServerError),
		Code:   "C",
		Title:  "T",
		Detail: "D",
		Source: &ErrorSource{Pointer: "PO", Parameter: "PA", Header: "H"},
		Meta:   map[string]string{"K": "V"},
	}
	errorsComplexSliceMany    = []Error{errorsSimpleStruct, errorsComplexStruct}
	errorsComplexSliceManyPtr = []*Error{&errorsSimpleStruct, &errorsComplexStruct}
	errorsWithInvalidMeta     = Error{ID: "1", Meta: "foo"} //nolint: errname
	errorsWithLinkObject      = Error{                      //nolint: errname
		Links: &ErrorLink{
			About: &LinkObject{
				Href: "A",
				Meta: map[string]any{
					"key_s": "B",
					"key_i": 420,
				},
			},
		},
	}
	errorsWithInvalidLink     = Error{Links: &ErrorLink{About: 1}}                                   //nolint: errname
	errorsWithInvalidLinkMeta = Error{Links: &ErrorLink{About: &LinkObject{Href: "A", Meta: "foo"}}} //nolint: errname

	// error bodies
	errorsSimpleStructBody     = `{"errors":[{"title":"T"}]}`
	errorsComplexStructBody    = `{"errors":[{"id":"1","links":{"about":"A","type":"TY"},"status":"500","code":"C","title":"T","detail":"D","source":{"pointer":"PO","parameter":"PA","header":"H"},"meta":{"K":"V"}}]}`
	errorsComplexSliceManyBody = `{"errors":[{"title":"T"},{"id":"1","links":{"about":"A","type":"TY"},"status":"500","code":"C","title":"T","detail":"D","source":{"pointer":"PO","parameter":"PA","header":"H"},"meta":{"K":"V"}}]}`
	errorsWithLinkObjectBody   = `{"errors":[{"links":{"about":{"href":"A","meta":{"key_i":420,"key_s":"B"}}}}]}`
)

type Article struct {
	ID    string `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`

	// Ignored is included to ensure un-tagged fields are ignored
	Ignored string `json:"ignored"`
}

type ArticleMetrics struct {
	Views int64 `json:"views"`
	Reads int64 `json:"reads"`
}

type ArticleInfo struct {
	PublishDate time.Time       `json:"publishDate"`
	Tags        []string        `json:"tags"`
	IsPublic    bool            `json:"isPublic"`
	Metrics     *ArticleMetrics `json:"metrics"`
}

type ArticleComplete struct {
	ID       string       `jsonapi:"primary,articles"`
	Title    string       `jsonapi:"attribute" json:"title"`
	SubTitle string       `jsonapi:"attribute" json:"subtitle,omitempty"`
	Info     *ArticleInfo `jsonapi:"attribute" json:"info"`

	// Ignored is included to ensure un-tagged fields are ignored
	Ignored string `json:"ignored"`
}

type ArticleWithMeta struct {
	ID    string          `jsonapi:"primary,articles"`
	Title string          `jsonapi:"attribute" json:"title"`
	Meta  *ArticleMetrics `jsonapi:"meta"`
}

type ArticleLinked struct {
	ID    string `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`

	// Ignored is included to ensure un-tagged fields are ignored
	Ignored string `json:"ignored"`
}

func (a *ArticleLinked) Link() *Link {
	return &Link{
		Self: fmt.Sprintf("https://example.com/articles/%s", a.ID),
		Related: &LinkObject{
			Href: fmt.Sprintf("https://example.com/articles/%s/comments", a.ID),
			Meta: map[string]int{"count": 10},
		},
	}
}

type ArticleLinkedOnlySelf struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedOnlySelf) Link() *Link {
	return &Link{Self: fmt.Sprintf("https://example.com/articles/%s", a.ID)}
}

type ArticleLinkedInvalidSelf struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedInvalidSelf) Link() *Link {
	return &Link{Self: 10}
}

type ArticleLinkedInvalidRelated struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedInvalidRelated) Link() *Link {
	return &Link{Related: 10}
}

type ArticleLinkedInvalidMissingFields struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedInvalidMissingFields) Link() *Link {
	return &Link{Self: nil, Related: nil}
}

type ArticleLinkedInvalidMissingFieldsEmptyValues struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedInvalidMissingFieldsEmptyValues) Link() *Link {
	var lo LinkObject
	return &Link{Self: "", Related: &lo}
}

type ArticleLinkedInvalidSelfMeta struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleLinkedInvalidSelfMeta) Link() *Link {
	return &Link{Self: &LinkObject{Href: "foo", Meta: "foo"}}
}

type ArticleOmitTitle struct {
	ID       string `jsonapi:"primary,articles"`
	Title    string `jsonapi:"attribute" json:"title,omitempty"`
	Subtitle string `jsonapi:"attribute" json:"subtitle,omitempty"`
}

type ArticleIntID struct {
	ID    int    `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`
}

func (a *ArticleIntID) MarshalID() string {
	return fmt.Sprintf("%d", a.ID)
}

func (a *ArticleIntID) UnmarshalID(id string) error {
	v, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	a.ID = v
	return nil
}

type IntID int

func (i IntID) String() string {
	return fmt.Sprintf("%d", i)
}

type ArticleIntIDID struct {
	ID    IntID  `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`
}

func (a *ArticleIntIDID) UnmarshalID(id string) error {
	v, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	a.ID = IntID(v)
	return nil
}

var (
	// ensure EncodingIntID implements encoding.[TextMarshaler|TextUnmarshaler]
	_ encoding.TextMarshaler   = (*EncodingIntID)(nil)
	_ encoding.TextUnmarshaler = (*EncodingIntID)(nil)
)

type EncodingIntID int

func (i EncodingIntID) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", i)), nil
}

func (i *EncodingIntID) UnmarshalText(text []byte) error {
	v, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	*i = EncodingIntID(v)
	return nil
}

type ArticleEncodingIntID struct {
	ID    EncodingIntID `jsonapi:"primary,articles"`
	Title string        `jsonapi:"attribute" json:"title"`
}

type ArticleWithResourceObjectMeta struct {
	ID    string         `jsonapi:"primary,articles"`
	Title string         `jsonapi:"attribute" json:"title"`
	Meta  map[string]any `jsonapi:"meta"`
}

type Comment struct {
	ID       string  `jsonapi:"primary,comments"`
	Body     string  `jsonapi:"attribute" json:"body"`
	Archived bool    `jsonapi:"attribute" json:"archived,omitempty"`
	Author   *Author `jsonapi:"relationship" json:"author,omitempty"`
}

func (c *Comment) LinkRelation(relation string) *Link {
	return &Link{
		Self:    fmt.Sprintf("http://example.com/comments/%s/relationships/%s", c.ID, relation),
		Related: fmt.Sprintf("http://example.com/comments/%s/%s", c.ID, relation),
	}
}

type Author struct {
	ID   string         `jsonapi:"primary,author"`
	Name string         `jsonapi:"attribute" json:"name"`
	Meta map[string]any `jsonapi:"meta"`
}

type ArticleRelated struct {
	ID       string     `jsonapi:"primary,articles"`
	Title    string     `jsonapi:"attribute" json:"title"`
	Author   *Author    `jsonapi:"relationship" json:"author,omitempty"`
	Comments []*Comment `jsonapi:"relationship" json:"comments,omitempty"`
}

type ArticleRelatedWithContact struct {
	ID       string     `jsonapi:"primary,articles"`
	Title    string     `jsonapi:"attribute" json:"title"`
	Contact  *Author    `jsonapi:"relationship" json:"contact,omitempty"`
	Author   *Author    `jsonapi:"relationship" json:"author,omitempty"`
	Comments []*Comment `jsonapi:"relationship" json:"comments,omitempty"`
}

func (a *ArticleRelated) LinkRelation(relation string) *Link {
	return &Link{
		Self:    fmt.Sprintf("http://example.com/articles/%s/relationships/%s", a.ID, relation),
		Related: fmt.Sprintf("http://example.com/articles/%s/%s", a.ID, relation),
	}
}

type ArticleRelatedNoOmitEmpty struct {
	ID       string     `jsonapi:"primary,articles"`
	Title    string     `jsonapi:"attribute" json:"title"`
	Author   *Author    `jsonapi:"relationship" json:"author"`
	Comments []*Comment `jsonapi:"relationship" json:"comments"`
}

type ArticleDoubleID struct {
	ID      string `jsonapi:"primary,articles"`
	Title   string `jsonapi:"attribute" json:"title"`
	OtherID string `jsonapi:"primary,article"`
}

type Metadata struct {
	LastModified time.Time `jsonapi:"attribute" json:"lastModified"`
}

type CommentFields struct {
	Body     string `jsonapi:"attribute" json:"body"`
	Archived bool   `jsonapi:"attribute" json:"archived,omitempty"`
	Author   Author `jsonapi:"relationship" json:"author,omitempty"`
}

type CommentFieldsPointer struct {
	Body     string  `jsonapi:"attribute" json:"body"`
	Archived bool    `jsonapi:"attribute" json:"archived,omitempty"`
	Author   *Author `jsonapi:"relationship" json:"author,omitempty"`
}

type CommentEmbedded struct {
	ID string `jsonapi:"primary,comments"`
	CommentFields
}

type CommentEmbeddedPointer struct {
	ID string `jsonapi:"primary,comments"`
	CommentFieldsPointer
}

type ArticleEmbedded struct {
	Metadata

	ID    string `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`
}

type ArticleEmbeddedPointer struct {
	*Metadata

	ID    string `jsonapi:"primary,articles"`
	Title string `jsonapi:"attribute" json:"title"`
}

type ArticleWithGenericMeta struct {
	ID   string         `jsonapi:"primary,articles"`
	Meta map[string]any `jsonapi:"meta"`
}

type AuthorWithInvalidTypeName struct {
	ID string `jsonapi:"primary,aut%hor"`
}

type AuthorWithInvalidAttributeName struct {
	ID   string `jsonapi:"primary,author"`
	Name string `jsonapi:"attribute" json:"na%me"`
}

type ArticleWithInvalidLinkMetaMemberName struct {
	ID string `jsonapi:"primary,articles"`
}

func (a *ArticleWithInvalidLinkMetaMemberName) Link() *Link {
	return &Link{
		Self: &LinkObject{
			Href: "foo",
			Meta: map[string]int{"foo%": 1},
		},
	}
}

type ArticleWithInvalidRelationshipName struct {
	ID     string  `jsonapi:"primary,articles"`
	Author *Author `jsonapi:"relationship" json:"aut%hor"`
}

type ArticleWithInvalidRelationshipTypeName struct {
	ID     string                     `jsonapi:"primary,articles"`
	Author *AuthorWithInvalidTypeName `jsonapi:"relationship" json:"author"`
}

type ArticleWithInvalidRelationshipAttributeName struct {
	ID     string                          `jsonapi:"primary,articles"`
	Author *AuthorWithInvalidAttributeName `jsonapi:"relationship" json:"author"`
}

type WebsiteWithInvalidNestedRelationshipTypeName struct {
	ID       string                                    `jsonapi:"primary,website"`
	Articles []*ArticleWithInvalidRelationshipTypeName `jsonapi:"relationship" json:"articles"`
}
