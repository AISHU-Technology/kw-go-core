package utils

import (
	"fmt"
	"testing"
)

type Student struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade int    `json:"grade"`
}

func TestGetFieldName(t *testing.T) {
	fmt.Println(GetFieldName(Student{Name: "1", Age: 2, Grade: 3}))
	fmt.Println(GetFieldName(&Student{Name: "4", Age: 5, Grade: 6}))
	fmt.Println(GetFieldName(""))
	fmt.Println(GetTagName(&Student{Name: "7", Age: 8, Grade: 9}))
}
