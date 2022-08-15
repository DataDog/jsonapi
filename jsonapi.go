// Package jsonapi implements encoding and decoding of JSON:API as defined in https://jsonapi.org/format/.
package jsonapi

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var dataArrayRegex *regexp.Regexp

func init() {
	// matches "data":[ (ignoring any space between : and [)
	dataArrayRegex = regexp.MustCompile(`"data":\s*\[`)
}

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
func (d *document) UnmarshalJSON(data []byte) error {
	if dataArrayRegex.Match(data) {
		type alias document
		aux := &struct {
			Data []*resourceObject `json:"data"`
			*alias
		}{
			alias: (*alias)(d),
		}
		if err := json.Unmarshal(data, &aux); err != nil {
			return err
		}
		d.hasMany = true
		d.DataMany = aux.Data
		return nil
	}

	type alias document
	aux := &struct {
		Data *resourceObject `json:"data"`
		*alias
	}{
		alias: (*alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	d.DataOne = aux.Data
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
