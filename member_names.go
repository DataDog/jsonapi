package jsonapi

import "regexp"

var (
	defaultNameRegex *regexp.Regexp
	strictNameRegex  *regexp.Regexp
)

func init() {
	defaultNameRegex = regexp.MustCompile(`^([a-zA-Z\d]|[^\x{0000}-\x{0080}])(([a-zA-Z\d]|[^\x{0000}-\x{0080}])|[-_ ]([a-zA-Z\d]|[^\x{0000}-\x{0080}]))*$`)
	// properties of the strict name regex:
	// - at least one lower case letter
	// - camel case, and must end with a lower case letter
	// - may have digits inside the word
	strictNameRegex = regexp.MustCompile(`^([a-z]|[a-z]+((\d)|([A-Z\d][a-z\d]+))*([A-Z\d][a-z\d]*[a-z]))$`)
}

type memberNameValidationMode int

const (
	// defaultValidation verifies that member names are valid according to the spec in
	// https://jsonapi.org/format/#document-member-names.
	//
	// Note that this validation mode allows for non-URL-safe member names.
	defaultValidation memberNameValidationMode = iota

	// disableValidation turns off member name validation for convenience and performance-saving
	// reasons.
	//
	// Note that this validation mode allows for member names to not conform to the JSON:API spec.
	disableValidation

	// strictValidation verifies that member names are both valid according to the spec in
	// https://jsonapi.org/format/#document-member-names, and follow recommendations from
	// https://jsonapi.org/recommendations/#naming.
	//
	// Note that these names are always URL-safe.
	strictValidation
)

func isValidMemberName(name string, mode memberNameValidationMode) bool {
	switch mode {
	case disableValidation:
		return true
	case strictValidation:
		return strictNameRegex.MatchString(name)
	default:
		return defaultNameRegex.MatchString(name)
	}
}
