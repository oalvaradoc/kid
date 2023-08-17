package masker

import (
	"fmt"
	"testing"
)


func TestOverlayMasker2(t *testing.T) {
	o := NewOverlayMasker()
	fmt.Println(len("สุขสง่าเจริญ"))
	res, _ := o.Do("", "สุขสง่าเจริญ", "0.2")
	t.Logf("1.The result of string mask:[%s]", res)

	res, _ = o.Do("", "อภิสิทธิ์", "0.2")
	t.Logf("1.The result of string mask:[%s]", res)

	res, _ = o.Do("", "อ", "0.2")
	t.Logf("1.The result of string mask:[%s]", res)

	res, _ = o.Do("", "อภิสิทธิ์", "0,2")
	t.Logf("1.The result of string mask:[%s]", res)
}

func TestOverlayMasker(t *testing.T) {
	o := NewOverlayMasker()

	res, _ := o.Do("", "1234567890", "0.3")
	t.Logf("1.The result of string mask:[%s]", res)

	res, _ = o.Do("", "1234567890", "0.4")
	t.Logf("2.The result of string mask:[%s]", res)

	res, _ = o.Do("", "1234567890", "2,7")
	t.Logf("3.The result of string mask:[%s]", res)

	res, _ = o.Do("", "1234567890", "0")
	t.Logf("4.The result of string mask:[%s]", res)

	res, _ = o.Do("", "1234567890", "2,3")
	t.Logf("5.The result of string mask:[%s]", res)

	passwordMasker := NewPasswordMasker()

	res, _ = passwordMasker.Do("1234567890", "")
	t.Logf("6.The result of string mask:[%s]", res)

	base64Masker := NewBase64Masker()
	res, _ = base64Masker.Do("1234567890", "")
	t.Logf("7.The result of string mask:[%s]", res)
}

func TestGetMasker(t *testing.T) {
	masker, _ := GetMasker(Base64)
	res, _ := masker.Do("", "1234567890", "")
	t.Logf("The result of mask:[%s]", res)

	masker, _ = GetMasker(Overlay)
	res, _ = masker.Do("", "1234567890", "0.3")
	t.Logf("The result of mask:[%s]", res)

	masker, _ = GetMasker(Password)
	res, _ = masker.Do("", "1234567890", "")
	t.Logf("The result of mask:[%s]", res)
}

type MyMasker struct{}

func (b *MyMasker) Do(keyPath, in string, parameters ...string) (string, error) {
	return "*****[" + in + "]*****", nil
}

func TestJsonBodyMask(t *testing.T) {
	RegisterMasker("MyMasker", &MyMasker{})
	sourceJson := `{
	"data":{
		"password": "123456",
		"card": "0123456789",
		"address": "this is a test address",
		"remark": "this is a test remark",
		"other1": "this is other field 1",
		"other2": "this is other field 2"
	}
}`
	finalJson, _, gerr := JsonBodyMask([]byte(sourceJson), []string{
		"data.password|password",
		"data.card|overlay|0.6",
		"data.address|overlay|4,8",
		"data.remark|base64",
		"data.other1|myMasker",
	})

	if nil != gerr {
		t.Errorf("the error result is:%++v", gerr)
	}

	t.Logf("The final json string:%s", string(finalJson))
}
