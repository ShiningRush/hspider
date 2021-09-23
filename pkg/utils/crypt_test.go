package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	g, err := Encrypt("ocIx0jgls4nRekK3PopR3aowJGF8")
	assert.NoError(t, err)
	assert.Equal(t, "lANR7lINpv9v4bjvrOSzO5n+r/O4DXfWaEzVja4qtqQchv0qspJJvnvY+2rwpOnpBmPf7R2dgwKpQ21Hgt2JK7gVi0/UPGrbMko1fJUEvhJYVEf8PoyBcLgpZBNVDhzPIEJ0tpch94B5JHtwWSEPlDfqCyblaJbpOEgob6eB1kk=", g)
}
