package logger

import (
	"context"
	"regexp"

	"strings"
	"unicode"

	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

var (
	passRegEx = regexp.MustCompile(`pass\s*:\s*"(?:[^"\\]|\\.)*"`)
)

// GQLLog logs graphql queries, mutations and errors.
func GQLLog(_ context.Context, params *graphql.Params, result *graphql.Result, _ []byte) {

	if result.HasErrors() {
		logMSG := "GQL: " + toOneLine(hidePassword(params.RequestString))
		var errs []string
		for _, err := range result.Errors {
			errs = append(errs, err.Message)
		}
		log.Error().Strs("error", errs).Msg(logMSG)
	} else if log.Debug().Enabled() {
		logMSG := "GQL: " + toOneLine(hidePassword(params.RequestString))
		log.Debug().Msg(logMSG)
	}
}

func toOneLine(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '\n' || r == '\r'
	}), " ")
}

func hidePassword(s string) string {
	return passRegEx.ReplaceAllString(s, `pass:"<hidden>"`)
}
