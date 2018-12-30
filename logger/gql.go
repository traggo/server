package logger

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
)

var (
	passRegEx = regexp.MustCompile(`pass\s*:\s*"(?:[^"\\]|\\.)*"`)
)

// GQLLog logs graphql queries, mutations and errors.
func GQLLog() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		result := next(ctx)

		reqCtx := graphql.GetRequestContext(ctx)

		if len(reqCtx.Errors) > 0 {
			var errs []string
			for _, err := range reqCtx.Errors {
				errs = append(errs, err.Message)
			}

			log.Error().Strs("error", errs).Msg("GQL: " + toOneLine(hidePassword(reqCtx.RawQuery)))
		} else if log.Debug().Enabled() {
			log.Debug().Msg("GQL: " + toOneLine(hidePassword(reqCtx.RawQuery)))
		}

		return result
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
