package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestMysqlEscape(t *testing.T) {
	res, err := MysqlEscape("")
	assert.NotNil(t, err)
	assert.Equal(t, res, "")

	res, err = MysqlEscape("SELECT * FROM test")
	assert.Nil(t, err)
	assert.Equal(t, res, "SELECT * FROM test")

	res, err = MysqlEscape("SELECT * FROM 'test'")
	assert.Nil(t, err)
	assert.Equal(t, res, "SELECT * FROM \\'test\\'")

	res, err = MysqlEscape("SELECT \"field1\" FROM 'test'")
	assert.Nil(t, err)
	assert.Equal(t, res, "SELECT \\\"field1\\\" FROM \\'test\\'")
}
