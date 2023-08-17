package util

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lpadZeros = [...]string{
	"",                     // 0
	"0",                    // 1
	"00",                   // 2
	"000",                  // 3
	"0000",                 // 4
	"00000",                // 5
	"000000",               // 6
	"0000000",              // 7
	"00000000",             // 8
	"000000000",            // 9
	"0000000000",           // 10
	"00000000000",          // 11
	"000000000000",         // 12
	"0000000000000",        // 13
	"00000000000000",       // 14
	"000000000000000",      // 15
	"0000000000000000",     // 16
	"00000000000000000",    // 17
	"000000000000000000",   // 18
	"0000000000000000000",  // 19
	"00000000000000000000", // 20
}

func fmtDigital(targetLength int, srcDigital int) string {
	str := strconv.Itoa(srcDigital)
	lengthOfStr := len(str)
	if targetLength > lengthOfStr {
		return lpadZeros[targetLength-lengthOfStr] + str
	}

	return str
}

// IndexOfItem is used to determine the subscript position of the value in the array
func IndexOfItem(value interface{}, array interface{}) int {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return i
			}
		}
	}

	return -1
}

// CurrentTime is used to get the formatted current time
func CurrentTime() string {
	t := time.Now()
	return FormatTime(t)
}

// FormatTime is used to convert the time to string
func FormatTime(t time.Time) string {
	// The format of following string is "%04d-%02d-%02d %02d:%02d:%02d.%09d"
	str := fmtDigital(4, t.Year()) + "-" +
		fmtDigital(2, int(t.Month())) + "-" +
		fmtDigital(2, t.Day()) + " " +
		fmtDigital(2, t.Hour()) + ":" +
		fmtDigital(2, t.Minute()) + ":" +
		fmtDigital(2, t.Second()) + "." +
		fmtDigital(9, t.Nanosecond())
	return str
}

// ToTime is used to parse string to time
func ToTime(s string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, s, loc)
	return theTime
}

// CurrentHost is used to get the IP of machine where the service is located.
func CurrentHost() (host string) {
	host = "localhost"
	netInterfaces, e := net.Interfaces()
	if e != nil {
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					host = ipnet.IP.String()
					return
				}
			}
		}
	}

	return
}

// Fnv32 is a tool function that hashing string to a uint32 value
func Fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// IntToBytes is used to encode an int value to byte array
func IntToBytes(x int) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, int32(x)); err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}

// BytesToInt is used to convert byte array to an int value
func BytesToInt(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return int(x), nil
}

var flagNo int64 = 0
var lock sync.Mutex

const seqMax = 999999

// GenerateSerialNo is used to generate an unique serial number
// org: organization
// wks: workspace
// env: environment
// suNo: SU number
// instanceID: instance ID
// tp : 0 - trace ID
//      1 - span ID
func GenerateSerialNo(org, wks, env, suNo, instanceID string, tp string) string {
	timestamp := time.Now().Unix()

	serialNo := tp + org + wks + env + suNo + instanceID + strconv.FormatInt(timestamp, 10)

	lock.Lock()
	defer lock.Unlock()
	flagNo++
	if flagNo > seqMax {
		flagNo = 1
	}
	flagNo2Str := strconv.FormatInt(flagNo, 10)
	flagNo2Str = lpadZeros[6-len(flagNo2Str)] + flagNo2Str
	serialNo = serialNo + flagNo2Str
	return serialNo
}

func getASCII() string {
	str := ""
	for i := 32; i <= 126; i++ {
		str += string(byte(i))
	}
	return str
}

var ascii = getASCII()

// RandomString is used to generate a random string number
func RandomString(l int) string {
	bytes := []byte(ascii)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetMapValueIgnoreCase is used to get map value with key ignore case
func GetMapValueIgnoreCase(m map[string]string, key string) string {
	if nil == m {
		return ""
	}
	value := ""
	if v, ok := m[key]; ok {
		value = v
	} else if v, ok := m[strings.ToLower(key)]; ok {
		value = v
	} else if v, ok := m[strings.ToUpper(key)]; ok {
		value = v
	}

	return value
}
