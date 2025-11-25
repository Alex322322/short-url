package random_test

import (
	"strings"
	"testing"

	"github.com/Alex322322/short-url/internal/lib/random"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAlias(t *testing.T) {
	t.Run("generates correct length", func(t *testing.T) {
		length := 10
		alias := random.GenerateAlias(length)

		assert.Equal(t, length, len(alias), "Alias should have correct length")
	})

	t.Run("generates different aliases", func(t *testing.T) {
		alias1 := random.GenerateAlias(8)
		alias2 := random.GenerateAlias(8)

		assert.NotEqual(t, alias1, alias2, "Subsequent calls should generate different aliases")
	})

	t.Run("contains only valid characters", func(t *testing.T) {
		alias := random.GenerateAlias(20)
		validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		for _, char := range alias {
			assert.True(t, strings.ContainsRune(validChars, char),
				"Alias should contain only valid characters: %s", string(char))
		}
	})

	t.Run("handles zero length", func(t *testing.T) {
		alias := random.GenerateAlias(0)

		assert.Equal(t, "", alias, "Zero length should return empty string")
	})

	t.Run("handles very short length", func(t *testing.T) {
		alias := random.GenerateAlias(1)

		assert.Equal(t, 1, len(alias), "Should handle single character alias")
	})

	t.Run("handles very long length", func(t *testing.T) {
		alias := random.GenerateAlias(1000)

		assert.Equal(t, 1000, len(alias), "Should handle long aliases")
	})
}
