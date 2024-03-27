package utils

import (
	"github.com/json-iterator/go"
	"reflect"
)

func MapToStruct(m map[string]interface{}, data interface{}) (err error) {
	arr, err := jsoniter.Marshal(m)
	if err != nil {
		return
	}
	err2 := jsoniter.Unmarshal(arr, &data)
	if err2 != nil {
		return
	}
	return
}

func StructToMap(obj interface{}) map[string]interface{} {
	types := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < types.NumField(); i++ {
		data[types.Field(i).Name] = values.Field(i).Interface()
	}
	return data
}

// MapWithInField - 从map中提取指定字段
func MapWithInField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for _, val := range field {
		result[val] = data[val]
	}
	return
}

// MapWithNotField -  从map中排除指定字段
func MapWithNotField[T map[string]any](data T, field []string) (result T) {
	result = make(T)
	for key, val := range data {
		if !IsArrayContain[string](key, field) {
			result[key] = val
		}
	}
	return
}
