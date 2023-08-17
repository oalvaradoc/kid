package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"os"
	"testing"
)

func TestIsInArray(t *testing.T) {
	assert.True(t, IsInArray([]string{"a", "b", "c"}, "b"))
	assert.False(t, IsInArray([]string{"a", "b", "c"}, "d"))
	assert.False(t, IsInArray(nil, "e"))
}

func TestIsArrayEqual(t *testing.T) {
	assert.True(t, IsArrayEqual(nil, nil))
	assert.False(t, IsArrayEqual([]string{}, nil))
	assert.False(t, IsArrayEqual(nil, []string{}))
	assert.True(t, IsArrayEqual([]string{"a", "b", "c"}, []string{"a", "b", "c"}))
	assert.False(t, IsArrayEqual([]string{"a", "b", "c"}, []string{"A", "B", "C"}))
}

func TestStringInList(t *testing.T) {
	assert.False(t, StringInList("", &[]string{"a"}))
	assert.False(t, StringInList("a", &[]string{}))
	assert.True(t, StringInList("a", &[]string{"a", "b", "c"}))
}

func TestUint64InList(t *testing.T) {
	assert.False(t, Uint64InList(0, []uint64{}))
	assert.False(t, Uint64InList(1, []uint64{}))
	assert.True(t, Uint64InList(1, []uint64{0, 1, 2}))
}

func TestEnvDefaultString(t *testing.T) {
	os.Unsetenv("the_key_does_not_exist")
	assert.Equal(t, EnvDefaultString("the_key_does_not_exist", "default_value_1"), "default_value_1")

	os.Setenv("test_key", "test_value")
	assert.Equal(t, EnvDefaultString("test_key", "default_value_2"), "test_value")
}

func TestEnvDefaultInt(t *testing.T) {
	os.Unsetenv("the_key_does_not_exist")
	assert.Equal(t, EnvDefaultInt("the_key_does_not_exist", 100), 100)

	os.Setenv("test_key2", "2")
	assert.Equal(t, EnvDefaultInt("test_key2", 200), 2)
}

func TestEnvSetString(t *testing.T) {
	assert.Nil(t, EnvSetString("TestEnvSetString_test_key", "value"))
	os.Unsetenv("TestEnvSetString_test_key")
}


func TestGetEither(t *testing.T) {
	m := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}
	assert.Equal(t, GetEither(m, "the_key_does_not_exist", "key1"), "value1")
	assert.Equal(t, GetEither(m, "key2", "key3"), "value2")
}

func TestIsNil(t *testing.T) {
	assert.True(t, IsNil(nil))
	assert.False(t, IsNil(""))
}

func TestMapToString(t *testing.T) {
	assert.Equal(t, MapToString(map[string]string{"key1": "value1", "key2": "value2"}), `[key1:value1 key2:value2]`)
	assert.Equal(t, MapToString(map[string]string{}), `[]`)
}
