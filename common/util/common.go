package util

import (
	"os"
	"reflect"
	"strconv"
)

// IsInArray is used to check the input string whether is in the string array
func IsInArray(array []string, iv string) bool {
	if 0 == len(array) {
		return false
	}

	for _, v := range array {
		if v == iv {
			return true
		}
	}

	return false
}

// IsArrayEqual is used to check whether the two arrays are equal
func IsArrayEqual(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !IsInArray(a, b[i]) {
			return false
		}
	}

	return true
}

// StringInList is used to check whether the string is in the string list
func StringInList(key string, strList *[]string) bool {
	for _, k := range *strList {
		if k == key {
			return true
		}
	}
	return false
}

// Uint64InList is used to check whether the uint64 is in the uint64 list
func Uint64InList(key uint64, list []uint64) bool {
	for _, v := range list {
		if key == v {
			return true
		}
	}
	return false
}

// EnvDefaultString is used to get string value from OS environment,
// returns default string if the key is not find in the OS environment.
func EnvDefaultString(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// EnvDefaultInt is used to get int value from OS environment, returns default int if the the key is not find
// in the OS environment or failed to parse the string value to int.
func EnvDefaultInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}

	return n
}

// EnvSetString sets the value of the environment variable named by the key.
func EnvSetString(key, defaultValue string) error {
	err := os.Setenv(key, defaultValue)
	if err != nil {
		return err
	}
	return nil
}

// GetEither gets the value in the map in order according to the two keys
func GetEither(m map[string]string, key1, key2 string) string {
	if nil == m {
		return ""
	}
	if v, ok := m[key1]; ok {
		return v
	}

	return m[key2]
}

// IsNil returns true if input parameter is nil
func IsNil(i interface{}) bool {
	if nil == i {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// MapToString format map as a string
func MapToString(m map[string]string) string {
	str := "["
	if nil == m {
		str += "]"
		return str
	}
	idx := 0
	for k, v := range m {
		idx++
		if idx > 1 {
			str += " "
		}
		str += k + ":" + v
	}

	str += "]"
	return str

}
