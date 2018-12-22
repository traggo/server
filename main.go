package main

import (
	"net/http"

	"github.com/traggo/server/database"
	"github.com/traggo/server/gql"
)

func main() {
	// TODO configurable
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", gql.Handler(db, 10))
	http.ListenAndServe(":3030", nil) // TODO configurable port
}
