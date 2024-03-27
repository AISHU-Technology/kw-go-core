package utils

import (
	"fmt"
	"testing"
)

type Dog struct {
	Name string `json:"name"`
	Vae  string `json:"vae"`
}

func TestStructToMap(t *testing.T) {
	stu := Dog{Name: "<下购>", Vae: "<女>"}
	values := StructToMap(stu)
	for k, v := range values {
		fmt.Println(fmt.Sprintf("===key===%v==value==%v", k, v))
	}

	var d Dog
	err := MapToStruct(values, &d)
	fmt.Println(d.Name, err)
}
