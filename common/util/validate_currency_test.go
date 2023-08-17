package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestValidateCurrency(t *testing.T) {
	assert.NotNil(t, ValidateCurrency("KKK"))
	assert.Nil(t, ValidateCurrency("CNY"))
	assert.Nil(t, ValidateCurrency("USD"))
	assert.Nil(t, ValidateCurrency("THB"))
	assert.NotNil(t, ValidateCurrency("NNN"))
}
