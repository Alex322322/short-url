package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	once sync.Once
	rnd  *rand.Rand
)

func init() {
	once.Do(func() {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
}

func GenerateAlias(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}

	return string(b)
}
