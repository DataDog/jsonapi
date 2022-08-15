package jsonapi_test

import (
	"fmt"

	"github.com/DataDog/jsonapi"
)

func ExampleMarshal() {
	type Article struct {
		ID    string `jsonapi:"primary,articles"`
		Title string `jsonapi:"attribute" json:"title"`
	}

	a := Article{ID: "1", Title: "Hello World"}

	b, err := jsonapi.Marshal(&a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
	// Output: {"data":{"id":"1","type":"articles","attributes":{"title":"Hello World"}}}
}

func ExampleMarshal_slice() {
	type Article struct {
		ID    string `jsonapi:"primary,articles"`
		Title string `jsonapi:"attribute" json:"title"`
	}

	a := []*Article{
		{ID: "1", Title: "Hello World"},
		{ID: "2", Title: "Hello Again"},
	}

	b, err := jsonapi.Marshal(&a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
	// Output: {"data":[{"id":"1","type":"articles","attributes":{"title":"Hello World"}},{"id":"2","type":"articles","attributes":{"title":"Hello Again"}}]}
}

func ExampleMarshal_meta() {
	type ArticleMeta struct {
		Views int `json:"views"`
	}
	type Article struct {
		ID    string       `jsonapi:"primary,articles"`
		Title string       `jsonapi:"attribute" json:"title"`
		Meta  *ArticleMeta `jsonapi:"meta"`
	}

	a := Article{ID: "1", Title: "Hello World", Meta: &ArticleMeta{Views: 10}}
	m := map[string]any{"foo": "bar"}

	b, err := jsonapi.Marshal(&a, jsonapi.MarshalMeta(m))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(b))
	// Output: {"data":{"id":"1","type":"articles","attributes":{"title":"Hello World"},"meta":{"views":10}},"meta":{"foo":"bar"}}
}
