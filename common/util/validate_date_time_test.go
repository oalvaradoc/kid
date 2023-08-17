package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestValidateDate(t *testing.T) {
	assert.Nil(t, ValidateDate("2021-07-12"))
	assert.NotNil(t, ValidateDate("2021/07/12"))
}

func TestValidateTime(t *testing.T) {
	assert.Nil(t, ValidateTime("11:59:30"))
	assert.NotNil(t, ValidateTime("11.59.30"))
}

func TestValidateDateTime(t *testing.T) {
	assert.Nil(t, ValidateDateTime("2021-07-12 11:59:30"))
	assert.NotNil(t, ValidateDateTime("2021/07/12 11.59.30"))
}
