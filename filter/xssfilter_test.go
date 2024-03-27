package filter

import (
	"fmt"
	"testing"
)

type Dog struct {
	Name string `json:"name"`
	Vae  string `json:"vae"`
}

func TestXssFilter(t *testing.T) {
	stu := Dog{Name: "<下购>", Vae: "<女>"}
	err4 := XssFilter(&stu)
	fmt.Println(stu, err4)

	//map
	mdata := map[string]interface{}{"殺殺殺": "<222>", "key2": map[string]interface{}{"熱恩": "<9988<https>hyy!!>"}}
	err1 := XssFilter(mdata)
	fmt.Println(mdata, err1)

	//slice
	sdata := []string{"<fasfsf>", "<rrrr!@#rr>"}
	err2 := XssFilter(sdata)
	fmt.Println(sdata, err2)

	str := DealNoStringType("<9988<https>hyy!uuuu8890")
	fmt.Println(str)
}
