package logger_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/test"
)

var (
	params = []struct {
		query    string
		expected string
	}{
		{
			query:    `mutation { createUser(name: "name", pass: "unicorn", admin:true) { name admin id } }`,
			expected: `GQL: mutation { createUser(name: "name", pass:"<hidden>", admin:true) { name admin id } }`,
		},
		{
			query:    `mutation { createUser(name: "name", pass     : "unicorn", admin:true) { name admin id } }`,
			expected: `GQL: mutation { createUser(name: "name", pass:"<hidden>", admin:true) { name admin id } }`,
		},
		{
			query:    `mutation { createUser(name: "name", pass:"unicorn", admin:true) { name admin id } }`,
			expected: `GQL: mutation { createUser(name: "name", pass:"<hidden>", admin:true) { name admin id } }`,
		},
		{
			query:    `mutation { createUser(name: "name", pass:"uni\"corn", admin:true) { name admin id } }`,
			expected: `GQL: mutation { createUser(name: "name", pass:"<hidden>", admin:true) { name admin id } }`,
		},
		{
			query:    `mutation { createUser(name: "name", pass:"23456789\n\r0!\"§$%><>|&/()=?*'ÄÜ'", admin:true) { name admin id } }`,
			expected: `GQL: mutation { createUser(name: "name", pass:"<hidden>", admin:true) { name admin id } }`,
		},
	}
)

func TestGQLLog_hidesPassword_withoutErrors(t *testing.T) {
	for i, param := range params {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			fakeLog := test.NewLogger(t)
			defer fakeLog.Dispose()
			gqlLog(param.query)
			fakeLog.AssertCount(1)
			fakeLog.AssertEntryExists(test.Entry{Level: zerolog.DebugLevel, Message: param.expected})
		})
	}
}

func TestGQLLog_hidesPassword_withErrors(t *testing.T) {
	for i, param := range params {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			fakeLog := test.NewLogger(t)
			defer fakeLog.Dispose()
			gqlLog(param.query, errors.New("oops"))
			fakeLog.AssertCount(1)
			fakeLog.AssertEntryExists(test.Entry{Level: zerolog.ErrorLevel, Message: param.expected})
		})
	}
}

func TestGQLLog_withoutErrors_removesNewlinesAndSpaces(t *testing.T) {
	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	query := `
query {
   users {
      id,
      name,
      admin
   }
}`

	gqlLog(query)

	expectedMessage := `GQL: query { users { id, name, admin } }`

	fakeLog.AssertCount(1)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.DebugLevel, Message: expectedMessage})
}

func TestGQLLog_withErrors_removesNewlinesAndSpaces(t *testing.T) {
	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	query := `
query {
   users {
      id,
      name,
      admin
   }
}`

	gqlLog(query, errors.New("oops"))

	expectedMessage := `GQL: query { users { id, name, admin } }`

	fakeLog.AssertCount(1)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.ErrorLevel, Message: expectedMessage})
}

func gqlLog(query string, errs ...error) {
	ctx := graphql.WithRequestContext(context.Background(), &graphql.RequestContext{
		RawQuery:       query,
		ErrorPresenter: graphql.DefaultErrorPresenter,
	})
	for _, err := range errs {
		graphql.AddError(ctx, err)
	}
	logger.GQLLog()(ctx, func(ctx context.Context) []byte {
		return []byte{}
	})
}
