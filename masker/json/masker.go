package masker

import (
	"encoding/base64"
	"errors"
	"fmt"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/buger/jsonparser"
	"math"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

// defines all supported mask type
const (
	Overlay  = "overlay"
	Base64   = "base64"
	Password = "password"
	Empty    = "empty"
)

type MaskOption func(*MaskOptions)

type MaskOptions struct {
	// OverlayString means the string being replaced, default value is "*"
	OverlayString string

	// Percent means the percent of string length will be overlay, default value is 0
	// if start/end not setted, will set to 0.5
	Percent float64

	// Start means the start index of overlay, default value is -1
	Start int

	// End means the end index of overlay, default value is -1
	End int
}

func init() {
	RegisterMasker(Overlay, NewOverlayMasker())
	RegisterMasker(Base64, NewBase64Masker())
	RegisterMasker(Password, NewPasswordMasker())
	RegisterMasker(Empty, NewEmptyMasker())
}

type MaskerInc interface {
	Do(keyPath, in string, parameters ...string) (string, error)
}

var maskers = make(map[string]MaskerInc, 0)
var locker sync.RWMutex

func GetMasker(key string) (MaskerInc, bool) {
	locker.RLock()
	defer locker.RUnlock()
	finalKey := strings.ToLower(key)

	res, ok := maskers[finalKey]
	return res, ok
}

// RegisterMaskFunc is used to register mask function into masker
func RegisterMasker(key string, masker MaskerInc) {
	locker.Lock()
	defer locker.Unlock()
	finalKey := strings.ToLower(key)

	maskers[finalKey] = masker
}

func overlay(str string, overlay string, start int, end int) (overlayed string) {
	r := []rune(str)
	l := len(r)

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l-1 {
		end = l - 1
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	if end+1 <= l {
		overlayed += overlay
		overlayed += string(r[end+1:])
	}
	return overlayed
}

type emptyMasker struct{}

func (p *emptyMasker) Do(_, in string, parameters ...string) (string, error) {
	// Do nothing, just return empty string
	return "", nil
}

type passwordMasker struct{}

func (p *passwordMasker) Do(_, in string, parameters ...string) (string, error) {
	opt := defaultPasswordMaskerOptions()

	return strings.Repeat(opt.OverlayString, len(in)), nil
}

type base64Masker struct{}

func (b *base64Masker) Do(_, in string, parameters ...string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(in)), nil
}

type overlayMasker struct{}

func defaultOverlayMaskerOptions() MaskOptions {
	return MaskOptions{
		OverlayString: constant.DefaultOverlayString,
	}
}

func defaultPasswordMaskerOptions() MaskOptions {
	return MaskOptions{
		OverlayString: constant.DefaultOverlayString,
	}
}

func (o *overlayMasker) Do(keyPath, in string, parameters ...string) (string, error) {
	if len(in) == 0 {
		return in, nil
	}
	opt := defaultOverlayMaskerOptions()
	var maskOptions []MaskOption
	if len(parameters) > 0 {
		parameterString := parameters[0]
		parameterFields := strings.Split(parameterString, ",")
		if len(parameterFields) == 1 {
			percent, gerr := strconv.ParseFloat(parameterFields[0], 64)
			if nil != gerr {
				log.Errorsf("Failed to parse string into float, key path:[%s], error:%++v", keyPath, gerr)
				return in, gerr
			}
			maskOptions = append(maskOptions, func(options *MaskOptions) {
				options.Percent = percent
			})
		} else if len(parameterFields) > 1 {
			start, gerr := strconv.Atoi(parameterFields[0])
			if nil != gerr {
				return "", gerr
			}

			end, gerr := strconv.Atoi(parameterFields[1])
			if nil != gerr {
				return "", gerr
			}
			maskOptions = append(maskOptions, func(options *MaskOptions) {
				options.Start = start
				options.End = end
			})
		}
	}

	for _, o := range maskOptions {
		o(&opt)
	}
	inLen := utf8.RuneCountInString(in)
	if opt.Percent > 0 {
		lenOfHeadAndTail := int((inLen - int(math.Floor(float64(inLen)*opt.Percent))) / 2)
		indexOfEnd := inLen - lenOfHeadAndTail

		return overlay(in, opt.OverlayString, lenOfHeadAndTail, indexOfEnd), nil
	}

	if opt.Start < 0 || opt.Start >= opt.End || opt.Start >= inLen {
		return in, nil
	}

	if opt.End > inLen {
		return overlay(in, opt.OverlayString, opt.Start, inLen), nil
	}

	return overlay(in, opt.OverlayString, opt.Start, opt.End), nil
}

func NewBase64Masker() *base64Masker {
	return &base64Masker{}
}

func NewPasswordMasker() *passwordMasker {
	return &passwordMasker{}
}

func NewOverlayMasker() *overlayMasker {
	return &overlayMasker{}
}

func NewEmptyMasker() *emptyMasker {
	return &emptyMasker{}
}

type Node struct {
	key      string
	value    []byte
	index    int
	last     *Node
	dataType jsonparser.ValueType
}

type JsonObject struct {
	leafs []*Node
}

func NewJsonObject() *JsonObject {
	return &JsonObject{}
}

var NotSupportTypeError = errors.New("Not support type error")

func (obj *JsonObject) scan(data []byte, keys []string, index int, node *Node) error {
	if index >= len(keys) {
		obj.leafs = append(obj.leafs, node.last)
		return nil
	}
	node.key = keys[index]

	searchKeys := keys[:index+1]
	value, dataType, _, err := jsonparser.Get(data, searchKeys...)
	if err != nil {
		return err
	}

	if jsonparser.String != dataType &&
		jsonparser.Array != dataType &&
		jsonparser.Object != dataType {
		return NotSupportTypeError
	}

	if node.key == keys[len(keys)-1] {
		node.dataType = dataType
		node.value = value
	}
	if dataType == jsonparser.Array {
		arrIndex := 0
		_, _ = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			next := &Node{
				last:  node,
				index: arrIndex,
			}
			_ = obj.scan(value, keys[index+1:], 0, next)
			arrIndex++
		}, searchKeys...)
		return nil
	} else {
		next := &Node{
			last:  node,
			index: -1,
		}
		return obj.scan(data, keys, index+1, next)
	}
}

