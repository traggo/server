package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryDB(t *testing.T) {
	db := InMemoryDB(t)
	defer db.Close()
	assert.NotNil(t, db)
}
