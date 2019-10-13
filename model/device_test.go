package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	dType := DeviceType("abc")
	require.Error(t, dType.Valid(), "unknown type")

	require.Equal(t, TypeShortExpiry.Seconds(), 3600)
	require.Equal(t, TypeLongExpiry.Seconds(), 2592000)
	require.Equal(t, TypeNoExpiry.Seconds(), 31536000)
	value, err := TypeNoExpiry.Value()
	require.NoError(t, err)
	require.Equal(t, value, "NoExpiry")
	require.Nil(t, TypeNoExpiry.Valid())

	var scan DeviceType

	err = scan.Scan([]byte("NoExpiry"))
	require.NoError(t, err)
	require.Equal(t, TypeNoExpiry, scan)
}
