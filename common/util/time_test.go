package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
	"time"
)

func TestValidDate(t *testing.T) {
	assert.Nil(t, ValidDate("2021-07-12"))
}

func TestGetFirstDateOfMonth(t *testing.T) {
	d := time.Now()
	timestamp := time.Now().AddDate(0, 0, -d.Day()+1).Unix()
	_ = time.Unix(timestamp, 0)
	//return tm.Format("2006-01-02")
}

func TestGetCurrentTimeFmt(t *testing.T) {
	assert.NotNil(t, GetCurrentTimeFmt())
}

func TestGetTimeFormat(t *testing.T) {
	assert.NotNil(t, GetTimeFormat(time.Now().Unix()))
}

func TestGetTimeParse(t *testing.T) {
	res, err := GetTimeParse("2021-07-12 11:54:01")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestGetCurTime(t *testing.T) {
	assert.NotNil(t, GetCurTime())
}

func TestGetCurrentTimeStamp(t *testing.T) {
	v1, v2 := GetCurrentTimeStamp()
	assert.NotNil(t, v1)
	assert.NotNil(t, v2)
}

func TestGetCurrentDateTime(t *testing.T) {
	assert.NotNil(t, GetCurrentDateTime())
}

func TestGetCurrentDate(t *testing.T) {
	assert.NotNil(t, GetCurrentDate())
}

func TestGetCurrentTime(t *testing.T) {
	assert.NotNil(t, GetCurrentTime())
}

func TestGetPreMonthDate(t *testing.T) {
	assert.NotNil(t, GetPreMonthDate(6))
}

func TestGetTHCurrentDateTime(t *testing.T) {
	assert.NotNil(t, GetTHCurrentDateTime())
}
