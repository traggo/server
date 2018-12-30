package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/model"
)

func TestAll_NoDuplicateEntries(t *testing.T) {
	all := model.All()
	var checkedItems []interface{}
	for _, item := range all {
		assert.NotContains(t, checkedItems, item)
		checkedItems = append(checkedItems, item)
	}
}
