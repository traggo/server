package rand

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	for i := 1; i < 10; i++ {
		assert.Len(t, Token(i), i)
	}
	assert.Len(t, Token(120), 120)
	assert.Len(t, Token(222), 222)
	assert.Len(t, Token(555665), 555665)
}

func TestToken_panics(t *testing.T) {
	old := randRead
	defer func() {
		randRead = old
	}()
	randRead = func(b []byte) (n int, err error) {
		return 0, errors.New("oops")
	}
	assert.Panics(t, func() {
		Token(1)
	})
}
