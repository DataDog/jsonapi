package jsonapi

import (
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

type resourceIdentifierMetaTestNode struct {
	ID               string `jsonapi:"primary,nodes"`
	IdentifierMeta   any
	RelationshipMeta any `jsonapi:"meta"`
}

func (n resourceIdentifierMetaTestNode) MarshalResourceIdentifierMeta() any {
	return n.IdentifierMeta
}

type resourceIdentifierMetaTestLayer struct {
	ID    string                           `jsonapi:"primary,layers"`
	Node  *resourceIdentifierMetaTestNode  `jsonapi:"relationship" json:"node,omitempty"`
	Nodes []resourceIdentifierMetaTestNode `jsonapi:"relationship" json:"nodes,omitempty"`
}

func TestMarshalResourceIdentifierMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		given       any
		expect      string
		expectError error
	}{
		{
			description: "to-one relationship",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID:             "node-a",
					IdentifierMeta: map[string]any{"index": 0},
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"node":{"data":{"id":"node-a","type":"nodes","meta":{"index":0}}}}}}`,
		},
		{
			description: "to-many relationship with distinct metadata",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Nodes: []resourceIdentifierMetaTestNode{
					{ID: "node-a", IdentifierMeta: map[string]any{"index": 0}},
					{ID: "node-b", IdentifierMeta: map[string]any{"index": 1}},
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"nodes":{"data":[{"id":"node-a","type":"nodes","meta":{"index":0}},{"id":"node-b","type":"nodes","meta":{"index":1}}]}}}}`,
		},
		{
			description: "nil metadata",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID: "node-a",
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"node":{"data":{"id":"node-a","type":"nodes"}}}}}`,
		},
		{
			description: "typed nil metadata",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID:             "node-a",
					IdentifierMeta: (*struct{})(nil),
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"node":{"data":{"id":"node-a","type":"nodes"}}}}}`,
		},
		{
			description: "empty metadata object",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID:             "node-a",
					IdentifierMeta: map[string]any{},
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"node":{"data":{"id":"node-a","type":"nodes","meta":{}}}}}}`,
		},
		{
			description: "invalid metadata",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID:             "node-a",
					IdentifierMeta: "invalid",
				},
			},
			expectError: &TypeError{Actual: "string", Expected: []string{"struct", "map"}},
		},
		{
			description: "resource identifier and relationship metadata",
			given: &resourceIdentifierMetaTestLayer{
				ID: "layer-a",
				Node: &resourceIdentifierMetaTestNode{
					ID:               "node-a",
					IdentifierMeta:   map[string]any{"index": 0},
					RelationshipMeta: map[string]any{"total": 2},
				},
			},
			expect: `{"data":{"id":"layer-a","type":"layers","relationships":{"node":{"data":{"id":"node-a","type":"nodes","meta":{"index":0}},"meta":{"total":2}}}}}`,
		},
		{
			description: "top-level resource ignores resource identifier metadata",
			given: &resourceIdentifierMetaTestNode{
				ID:             "node-a",
				IdentifierMeta: map[string]any{"index": 0},
			},
			expect: `{"data":{"id":"node-a","type":"nodes"}}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

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
