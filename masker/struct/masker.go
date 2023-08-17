package masker

import (
	"reflect"
	"sync"
)

const tagName = "mask"

// MaskType is an alias of mask type
type MaskType string

// defines all supported mask type
const (
	MStruct MaskType = "struct"
)

var maskers = make(map[MaskType]func(string) string, 0)
var locker sync.RWMutex

// Struct performs the mask logic
func Struct(s interface{}) interface{} {
	if s == nil {
		return nil
	}

	var selem, tptr reflect.Value

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {
		tptr = reflect.New(st.Elem())
		selem = reflect.ValueOf(s).Elem()
	} else {
		tptr = reflect.New(st)
		selem = reflect.ValueOf(s)
	}

	for i := 0; i < selem.NumField(); i++ {
		mtag := selem.Type().Field(i).Tag.Get(tagName)
		if len(mtag) == 0 {
			tptr.Elem().Field(i).Set(selem.Field(i))
			continue
		}
		switch selem.Field(i).Type().Kind() {
		default:
			tptr.Elem().Field(i).Set(selem.Field(i))
		case reflect.String:
			tptr.Elem().Field(i).SetString(String(MaskType(mtag), selem.Field(i).String()))
		case reflect.Struct:
			if MaskType(mtag) == MStruct {
				_t := Struct(selem.Field(i).Interface())
				if _t == nil {
					return nil
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			}
		case reflect.Ptr:
			if selem.Field(i).IsNil() {
				continue
			}
			if MaskType(mtag) == MStruct {
				_t := Struct(selem.Field(i).Interface())
				if _t == nil {
					return nil
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		case reflect.Slice:
			if selem.Field(i).IsNil() {
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.String {
				orgval := selem.Field(i).Interface().([]string)
				newval := make([]string, len(orgval))
				for i, val := range selem.Field(i).Interface().([]string) {
					newval[i] = String(MaskType(mtag), val)
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(newval))
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Struct && MaskType(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n := Struct(selem.Field(i).Index(j).Interface())
					if _n == nil {
						return nil
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Ptr && MaskType(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n := Struct(selem.Field(i).Index(j).Interface())
					if _n == nil {
						return nil
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n))
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Interface && MaskType(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n := Struct(selem.Field(i).Index(j).Interface())
					if _n == nil {
						return nil
					}
					if reflect.TypeOf(selem.Field(i).Index(j).Interface()).Kind() != reflect.Ptr {
						newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
					} else {
						newval = reflect.Append(newval, reflect.ValueOf(_n))
					}
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
		case reflect.Interface:
			if selem.Field(i).IsNil() {
				continue
			}
			if MaskType(mtag) != MStruct {
				continue
			}
			_t := Struct(selem.Field(i).Interface())
			if _t == nil {
				return nil
			}
			if reflect.TypeOf(selem.Field(i).Interface()).Kind() != reflect.Ptr {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			} else {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		}
	}

	return tptr.Interface()
}

// String is a mask function for string
func String(t MaskType, i string) string {
	if fn, ok := maskers[t]; ok {
		return fn(i)
	}

	return i
}

// RegisterMaskFunc is used to register mask function into masker
func RegisterMaskFunc(key MaskType, fn func(string) string) {
	locker.Lock()
	defer locker.Unlock()
	maskers[key] = fn
}

// Do is an alias function name of Struct
func Do(s interface{}) interface{} {
	return Struct(s)
}
