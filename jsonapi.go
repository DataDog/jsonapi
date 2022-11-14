// Package jsonapi implements encoding and decoding of JSON:API as defined in https://jsonapi.org/format/.
package jsonapi

import (
	"encoding/json"
	"fmt"
)

// ResourceObject is a JSON:API resource object as defined by https://jsonapi.org/format/1.0/#document-resource-objects
type resourceObject struct {
	ID            string               `json:"id"`
	Type          string               `json:"type"`
	Attributes    map[string]any       `json:"attributes,omitempty"`
	Relationships map[string]*document `json:"relationships,omitempty"`
	Meta          any                  `json:"meta,omitempty"`
	Links         *Link                `json:"links,omitempty"`
}

// JSONAPI is a JSON:API object as defined by https://jsonapi.org/format/1.0/#document-jsonapi-object.
type jsonAPI struct {
	Version string `json:"version"`
	Meta    any    `json:"meta,omitempty"`
}

// ErrorLink represents a JSON:API error links object as defined by https://jsonapi.org/format/1.0/#error-objects.
type ErrorLink struct {
	About string `json:"about,omitempty"`
}

// ErrorSource represents a JSON:API Error.Source as defined by https://jsonapi.org/format/1.0/#error-objects.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// Error represents a JSON:API error object as defined by https://jsonapi.org/format/1.0/#error-objects.
type Error struct {
	ID     string       `json:"id,omitempty"`
	Links  *ErrorLink   `json:"links,omitempty"`
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
	Meta   any          `json:"meta,omitempty"`
}

// LinkObject is a links object as defined by https://jsonapi.org/format/1.0/#document-links
type LinkObject struct {
	Href string `json:"href,omitempty"`
	Meta any    `json:"meta,omitempty"`
}

