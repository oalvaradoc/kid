package json

import (
	jsoniter "github.com/json-iterator/go"
)

var jsonite = jsoniter.ConfigCompatibleWithStandardLibrary

// Unmarshal decodes the json byte array into target object
func Unmarshal(data []byte, v interface{}) error {
	return jsonite.Unmarshal(data, v)
}

// UnmarshalFromString decodes the json string into target object
func UnmarshalFromString(str string, v interface{}) error {
	return jsonite.UnmarshalFromString(str, v)
}

// MarshalToString encodes object into json string
func MarshalToString(v interface{}) (string, error) {
	return jsonite.MarshalToString(v)
}

// Marshal encodes object into json byte array
func Marshal(v interface{}) ([]byte, error) {
	return jsonite.Marshal(v)
}

// MarshalIndent encodes object into json byte array with indent
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsonite.MarshalIndent(v, prefix, indent)
}

// Get gets the value of json with specified path
func Get(data []byte, path ...interface{}) jsoniter.Any {
	return jsonite.Get(data, path...)
}

// Valid is used to check whether the json data is valid
func Valid(data []byte) bool {
	return jsonite.Valid(data)
}
