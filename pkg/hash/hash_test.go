package hash_test

import (
	"github.com/stretchr/testify/assert"
	"golang.project/go-fundamentals/gameapp/pkg/hash"
	"testing"
)

func TestHash(t *testing.T) {
	plain := "hash_test_string"

	hashPlain, err := hash.Hash(plain)
	assert.NoError(t, err, "error hashing function")
	assert.NotEmpty(t, hashPlain, "hashing plain, is empty")

	matched := hash.CheckHash(plain, hashPlain)
	assert.True(t, matched, "should be must validate hashPlain")
}
