package masker

import (
	"strings"
	"testing"
)

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
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	overlayed += overlay
	overlayed += string(r[end:])
	return overlayed
}

func Mobile(i string) string {
	if len(i) == 0 {
		return ""
	}
	return overlay(i, "****", 3, 7)
}

func Password(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return strings.Repeat("*", len(i))
}

const (
	MPassword MaskType = "password"
	MMobile            = "mobile"
)

func TestStruct(t *testing.T) {
	RegisterMaskFunc(MPassword, Password)
	RegisterMaskFunc(MMobile, Mobile)

	type InnerStruct struct {
		Name     string
		Password string  `json:"Password" validate:"required" mask:"password"`
		Amount   float64 `json:"amount"`
	}
	type User struct {
		Name     string
		Mobile   string       `mask:"mobile"`
		Password string       `json:"Password" validate:"required" mask:"password"`
		Inner1   *InnerStruct `mask:"struct"`
		Inner2   InnerStruct  `mask:"struct"`
	}

	user := User{
		Name:     "User 1",
		Mobile:   "15019491171",
		Password: "test password1",
		//Inner1: &InnerStruct{
		//	Name: "Inner struct 1",
		//	Password: "test password 2",
		//},
		Inner2: InnerStruct{
			Name:     "Inner struct 2",
			Password: "test password 3",
			Amount:   100,
		},
	}

	got := Do(user)
	t.Logf("got:[%++v]", got)

}
