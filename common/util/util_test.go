package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestIndexOfItem(t *testing.T) {
	array := []string{
		"1", "2", "3", "4", "5",
	}
	assert.Equal(t, IndexOfItem("3", array), 2)
}

func TestCurrentTime(t *testing.T) {
	ct := CurrentTime()
	assert.NotNil(t, ct)
	t.Logf("the current time:%s", ct)
}

func TestFormatTime(t *testing.T) {
	tv := ToTime("2021-07-12 16:37:05")
	ft := FormatTime(tv)
	assert.Equal(t, ft, "2021-07-12 16:37:05.000000000")
}

func TestToTime(t *testing.T) {
	v := ToTime("2021-07-12 16:37:05")
	assert.NotNil(t, v)
	t.Logf("the result of to time:%++v", v)
}

func TestCurrentHost(t *testing.T) {
	host := CurrentHost()
	assert.NotNil(t, host)
	t.Logf("current host:%s", host)
}

func TestFnv32(t *testing.T) {
	fnv32 := Fnv32("this is a test string")
	t.Logf("Fnv32:%d", fnv32)
}

func TestIntToBytes(t *testing.T) {
	bs, err := IntToBytes(10)
	assert.Nil(t, err)
	assert.NotNil(t, bs)
}

func TestBytesToInt(t *testing.T) {
	bs, err := IntToBytes(1024)
	assert.Nil(t, err)
	assert.NotNil(t, bs)

	v, err := BytesToInt(bs)
	assert.Nil(t, err)
	assert.Equal(t, 1024, v)
}

func TestGenerateSerialNo(t *testing.T) {
	serialNo := GenerateSerialNo("org", "wks", "env", "su", "service-0", "0")
	assert.NotNil(t, serialNo)
	t.Logf("serial number:%s", serialNo)
}

func TestRandomString(t *testing.T) {
	str := RandomString(10)
	assert.Equal(t, len(str), 10)

	assert.False(t, RandomString(10) == RandomString(10))
}

func TestGetMapValueIgnoreCase(t *testing.T) {
	m := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	v := GetMapValueIgnoreCase(m, "K1")
	assert.NotNil(t, v)
}