// Link is the top-level links object as defined by https://jsonapi.org/format/1.0/#document-top-level.
// First|Last|Next|Previous are provided to support pagination as defined by https://jsonapi.org/format/1.0/#fetching-pagination.
type Link struct {
	Self    any `json:"self,omitempty"`
	Related any `json:"related,omitempty"`

	First    string `json:"first,omitempty"`
	Last     string `json:"last,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

func checkLinkValue(linkValue any) (bool, *TypeError) {
	var isEmpty bool

	switch lv := linkValue.(type) {
	case *LinkObject:
		isEmpty = (lv.Href == "")
	case string:
		isEmpty = (lv == "")
	case nil:
		isEmpty = true
	default:
		return false, &TypeError{Actual: fmt.Sprintf("%T", lv), Expected: []string{"*LinkObject", "string"}}
	}

	return isEmpty, nil
}

func (l *Link) check() error {
	selfIsEmpty, err := checkLinkValue(l.Self)
	if err != nil {
		return err
	}

	relatedIsEmpty, err := checkLinkValue(l.Related)
	if err != nil {
		return err
	}

	// if both are empty then fail, and if one is empty, it must be set to nil to satisfy omitempty
	switch {
	case selfIsEmpty && relatedIsEmpty:
		return ErrMissingLinkFields
	case selfIsEmpty:
		l.Self = nil
	case relatedIsEmpty:
		l.Related = nil
	}

	return nil
}

// Document is a JSON:API document as defined by https://jsonapi.org/format/1.0/#document-top-level
type document struct {
	// Data is a ResourceObject as defined by https://jsonapi.org/format/1.0/#document-resource-objects.
	// DataOne/DataMany are translated to Data in document.MarshalJSON
	hasMany  bool
	DataOne  *resourceObject   `json:"-"`
	DataMany []*resourceObject `json:"-"`

	// Meta is Meta Information as defined by https://jsonapi.org/format/1.0/#document-meta.
	Meta any `json:"meta,omitempty"`

	// JSONAPI is a JSON:API object as defined by https://jsonapi.org/format/1.0/#document-jsonapi-object.
	JSONAPI *jsonAPI `json:"jsonapi,omitempty"`

	// Errors is a list of JSON:API error objects as defined by https://jsonapi.org/format/1.0/#error-objects.
	Errors []*Error `json:"errors,omitempty"`

	// Links is the top-level links object as defined by https://jsonapi.org/format/1.0/#document-top-level.
	Links *Link `json:"links,omitempty"`

	// Includes contains ResourceObjects creating a compound document as defined by https://jsonapi.org/format/#document-compound-documents.
	Included []*resourceObject `json:"included,omitempty"`
}

func newDocument() *document {
	return &document{
		DataMany: make([]*resourceObject, 0),
		Errors:   make([]*Error, 0),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (d *document) MarshalJSON() ([]byte, error) {
	// if we get errors, force exclusion of the Data field
	if len(d.Errors) > 0 {
		type alias document
		return json.Marshal(&struct{ *alias }{alias: (*alias)(d)})
	}

	// if DataMany is populated Data is a []*resourceObject
	if d.hasMany {
		type alias document
		return json.Marshal(&struct {
			Data []*resourceObject `json:"data"`
			*alias
		}{
			Data:  d.DataMany,
			alias: (*alias)(d),
		})
	}

	type alias document
	return json.Marshal(&struct {
		Data *resourceObject `json:"data"`
		*alias
	}{
		Data:  d.DataOne,
		alias: (*alias)(d),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *document) UnmarshalJSON(data []byte) (err error) {
	type alias document

	// Since there is no simple regular expression to capture only that the primary data is an
	// array, try unmarshaling both ways
	auxMany := &struct {
		Data []*resourceObject `json:"data"`
		*alias
	}{
		alias: (*alias)(d),
	}
	if err = json.Unmarshal(data, &auxMany); err == nil {
		d.hasMany = true
		d.DataMany = auxMany.Data
		return
	}

	auxOne := &struct {
		Data *resourceObject `json:"data"`
		*alias
	}{
		alias: (*alias)(d),
	}
	if err = json.Unmarshal(data, &auxOne); err == nil {
		d.DataOne = auxOne.Data
	}

	return
}

// verifyFullLinkage returns an error if the given compound document is not fully-linked as
// described by https://jsonapi.org/format/1.1/#document-compound-documents. That is, there must be
// a chain of relationships linking all included data to primary data transitively.
func (d *document) verifyFullLinkage() error {
	if len(d.Included) == 0 {
		return nil
	}

	getResourceObjectSlice := func(d *document) []*resourceObject {
		if d.hasMany {
			return d.DataMany
		}
		if d.DataOne == nil {
			return nil
		}
		return []*resourceObject{d.DataOne}
	}

	resourceIdentifier := func(ro *resourceObject) string {
		return fmt.Sprintf("{Type: %v, ID: %v}", ro.Type, ro.ID)
	}

	// a list of related resource identifiers, and a flag to mark nodes as visited
	type includeNode struct {
		relatedTo []string
		visited   bool
	}

	// compute a graph of relationships between just the included resources
	includeGraph := make(map[string]*includeNode)
	for _, included := range d.Included {
		relatedTo := make([]string, 0)

		for _, relationship := range included.Relationships {
			for _, ro := range getResourceObjectSlice(relationship) {
				relatedTo = append(relatedTo, resourceIdentifier(ro))
			}
		}

		includeGraph[resourceIdentifier(included)] = &includeNode{relatedTo: relatedTo}
	}

	// helper to traverse the graph from a given key and mark nodes as visited
	var visit func(identifier string)
	visit = func(identifier string) {
		node, ok := includeGraph[identifier]
		if !ok || node.visited {
			return
		}
		node.visited = true
		for _, related := range node.relatedTo {
			visit(related)
		}
	}

	// visit all include nodes that are accessible from the primary data
	primaryData := getResourceObjectSlice(d)
	for _, data := range primaryData {
		for _, relationship := range data.Relationships {
			for _, ro := range getResourceObjectSlice(relationship) {
				visit(resourceIdentifier(ro))
			}
		}
	}

	invalidResources := make([]string, 0)
	for identifier, node := range includeGraph {
		if !node.visited {
			invalidResources = append(invalidResources, identifier)
		}
	}

	if len(invalidResources) > 0 {
		return &PartialLinkageError{invalidResources}
	}

	return nil
}

// Linkable can be implemented to marshal resource object links as defined by https://jsonapi.org/format/1.0/#document-resource-object-links.
type Linkable interface {
	Link() *Link
}

// LinkableRelation can be implemented to marshal resource object related resource links as defined by https://jsonapi.org/format/1.0/#document-resource-object-related-resource-links.
type LinkableRelation interface {
	LinkRelation(relation string) *Link
}

// MarshalIdentifier can be optionally implemented to control marshaling of the primary field to a string.
//
// The order of operations for marshaling the primary field is:
//
//  1. Use MarshalIdentifier if it is implemented
//  2. Use the value directly if it is a string
//  3. Use fmt.Stringer if it is implemented
//  4. Fail
type MarshalIdentifier interface {
	MarshalID() string
}

// UnmarshalIdentifier can be optionally implemented to control unmarshaling of the primary field from a string.
//
// The order of operations for unmarshaling the primary field is:
//
//  1. Use UnmarshalIdentifier if it is implemented
//  2. Use the value directly if it is a string
//  3. Fail
type UnmarshalIdentifier interface {
	UnmarshalID(id string) error
}
