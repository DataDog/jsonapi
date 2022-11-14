package jsonapi

import (
	"fmt"
	"strconv"
	"time"
)

var (
	// authors
	authorA         = Author{ID: "1", Name: "A"}
	authorAWithMeta = Author{ID: "1", Name: "A", Meta: map[string]any{"count": 10}}

	// comments
	commentA           = Comment{ID: "1", Body: "A"}
	commentB           = Comment{ID: "2", Body: "B"}
	commentAWithAuthor = Comment{ID: "1", Body: "A", Author: &authorA}
	commentArchived    = Comment{ID: "1", Body: "A", Archived: true}
	commentsAB         = []*Comment{&commentA, &commentB}

	// articles
	articleA        = Article{ID: "1", Title: "A"}
	articleANoID    = Article{Title: "A"}
	articleB        = Article{ID: "2", Title: "B"}
	articlesAB      = []Article{articleA, articleB}
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
	articleOmitTitleFull                         = ArticleOmitTitle{ID: "1"}
	articleOmitTitlePartial                      = ArticleOmitTitle{ID: "1", Subtitle: "A"}
	articleAIntID                                = ArticleIntID{ID: 1, Title: "A"}
	articleBIntID                                = ArticleIntID{ID: 2, Title: "B"}
	articlesIntIDABPtr                           = []*ArticleIntID{&articleAIntID, &articleBIntID}
	articleAIntIDID                              = ArticleIntIDID{ID: IntID(1), Title: "A"}
	articleBIntIDID                              = ArticleIntIDID{ID: IntID(2), Title: "B"}
	articlesIntIDIDABPtr                         = []*ArticleIntIDID{&articleAIntIDID, &articleBIntIDID}

	// articles with optional meta
	articleAWithMeta              = ArticleWithMeta{ID: "1", Title: "A", Meta: &ArticleMetrics{Views: 10, Reads: 4}}
	articleWithResourceObjectMeta = ArticleWithResourceObjectMeta{
		ID:    "1",
		Title: "A",
		Meta:  map[string]any{"count": 10},
	}
	articleWithoutResourceObjectMeta = ArticleWithResourceObjectMeta{ID: "1", Title: "A"}

	// articles with relationships
	articleRelated                 = ArticleRelated{ID: "1", Title: "A"}
	articleRelatedNoOmitEmpty      = ArticleRelatedNoOmitEmpty{ID: "1", Title: "A"}
	articleRelatedAuthor           = ArticleRelated{ID: "1", Title: "A", Author: &authorA}
	articleRelatedAuthorWithMeta   = ArticleRelated{ID: "1", Title: "A", Author: &authorAWithMeta}
	articleRelatedComments         = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentA}}
	articleRelatedCommentsArchived = ArticleRelated{ID: "1", Title: "A", Comments: []*Comment{&commentArchived}}
	articleRelatedComplete         = ArticleRelated{ID: "1", Title: "A", Author: &authorAWithMeta, Comments: commentsAB}

	// article bodies
	emptyBody                         = `{"data":[]}`
	articleABody                      = `{"data":{"type":"articles","id":"1","attributes":{"title":"A"}}}`
	articleANoIDBody                  = `{"data":{"type":"articles","attributes":{"title":"A"}}}`
	articleAInvalidTypeBody           = `{"data":{"type":"not-articles","id":"1","attributes":{"title":"A"}}}`
	articleOmitTitleFullBody          = `{"data":{"type":"articles","id":"1"}}`
	articleOmitTitlePartialBody       = `{"data":{"type":"articles","id":"1","attributes":{"subtitle":"A"}}}`
	articlesABBody                    = `{"data":[{"type":"articles","id":"1","attributes":{"title":"A"}},{"type":"articles","id":"2","attributes":{"title":"B"}}]}`
	articleCompleteBody               = `{"data":{"id":"1","type":"articles","attributes":{"info":{"publishDate":"1989-06-15T00:00:00Z","tags":["a","b"],"isPublic":true,"metrics":{"views":10,"reads":4}},"title":"A","subtitle":"AA"}}}`
	articleALinkedBody                = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"links":{"self":"https://example.com/articles/1","related":{"href":"https://example.com/articles/1/comments","meta":{"count":10}}}}}`
	articleLinkedOnlySelfBody         = `{"data":{"id":"1","type":"articles","links":{"self":"https://example.com/articles/1"}}}`
	articleWithResourceObjectMetaBody = `{"data":{"type":"articles","id":"1","attributes":{"title":"A"},"meta":{"count":10}}}`
	articleAWithMetaBody              = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"meta":{"views":10,"reads":4}}}`

	// articles with relationships bodies
	articleRelatedNoOmitEmptyBody               = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":null},"comments":{"data":[]}}}}`
	articleRelatedAuthorBody                    = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}}}}}`
	articleRelatedAuthorWithMetaBody            = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}}}}}`
	articleRelatedCommentsBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}}}`
	articleRelatedCommentsWithIncludeBody       = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},"included":[{"id":"1","type":"comments","attributes":{"body":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/comments/1/relationships/author","related":"http://example.com/comments/1/author"}}}}]}`
	articleRelatedCompleteBody                  = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"meta":{"count":10},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}},"comments":{"data":[{"id":"1","type":"comments"},{"id":"2","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}}}`
	articleRelatedCommentsNestedWithIncludeBody = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"},"relationships":{"comments":{"data":[{"id":"1","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}},"included":[{"id":"1","type":"comments","attributes":{"body":"A"},"relationships":{"author":{"data":{"id":"1","type":"author"},"links":{"self":"http://example.com/comments/1/relationships/author","related":"http://example.com/comments/1/author"}}}},{"id":"1","type":"author","attributes":{"name":"A"}}]}`

	articleWithIncludeOnlyBody = `{"data":{"id":"1","type":"articles","attributes":{"title":"A"}},"included":[{"id":"1","type":"author","attributes":{"name":"A"}}]}`

	// error structs
	errorsSimpleStruct         = Error{Title: "T"}
	errorsSimpleSliceSingle    = []Error{errorsSimpleStruct}
	errorsSimpleSliceSinglePtr = []*Error{&errorsSimpleStruct}
	errorsComplexStruct        = Error{
		ID:     "1",
		Links:  &ErrorLink{About: "A"},
		Status: "S",
		Code:   "C",
		Title:  "T",
		Detail: "D",
		Source: &ErrorSource{Pointer: "PO", Parameter: "PA"},
		Meta:   map[string]string{"K": "V"},
	}
	errorsComplexSliceMany    = []Error{errorsSimpleStruct, errorsComplexStruct}
	errorsComplexSliceManyPtr = []*Error{&errorsSimpleStruct, &errorsComplexStruct}

	// error bodies
	errorsSimpleStructBody     = `{"errors":[{"title":"T"}]}`
	errorsComplexStructBody    = `{"errors":[{"id":"1","links":{"about":"A"},"status":"S","code":"C","title":"T","detail":"D","source":{"pointer":"PO","parameter":"PA"},"meta":{"K":"V"}}]}`
	errorsComplexSliceManyBody = `{"errors":[{"title":"T"},{"id":"1","links":{"about":"A"},"status":"S","code":"C","title":"T","detail":"D","source":{"pointer":"PO","parameter":"PA"},"meta":{"K":"V"}}]}`
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