func (obj JsonObject) toKeys() [][]string {
	pathes := make([][]string, 0)
	for i := 0; i < len(obj.leafs); i++ {
		path := make([]string, 0)
		y := obj.leafs[i]
		for y != nil {
			path = append(path, y.key)
			if y.index >= 0 {
				path = append(path, "["+strconv.Itoa(y.index)+"]")
			}
			y = y.last
		}

		for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
			path[left], path[right] = path[right], path[left]
		}
		pathes = append(pathes, path)
	}
	return pathes
}

func MapValueMask(source map[string]string, maskRules []string, returnErrorIfMissingField ...bool) (map[string]string, bool, error) {
	isExistError := false
	finalMap := make(map[string]string, 0)
	for k, v := range source {
		finalMap[k] = v
	}

	for _, rule := range maskRules {
		fields := strings.Split(rule, "|")
		if len(fields) >= 2 {
			keyPath := fields[0]
			maskerName := strings.ToLower(fields[1])

			var parameters []string
			if len(fields) >= 3 {
				parameters = make([]string, 0)
				parameters = append(parameters, fields[2])
			}
			masker, ok := GetMasker(maskerName)
			if !ok {
				log.Errorsf("Cannot found the masker with name,key path:[%s], masker name:%s", keyPath, maskerName)
				isExistError = true
				continue
			}

			value, ok := finalMap[keyPath]
			if !ok && len(returnErrorIfMissingField) > 0 && returnErrorIfMissingField[0] {
				log.Errorsf("Failed to scan, key path:[%s]", keyPath)
				return nil, true, errors.New(fmt.Sprintf("failed to scan, key path:[%s]", keyPath))
			}

			finalValue, gerr := masker.Do(keyPath, value, parameters...)
			if nil != gerr {
				log.Errorsf("Failed to do the masking, key path:[%s], error:%++v", keyPath, gerr)
				isExistError = true
				continue
			}

			finalMap[keyPath] = finalValue
		}

	}

	return finalMap, isExistError, nil
}

