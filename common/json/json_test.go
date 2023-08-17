package json

import (
	uuid "github.com/satori/go.uuid"
	"os"
	"testing"
)

func TestMarshal(t *testing.T) {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := Marshal(group)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	os.Stdout.Write(b)
}

func TestMarshalToString(t *testing.T) {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := MarshalToString(group)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	os.Stdout.WriteString(b)
}

func TestUnmarshal(t *testing.T) {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := Unmarshal(jsonBlob, &animals)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	t.Logf("%+v", animals)
}

func TestUnmarshalFromString(t *testing.T) {
	var jsonBlob = `[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := UnmarshalFromString(jsonBlob, &animals)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	t.Logf("%+v", animals)
}

func TestGet(t *testing.T) {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	str := Get(val, "Colors", 0).ToString()

	t.Logf("Get json value:%v", str)

}

func TestValid(t *testing.T) {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)

	v := Valid(val)

	if !v {
		t.Errorf("valid failed!")
	}

	t.Logf("valid result:%t", v)
}

func TestMarshalIndent(t *testing.T) {
	u := uuid.NewV4()

	s := struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Age     uint   `json:"age"`
		Sex     bool   `json:"sex"`
		Address string `json:"address"`
	}{u.String(), "test", 26, true, "anywhere"}

	resp2, _ := MarshalIndent(s, "", "    ")

	t.Logf("resp2:%s\n", string(resp2))
}
