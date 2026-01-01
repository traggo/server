package graphql

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestHandler_jsonOverHtml(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := NewResolver(db.DB, 4, model.Version{})
	handler := Handler("/gql", resolver, NewDirective())
	req := httptest.NewRequest("GET", "/gql?query="+url.QueryEscape("query {currentUser {name}}"), strings.NewReader(""))
	req.Header.Set("Accept", "text/html;application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	require.Equal(t, "application/graphql-response+json", recorder.Header().Get("Content-Type"))
	require.JSONEq(t, `
{
   "data": { "currentUser": null }
}
`, recorder.Body.String())
}

func TestHandler_htmlIfNotJson(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := NewResolver(db.DB, 4, model.Version{})
	handler := Handler("/gql", resolver, NewDirective())
	req := httptest.NewRequest("get", "/gql", strings.NewReader(""))
	req.Header.Set("Accept", "text/html;application/xml")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	require.Contains(t, recorder.Body.String(), "Traggo Playground")
}
