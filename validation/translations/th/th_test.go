package th

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {

	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("th")

	validate := validator.New()

	err := RegisterDefaultTranslations(validate, trans)
	Equal(t, err, nil)

	type Inner struct {
		EqCSFieldString  string
		NeCSFieldString  string
		GtCSFieldString  string
		GteCSFieldString string
		LtCSFieldString  string
		LteCSFieldString string
	}

	type Test struct {
		Inner             Inner
		RequiredString    string    `validate:"required"`
		RequiredNumber    int       `validate:"required"`
		RequiredMultiple  []string  `validate:"required"`
		LenString         string    `validate:"len=1"`
		LenNumber         float64   `validate:"len=1113.00"`
		LenMultiple       []string  `validate:"len=7"`
		MinString         string    `validate:"min=1"`
		MinNumber         float64   `validate:"min=1113.00"`
		MinMultiple       []string  `validate:"min=7"`
		MaxString         string    `validate:"max=3"`
		MaxNumber         float64   `validate:"max=1113.00"`
		MaxMultiple       []string  `validate:"max=7"`
		EqString          string    `validate:"eq=3"`
		EqNumber          float64   `validate:"eq=2.33"`
		EqMultiple        []string  `validate:"eq=7"`
		NeString          string    `validate:"ne="`
		NeNumber          float64   `validate:"ne=0.00"`
		NeMultiple        []string  `validate:"ne=0"`
		LtString          string    `validate:"lt=3"`
		LtNumber          float64   `validate:"lt=5.56"`
		LtMultiple        []string  `validate:"lt=2"`
		LtTime            time.Time `validate:"lt"`
		LteString         string    `validate:"lte=3"`
		LteNumber         float64   `validate:"lte=5.56"`
		LteMultiple       []string  `validate:"lte=2"`
		LteTime           time.Time `validate:"lte"`
		GtString          string    `validate:"gt=3"`
		GtNumber          float64   `validate:"gt=5.56"`
		GtMultiple        []string  `validate:"gt=2"`
		GtTime            time.Time `validate:"gt"`
		GteString         string    `validate:"gte=3"`
		GteNumber         float64   `validate:"gte=5.56"`
		GteMultiple       []string  `validate:"gte=2"`
		GteTime           time.Time `validate:"gte"`
		EqFieldString     string    `validate:"eqfield=MaxString"`
		EqCSFieldString   string    `validate:"eqcsfield=Inner.EqCSFieldString"`
		NeCSFieldString   string    `validate:"necsfield=Inner.NeCSFieldString"`
		GtCSFieldString   string    `validate:"gtcsfield=Inner.GtCSFieldString"`
		GteCSFieldString  string    `validate:"gtecsfield=Inner.GteCSFieldString"`
		LtCSFieldString   string    `validate:"ltcsfield=Inner.LtCSFieldString"`
		LteCSFieldString  string    `validate:"ltecsfield=Inner.LteCSFieldString"`
		NeFieldString     string    `validate:"nefield=EqFieldString"`
		GtFieldString     string    `validate:"gtfield=MaxString"`
		GteFieldString    string    `validate:"gtefield=MaxString"`
		LtFieldString     string    `validate:"ltfield=MaxString"`
		LteFieldString    string    `validate:"ltefield=MaxString"`
		AlphaString       string    `validate:"alpha"`
		AlphanumString    string    `validate:"alphanum"`
		NumericString     string    `validate:"numeric"`
		NumberString      string    `validate:"number"`
		HexadecimalString string    `validate:"hexadecimal"`
		HexColorString    string    `validate:"hexcolor"`
		RGBColorString    string    `validate:"rgb"`
		RGBAColorString   string    `validate:"rgba"`
		HSLColorString    string    `validate:"hsl"`
		HSLAColorString   string    `validate:"hsla"`
		Email             string    `validate:"email"`
		URL               string    `validate:"url"`
		URI               string    `validate:"uri"`
		Base64            string    `validate:"base64"`
		Contains          string    `validate:"contains=purpose"`
		ContainsAny       string    `validate:"containsany=!@#$"`
		Excludes          string    `validate:"excludes=text"`
		ExcludesAll       string    `validate:"excludesall=!@#$"`
		ExcludesRune      string    `validate:"excludesrune=☻"`
		ISBN              string    `validate:"isbn"`
		ISBN10            string    `validate:"isbn10"`
		ISBN13            string    `validate:"isbn13"`
		UUID              string    `validate:"uuid"`
		UUID3             string    `validate:"uuid3"`
		UUID4             string    `validate:"uuid4"`
		UUID5             string    `validate:"uuid5"`
		ASCII             string    `validate:"ascii"`
		PrintableASCII    string    `validate:"printascii"`
		MultiByte         string    `validate:"multibyte"`
		DataURI           string    `validate:"datauri"`
		Latitude          string    `validate:"latitude"`
		Longitude         string    `validate:"longitude"`
		SSN               string    `validate:"ssn"`
		IP                string    `validate:"ip"`
		IPv4              string    `validate:"ipv4"`
		IPv6              string    `validate:"ipv6"`
		CIDR              string    `validate:"cidr"`
		CIDRv4            string    `validate:"cidrv4"`
		CIDRv6            string    `validate:"cidrv6"`
		TCPAddr           string    `validate:"tcp_addr"`
		TCPAddrv4         string    `validate:"tcp4_addr"`
		TCPAddrv6         string    `validate:"tcp6_addr"`
		UDPAddr           string    `validate:"udp_addr"`
		UDPAddrv4         string    `validate:"udp4_addr"`
		UDPAddrv6         string    `validate:"udp6_addr"`
		IPAddr            string    `validate:"ip_addr"`
		IPAddrv4          string    `validate:"ip4_addr"`
		IPAddrv6          string    `validate:"ip6_addr"`
		UinxAddr          string    `validate:"unix_addr"` // can't fail from within Go's net package currently, but maybe in the future
		MAC               string    `validate:"mac"`
		IsColor           string    `validate:"iscolor"`
		StrPtrMinLen      *string   `validate:"min=10"`
		StrPtrMaxLen      *string   `validate:"max=1"`
		StrPtrLen         *string   `validate:"len=2"`
		StrPtrLt          *string   `validate:"lt=1"`
		StrPtrLte         *string   `validate:"lte=1"`
		StrPtrGt          *string   `validate:"gt=10"`
		StrPtrGte         *string   `validate:"gte=10"`
		OneOfString       string    `validate:"oneof=red green"`
		OneOfInt          int       `validate:"oneof=5 63"`
		JSONString        string    `validate:"json"`
		LowercaseString   string    `validate:"lowercase"`
		UppercaseString   string    `validate:"uppercase"`
		Datetime          string    `validate:"datetime=2006-01-02"`
	}

	var test Test

	test.Inner.EqCSFieldString = "1234"
	test.Inner.GtCSFieldString = "1234"
	test.Inner.GteCSFieldString = "1234"

	test.MaxString = "1234"
	test.MaxNumber = 2000
	test.MaxMultiple = make([]string, 9)

	test.LtString = "1234"
	test.LtNumber = 6
	test.LtMultiple = make([]string, 3)
	test.LtTime = time.Now().Add(time.Hour * 24)

	test.LteString = "1234"
	test.LteNumber = 6
	test.LteMultiple = make([]string, 3)
	test.LteTime = time.Now().Add(time.Hour * 24)

	test.LtFieldString = "12345"
	test.LteFieldString = "12345"

	test.LtCSFieldString = "1234"
	test.LteCSFieldString = "1234"

	test.AlphaString = "abc3"
	test.AlphanumString = "abc3!"
	test.NumericString = "12E.00"
	test.NumberString = "12E"

	test.Excludes = "this is some test text"
	test.ExcludesAll = "This is Great!"
	test.ExcludesRune = "Love it ☻"

	test.ASCII = "ｶﾀｶﾅ"
	test.PrintableASCII = "ｶﾀｶﾅ"

	test.MultiByte = "1234feerf"

	s := "toolong"
	test.StrPtrMaxLen = &s
	test.StrPtrLen = &s

	test.JSONString = "{\"foo\":\"bar\",}"

	test.LowercaseString = "ABCDEFG"
	test.UppercaseString = "abcdefg"

	test.Datetime = "20060102"

	err = validate.Struct(test)
	NotEqual(t, err, nil)

	errs, ok := err.(validator.ValidationErrors)
	Equal(t, ok, true)

	tests := []struct {
		ns       string
		expected string
	}{
		{
			ns:       "Test.IsColor",
			expected: "IsColor ต้องเป็นสีที่ถูกต้อง",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC ต้องเป็นที่อยู่ MAC ที่ถูกต้อง",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr ต้องเป็นที่อยู่ IP ที่ถูกต้อง",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4 ต้องเป็นที่อยู่ IPv4 ที่ถูกต้อง",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6 ต้องเป็นที่อยู่ IPv6 ที่ถูกต้อง",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr ต้องเป็นที่อยู่ UDP ที่ถูกต้อง",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4 ต้องเป็นที่อยู่ IPv4 UDP ที่ถูกต้อง",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6 ต้องเป็นที่อยู่ IPv6 UDP ที่ถูกต้อง",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr ต้องเป็นที่อยู่ TCP ที่ถูกต้อง",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4 ต้องเป็นที่อยู่ IPv4 TCP ที่ถูกต้อง",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6 ต้องเป็นที่อยู่ IPv6 TCP ที่ถูกต้อง",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR ต้องเป็นเส้นทางระหว่างโดเมนแบบไม่มีคลาส (CIDR) ที่ถูกต้อง",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4 ต้องเป็นเส้นทางระหว่างโดเมน (CIDR) แบบไม่มีคลาสที่ถูกต้องซึ่งมีที่อยู่ IPv4",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6 ต้องเป็นเส้นทางระหว่างโดเมน (CIDR) แบบไม่มีคลาสที่ถูกต้องซึ่งมีที่อยู่ IPv6",
		},
		{
			ns:       "Test.SSN",
			expected: "SSNต้องเป็นหมายเลขประกันสังคมที่ถูกต้อง(SSN)",
		},
		{
			ns:       "Test.IP",
			expected: "IPต้องเป็นที่อยู่ IP ที่ถูกต้อง",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4ต้องเป็นที่อยู่ IPv4 ที่ถูกต้อง",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6ต้องเป็นที่อยู่ IPv6 ที่ถูกต้อง",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURIต้องมีข้อมูลที่ถูกต้องURI",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitudeต้องมีพิกัดละติจูดที่ถูกต้อง",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitudeต้องมีพิกัดลองจิจูดที่ถูกต้อง",
		},
		{
			ns:       "Test.MultiByte",
			expected: "MultiByteต้องมีอักขระหลายไบต์",
		},
		{
			ns:       "Test.ASCII",
			expected: "ASCIIต้องมีอักขระ ascii เท่านั้น",
		},
		{
			ns:       "Test.PrintableASCII",
			expected: "PrintableASCIIต้องมีอักขระ ascii ที่พิมพ์ได้เท่านั้น",
		},
		{
			ns:       "Test.UUID",
			expected: "UUIDต้องเป็นไฟล์UUID",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3ต้องเป็นไฟล์V3 UUID",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4ต้องเป็น V4 ที่ถูกต้อง UUID",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5ต้องเป็น V5 ที่ถูกต้อง UUID",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBNต้องเป็นหมายเลข ISBN ที่ถูกต้อง",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10ต้องเป็นหมายเลข ISBN-10 ที่ถูกต้อง",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13ต้องเป็นหมายเลข ISBN-13 ที่ถูกต้อง",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludesไม่สามารถมีข้อความได้'text'",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAllต้องไม่มีอักขระใด ๆ ต่อไปนี้'!@#$'",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRuneไม่สามารถมี'☻'",
		},
		{
			ns:       "Test.ContainsAny",
			expected: "ContainsAnyต้องมีอักขระต่อไปนี้อย่างน้อยหนึ่งตัว'!@#$'",
		},
		{
			ns:       "Test.Contains",
			expected: "Containsต้องมีข้อความ'purpose'",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64ต้องเป็นสตริง Base64 ที่ถูกต้อง",
		},
		{
			ns:       "Test.Email",
			expected: "Emailต้องเป็นกล่องจดหมายที่ถูกต้อง",
		},
		{
			ns:       "Test.URL",
			expected: "URLต้องเป็น URL ที่ถูกต้อง",
		},
		{
			ns:       "Test.URI",
			expected: "URIต้องเป็น URI ที่ถูกต้อง",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorStringต้องเป็นสี RGB ที่ถูกต้อง",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorStringต้องเป็นสี RGBA ที่ถูกต้อง",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorStringต้องเป็นสี HSL ที่ถูกต้อง",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorStringต้องเป็นสี HSLA ที่ถูกต้อง",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalStringต้องเป็นเลขฐานสิบหกที่ถูกต้อง",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorStringต้องเป็นเลขฐานสิบหกที่ถูกต้องสี",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberStringต้องเป็นตัวเลขที่ถูกต้อง",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericStringต้องเป็นค่าที่ถูกต้อง",
		},
		{
			ns:       "Test.AlphanumString",
			expected: "AlphanumStringมีได้เฉพาะตัวอักษรและตัวเลข",
		},
		{
			ns:       "Test.AlphaString",
			expected: "AlphaStringมีได้เฉพาะตัวอักษร",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldStringต้องน้อยกว่าMaxString",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldStringต้องน้อยกว่าหรือเท่ากับMaxString",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldStringต้องมากกว่าMaxString",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldStringต้องมากกว่าหรือเท่ากับMaxString",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldStringไม่สามารถเท่ากับEqFieldString",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldStringต้องน้อยกว่าInner.LtCSFieldString",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldStringต้องน้อยกว่าหรือเท่ากับInner.LteCSFieldString",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldStringต้องมากกว่าInner.GtCSFieldString",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldStringต้องมากกว่าหรือเท่ากับInner.GteCSFieldString",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldStringไม่สามารถเท่ากับInner.NeCSFieldString",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldStringต้องเท่ากับInner.EqCSFieldString",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldStringต้องเท่ากับMaxString",
		},
		{
			ns:       "Test.GteString",
			expected: "GteString ต้องมีความยาวอย่างน้อย 3 อักขระ",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber ต้องมากกว่าหรือเท่ากับ 5.56",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple ต้องมีอย่างน้อย 2 รายการ",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime ต้องมากกว่าหรือเท่ากับวันที่และเวลาปัจจุบัน",
		},
		{
			ns:       "Test.GtString",
			expected: "GtString ต้องมากกว่า 3 อักขระ",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber ต้องมากกว่า 5.56",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple ต้องมากกว่า 2 รายการ",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime ต้องมากกว่าวันที่และเวลาปัจจุบัน",
		},
		{
			ns:       "Test.LteString",
			expected: "LteString ความยาวต้องไม่เกิน 3 อักขระ",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber ต้องน้อยกว่าหรือเท่ากับ 5.56",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple สามารถมีได้มากที่สุด 2 รายการ",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime ต้องน้อยกว่าหรือเท่ากับวันที่และเวลาปัจจุบัน",
		},
		{
			ns:       "Test.LtString",
			expected: "LtString ต้องมีความยาวน้อยกว่า 3 อักขระ",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber ต้องน้อยกว่า 5.56",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple ต้องมีน้อยกว่า 2 รายการ",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime ต้องน้อยกว่าวันที่และเวลาปัจจุบัน",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString ต้องไม่เท่ากับ ",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber ต้องไม่เท่ากับ 0.00",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple ต้องไม่เท่ากับ 0",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString ไม่เท่ากับ 3",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber ไม่เท่ากับ 2.33",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple ไม่เท่ากับ 7",
		},
		{
			ns:       "Test.MaxString",
			expected: "MaxString ความยาวต้องไม่เกิน 3 อักขระ",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber ต้องน้อยกว่าหรือเท่ากับ 1,113.00",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple สามารถมีได้มากที่สุด 7 รายการ",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString ต้องมีความยาวอย่างน้อย 1 อักขระ",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber ขั้นต่ำสามารถทำได้เพียง 1,113.00",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple ต้องมีอย่างน้อย 7 รายการ",
		},
		{
			ns:       "Test.LenString",
			expected: "LenStringความยาวต้องเป็น1อักขระ",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumberต้องเท่ากับ1,113.00",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultipleต้องเท่ากับ7สิ่งของ",
		},
		{
			ns:       "Test.RequiredString",
			expected: "RequiredStringเป็นฟิลด์บังคับ",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "RequiredNumberเป็นฟิลด์บังคับ",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "RequiredMultipleเป็นฟิลด์บังคับ",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen ต้องมีความยาวอย่างน้อย 10 อักขระ",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen ความยาวต้องไม่เกิน 1 อักขระ",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLenความยาวต้องเป็น2อักขระ",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt ต้องมีความยาวน้อยกว่า 1 อักขระ",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte ความยาวต้องไม่เกิน 1 อักขระ",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt ต้องมากกว่า 10 อักขระ",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte ต้องมีความยาวอย่างน้อย 10 อักขระ",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString ต้องเป็นหนึ่งใน [red green]",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt ต้องเป็นหนึ่งใน [5 63]",
		},
		{
			ns:       "Test.JSONString",
			expected: "JSONString ต้องเป็นสตริง JSON",
		},
		{
			ns:       "Test.LowercaseString",
			expected: "LowercaseString ต้องเป็นอักษรตัวพิมพ์เล็ก",
		},
		{
			ns:       "Test.UppercaseString",
			expected: "UppercaseString ต้องเป็นอักษรตัวพิมพ์ใหญ่",
		},
		{
			ns:       "Test.Datetime",
			expected: "Datetime ต้องเป็น 2006-01-02",
		},
	}

	for _, tt := range tests {

		var fe validator.FieldError

		for _, e := range errs {
			if tt.ns == e.Namespace() {
				fe = e
				break
			}
		}

		NotEqual(t, fe, nil)
		//t.Log(fe.Translate(trans))
		Equal(t, tt.expected, fe.Translate(trans))
	}

}
