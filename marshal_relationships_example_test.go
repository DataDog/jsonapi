package jsonapi_test

import (
	"fmt"

	"github.com/DataDog/jsonapi"
)

type Author struct {
	ID   string `jsonapi:"primary,author"`
	Name string `jsonapi:"attribute" json:"name"`
}

type Comment struct {
	ID     string  `jsonapi:"primary,comments"`
	Body   string  `jsonapi:"attribute" json:"comment"`
	Author *Author `jsonapi:"relationship"`
}

func (c *Comment) LinkRelation(relation string) *jsonapi.Link {
	return &jsonapi.Link{
		Self:    fmt.Sprintf("http://example.com/comments/%s/relationships/%s", c.ID, relation),
		Related: fmt.Sprintf("http://example.com/comments/%s/%s", c.ID, relation),
	}
}

type Article struct {
	ID       string     `jsonapi:"primary,articles"`
	Title    string     `jsonapi:"attribute" json:"title"`
	Author   *Author    `jsonapi:"relationship" json:"author,omitempty"`
	Comments []*Comment `jsonapi:"relationship" json:"comments,omitempty"`
}

func (a *Article) LinkRelation(relation string) *jsonapi.Link {
	return &jsonapi.Link{
		Self:    fmt.Sprintf("http://example.com/articles/%s/relationships/%s", a.ID, relation),
		Related: fmt.Sprintf("http://example.com/articles/%s/%s", a.ID, relation),
	}
}

func ExampleMarshal_relationships() {
	authorA := &Author{ID: "AA", Name: "Cool Author"}
	authorB := &Author{ID: "AB", Name: "Cool Commenter"}
	authorC := &Author{ID: "AC", Name: "Neat Commenter"}
	commentA := &Comment{ID: "CA", Body: "Very cool", Author: authorB}
	commentB := &Comment{ID: "CB", Body: "Super neat", Author: authorC}
	article := Article{
		ID:       "1",
		Title:    "Hello World",
		Author:   authorA,
		Comments: []*Comment{commentA, commentB},
	}

	b, err := jsonapi.Marshal(&article)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
	// Output: {"data":{"id":"1","type":"articles","attributes":{"title":"Hello World"},"relationships":{"author":{"data":{"id":"AA","type":"author"},"links":{"self":"http://example.com/articles/1/relationships/author","related":"http://example.com/articles/1/author"}},"comments":{"data":[{"id":"CA","type":"comments"},{"id":"CB","type":"comments"}],"links":{"self":"http://example.com/articles/1/relationships/comments","related":"http://example.com/articles/1/comments"}}}}}
}
