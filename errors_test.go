package jsonapi

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestErrorMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	expected := []byte(`{"id":"1","links":{"about":"A"},"status":"500","code":"C","title":"T","detail":"D","source":{"pointer":"PO","parameter":"PA"},"meta":{"K":"V"}}`)

	var e Error
	err := json.Unmarshal(expected, &e)
	is.MustNoError(t, err)

	actual, err := json.Marshal(&e)
	is.MustNoError(t, err)
	is.EqualJSON(t, string(expected), string(actual))
}
