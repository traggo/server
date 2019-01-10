package device

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func assertDeviceExist(t *testing.T, db *test.Database, expected model.Device) {
	foundDevice := new(model.Device)
	find := db.Find(foundDevice, expected.ID)
	require.Nil(t, find.Error)
	require.NotNil(t, foundDevice)
	require.Equal(t, expected, *foundDevice)
}

func assertDeviceCount(t *testing.T, db *test.Database, expected int) {
	count := new(int)
	db.Model(new(model.Device)).Count(count)
	require.Equal(t, expected, *count)
}
