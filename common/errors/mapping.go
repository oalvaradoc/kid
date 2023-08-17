package errors

import "strings"

type errorCodeMapping struct  {
	cache map[string]string
}

var _cache = &errorCodeMapping{
	cache: map[string]string{},
}

func SetErrorCodeMapping(mapping map[string]string) {
	if nil != mapping {
		for k,v := range mapping {
			_cache.cache[k] =v
		}
	}
}

func GetFinalErrorCode(code string) string {
	if len(_cache.cache) > 0 {
		if v, ok := _cache.cache[strings.ToLower(code)]; ok {
			return v
		}
	}
	return code
}
