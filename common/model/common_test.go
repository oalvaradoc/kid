package model

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestBuildErrorResponse(t *testing.T) {
	er := BuildErrorResponse("error message from unit test")
	assert.NotNil(t, er)
	assert.Equal(t, er.ErrorCode, -1)
	assert.Equal(t, er.ErrorMsg, "error message from unit test")
}

func TestBuildErrorResponseWithErrorCode(t *testing.T) {
	er := BuildErrorResponseWithErrorCode(100, "error message from unit test")
	assert.NotNil(t, er)
	assert.Equal(t, er.ErrorCode, 100)
	assert.Equal(t, er.ErrorMsg, "error message from unit test")
}

func TestBuildResponse(t *testing.T) {
	br := BuildResponse("test")
	assert.NotNil(t, br)
	assert.Equal(t, br.ErrorCode, 0)
	assert.Equal(t, br.ErrorMsg, "")
}

func TestBuildPagingResponse(t *testing.T) {
	pr := BuildPagingResponse(100, []string{"a", "b", "c"})
	assert.NotNil(t, pr)
	assert.Equal(t, pr.ErrorCode, 0)
	assert.Equal(t, pr.ErrorMsg, "")
	assert.NotNil(t, pr.Data)
}
