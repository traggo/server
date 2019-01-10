package rand

import (
	"crypto/rand"

	"github.com/rs/zerolog/log"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randRead = rand.Read
)

func init() {
	// exits with 1 when crypto/rand is not available
	Token(1)
}

func randBytes(count int) ([]byte, error) {
	bytes := make([]byte, count)
	_, err := randRead(bytes)
	return bytes, err
}

// Token returns a random token.
func Token(count int) string {
	bytes, err := randBytes(count)
	if err != nil {
		log.Panic().Msg("crypto/rand is not available")
	}
	return bytesToString(bytes, count)
}

func bytesToString(bytes []byte, count int) string {
	token := make([]rune, count)

	for index, item := range bytes {
		token[index] = chars[int(item)%len(chars)]
	}

	return string(token)
}