func JsonBodyMask(source []byte, maskRules []string, returnErrorIfMissingField ...bool) ([]byte, bool, error) {
	finalResult := source
	isExistError := false
	// rule format:
	//    key | masker name | parameters
	for _, rule := range maskRules {
		fields := strings.Split(rule, "|")
		if len(fields) >= 2 {
			keyPath := fields[0]
			maskerName := strings.ToLower(fields[1])

			var parameters []string
			if len(fields) >= 3 {
				parameters = make([]string, 0)
				parameters = append(parameters, fields[2])
			}
			masker, ok := GetMasker(maskerName)
			if !ok {
				log.Errorsf("Cannot found the masker with name,key path:[%s], masker name:%s", keyPath, maskerName)
				isExistError = true
				continue
			}
			finalKeyPath := strings.ReplaceAll(keyPath, "[*]", "")
			keys := strings.Split(finalKeyPath, ".")

			jsonObject := NewJsonObject()
			root := &Node{
				last:  nil,
				index: -1,
			}
			err := jsonObject.scan(source, keys, 0, root)
			if err != nil {
				if len(returnErrorIfMissingField) > 0 && returnErrorIfMissingField[0] {
					log.Errorsf("Failed to scan, key path:[%s] error:%++v", keyPath, err)
					return nil, true, errors.New(fmt.Sprintf("failed to scan, key path:[%s] error:%++v", keyPath, err))
				}
				//log.Errorsf("Failed to scan, error:%++v", err)
				continue
			}

			converted := make(map[string]bool, 0)
			for i := 0; i < len(jsonObject.leafs); i++ {
				path := make([]string, 0)
				leaf := jsonObject.leafs[i]
				if len(leaf.value) == 0 {
					continue
				}
				node := jsonObject.leafs[i]
				for node != nil {
					path = append(path, node.key)
					if node.index >= 0 {
						path = append(path, "["+strconv.Itoa(node.index)+"]")
					}
					node = node.last
				}

				for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
					path[left], path[right] = path[right], path[left]
				}

				pathStr := strings.Join(path, ",")
				if _, ok := converted[pathStr]; ok {
					continue
				}

				var convertedValue string
				switch leaf.dataType {
				case jsonparser.Array:
					values := strings.Split(string(leaf.value[1:len(leaf.value)-1]), ",")
					convertedValue = string(constant.ArrPreffix)
					for i, value := range values {
						value = strings.TrimSpace(value)
						r, gerr := masker.Do(keyPath, value, parameters...)
						if nil != gerr {
							log.Errorsf("Failed to do the masking, key path:[%s], error:%++v", keyPath, gerr)
							isExistError = true
							continue
						}
						convertedValue += r
						if i < len(values)-1 {
							convertedValue += string(constant.Seperator)
						}
					}
					convertedValue += string(constant.ArrSuffix)
				case jsonparser.Object:
					convertedValue = "{*}"
				default:
					r, gerr := masker.Do(keyPath, string(leaf.value), parameters...)
					if nil != gerr {
						log.Errorsf("Failed to do the masking, key path:[%s], error:%++v", keyPath, gerr)
						isExistError = true
						continue
					}
					convertedValue = "\"" + r + "\""
				}
				finalResult, err = jsonparser.Set(finalResult, []byte(convertedValue), path...)
				if err != nil {
					return source, isExistError, err
				}
			}
		}
	}

	return finalResult, isExistError, nil
}
