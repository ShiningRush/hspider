package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLoginUrl(t *testing.T) {
	url, err := getLoginUrl("ocIx0jgls4nRekK3PopR3aowJGF8")
	assert.NoError(t, err)
	assert.Equal(t, "", url)
}
