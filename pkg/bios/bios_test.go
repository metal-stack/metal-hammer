package bios

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBios(t *testing.T) {

	bios := Bios()
	assert.NotNil(t, bios)
	assert.NotEqual(t, bios.Version, "")
	assert.NotEqual(t, bios.Vendor, "")
	assert.NotEqual(t, bios.Date, "")
	assert.NotEqual(t, bios.String(), "")
	assert.True(t, strings.HasPrefix(bios.String(), "version"))
}
