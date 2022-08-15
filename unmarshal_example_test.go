package jsonapi_test

import (
	"fmt"

	"github.com/DataDog/jsonapi"
)

func ExampleUnmarshal() {
	body := `{"data":{"id":"1","type":"articles","attributes":{"title":"Hello World"}}}`

	type Article struct {
		ID    string `jsonapi:"primary,articles"`
		Title string `jsonapi:"attribute" json:"title"`
	}

	var a Article
	err := jsonapi.Unmarshal([]byte(body), &a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", &a)
	// Output: &{ID:1 Title:Hello World}
}

func ExampleUnmarshal_slice() {
	body := `{"data":[{"id":"1","type":"articles","attributes":{"title":"Hello World"}},{"id":"2","type":"articles","attributes":{"title":"Hello Again"}}]}`

	type Article struct {
		ID    string `jsonapi:"primary,articles"`
		Title string `jsonapi:"attribute" json:"title"`
	}

	var a []*Article
	err := jsonapi.Unmarshal([]byte(body), &a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v %+v", a[0], a[1])
	// Output: &{ID:1 Title:Hello World} &{ID:2 Title:Hello Again}
}

func ExampleUnmarshalMeta() {
	body := `{"data":{"id":"1","type":"articles","attributes":{"title":"Hello World"},"meta":{"views":10}},"meta":{"foo":"bar"}}`

	type ArticleMeta struct {
		Views int `json:"views"`
	}
	type Article struct {
		ID    string       `jsonapi:"primary,articles"`
		Title string       `jsonapi:"attribute" json:"title"`
		Meta  *ArticleMeta `jsonapi:"meta"`
	}

	var (
		a Article
		m map[string]any
	)
	err := jsonapi.Unmarshal([]byte(body), &a, jsonapi.UnmarshalMeta(&m))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s %+v %+v", a.ID, a.Title, a.Meta, m)
	// Output: 1 Hello World &{Views:10} map[foo:bar]
}
