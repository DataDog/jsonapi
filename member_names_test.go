package jsonapi

import (
	"testing"

	"github.com/DataDog/jsonapi/internal/is"
)

func TestIsValidMemberName(t *testing.T) {
	t.Parallel()

	// associate member name strings with the strictest validation mode they should pass
	testValidations := map[MemberNameValidationMode][]string{
		StrictValidation: {
			"a",
			"lowercase1with2numerals",
			"camelCase",
			"camel12Case9WithNumera1s",
		},
		DefaultValidation: {
			"A",
			"9camelCaseWithNumeralPrefix",
			"camelCaseWithNumeralSuffix10",
			"4camelCaseWithSurroundingNumerals5",
			"camelC",
			"PascalCase",
			"dash-case",
			"snake_case",
			"space case",
			"cRaZyCasE",
			"12",
			"ƒunky unicode",
		},
		DisableValidation: {
			"bad%character",
		},
	}

	for mode, names := range testValidations {
		mode := mode
		for _, name := range names {
			name := name
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				passesStrict := isValidMemberName(name, StrictValidation)
				passesDefault := isValidMemberName(name, DefaultValidation)

				switch mode {
				case StrictValidation:
					is.Equal(t, true, passesStrict)
					is.Equal(t, true, passesDefault)
				case DefaultValidation:
					is.Equal(t, false, passesStrict)
					is.Equal(t, true, passesDefault)
				case DisableValidation:
					is.Equal(t, false, passesStrict)
					is.Equal(t, false, passesDefault)
				}
			})
		}
	}
}
