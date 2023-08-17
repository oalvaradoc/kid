package v2

import "time"

// SkmPublicPemRequest is a request modle to stores the request of get public key of SKM communction
type SkmPublicPemRequest struct {
	Algorithm string `json:"algorithm"` //rsa,sm2
}

// SkmPublicPemResponseData is a response model to stores the response of get public key of SKM communction
type SkmPublicPemResponseData struct {
	PublicPem string `json:"publicPem"`
}

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

type Date time.Time

func Now() Date {
	return Date(time.Now())
}

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Date(now)
	return
}

func (t Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Date) String() string {
	return time.Time(t).Format(TimeFormat)
}

// KeyType is alias of string
type KeyType string

// SkmKey is used to store the details of each key
type SkmKey struct {
	ServiceId  string  `json:"serviceId"`
	SystemType string  `json:"systemType"`
	KeyType    KeyType `json:"keyType"`
	KeyId      string  `json:"keyId"`
	KeyVersion string  `json:"keyVersion"`
	KeyValue   string  `json:"keyValue"`
	State      int     `json:"state"`
	EffectTime Date    `json:"effectTime"`
	ExpireTime Date    `json:"expireTime"`
	Creator    string  `json:"creator"`
	CreateTime Date    `json:"createTime"`
	Modifier   string  `json:"modifier"`
	LastUpdate Date    `json:"lastUpdate"`
}

// A type, typically a collection, that satisfies sort.Interface can be
// sorted by the routines in this package. The methods require that the
// elements of the collection be enumerated by an integer index.
type SkmKeySlice []SkmKey

// Len is the number of elements in the collection.
func (s SkmKeySlice) Len() int { return len(s) }

// Swap swaps the elements with indexes i and j.
func (s SkmKeySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SkmKeySlice) Less(i, j int) bool {
	return time.Time(s[i].EffectTime).Before(time.Time(s[j].EffectTime))
}

// SkmKeyResult is used to store key information
type SkmKeyResult struct {
	ServiceId string      `json:"serviceId"`
	Result    SkmKeySlice `json:"result"`
}

// SkmServiceRequest is used to store the request structure that the client requests
// to obtain a single service key information from the SKM adapter.
type SkmServiceRequest struct {
	KeyType     KeyType `json:"keyType"`
	ServiceId   string  `json:"serviceId"`
	SystemType  string  `json:"systemType"`
	EffectTime  Date    `json:"effectTime"`
	ExpireTime  Date    `json:"expireTime"`
	Operator    string  `json:"operator"`
	RequestTime Date    `json:"requestTime"`
}

// A request structure for storing a client requesting a SKM adapter
// to obtain a plurality of different key types and different service key information.
type SkmKeyValuesRequest struct {
	Services []SkmServiceRequest `json:"services"`
}

// SkmKeysResponse is used to store multiple service key information
type SkmKeysResponse struct {
	//Total  int              `json:"total"`
	//Succ   int              `json:"succ"`
	Result []SkmKeyResult `json:"result"`
}

// SkmRequestCrypto for all request that need to be encrypted.
type SkmRequestCrypto struct {
	Request         string `json:"request"`
	Key             string `json:"key"`
	CryptoAlgorithm string `json:"cryptoAlgorithm"`
}

// SkmResponseCrypto for all response that need to be decrypted
type SkmResponseCrypto struct {
	Response string `json:"response"`
}
