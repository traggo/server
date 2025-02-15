package logger

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
)

var (
	passRegEx = regexp.MustCompile(`pass\s*:\s*"(?:[^"\\]|\\.)*"`)
)

// GQLLog logs graphql queries, mutations and errors.
func GQLLog() graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		start := time.Now()
		result := next(ctx)
		elapsed := time.Now().Sub(start)
		errs := graphql.GetErrors(ctx)
		rawQuery := graphql.GetRequestContext(ctx).RawQuery

		if len(errs) > 0 {
			var errorStrings []string
			for _, err := range errs {
				errorStrings = append(errorStrings, err.Error())
			}

			log.Error().Strs("error", errorStrings).Str("took", elapsed.String()).Msg("GQL: " + toOneLine(hidePassword(rawQuery)))
		} else if log.Debug().Enabled() {
			log.Debug().Str("took", elapsed.String()).Msg("GQL: " + toOneLine(hidePassword(rawQuery)))
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
