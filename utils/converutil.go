package utils

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/common"
	"github.com/AISHU-Technology/kw-go-core/errorx"
	"github.com/BurntSushi/toml"
	"github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func ToString(i interface{}) (str string) {
	switch i.(type) {
	case string:
		str = i.(string)
	case int:
		str = strconv.Itoa(i.(int))
	case int8:
		str = fmt.Sprint(i.(int8))
	case int16:
		str = fmt.Sprint(i.(int16))
	case int32:
		str = string(i.(int32))
	case int64:
		str = strconv.FormatInt(i.(int64), 10)
	case uint:
		str = strconv.Itoa(i.(int))
	case uint8:
		str = string(i.(uint8))
	case uint16:
		str = fmt.Sprint(i.(uint16))
	case uint32:
		str = fmt.Sprint(i.(uint32))
	case uint64:
		str = strconv.FormatUint(i.(uint64), 10)
	case float32:
		str = fmt.Sprintf("%f", i.(float32))
	case float64:
		str = strconv.FormatFloat(i.(float64), 'f', -1, 32)
	case time.Time:
		str = i.(time.Time).Format("2006-01-02 15:04:05")
	case []byte:
		b := i.([]byte)
		str = *(*string)(unsafe.Pointer(&b))
	case error:
		str = i.(error).Error()
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func StrToInt(i interface{}) (num int, err error) {
	switch i.(type) {
	case string:
		num, err = strconv.Atoi(i.(string))
	}
	return
}

func ToInt32(i interface{}) (num int32, err error) {
	switch i.(type) {
	case int:
		num = int32(i.(int))
	case int32:
		num = i.(int32)
	case int64:
		// 有可能造成精度丢失
		num = int32(i.(int64))
	case float32:
		// 有可能造成精度丢失
		num = int32(i.(float32))
	case float64:
		// 有可能造成精度丢失
		num = int32(i.(float64))
	case string:
		n, e := strconv.Atoi(i.(string))
		num = int32(n)
		err = e
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func ToInt64(i interface{}) (num int64, err error) {
	switch i.(type) {
	case int:
		num = int64(i.(int))
	case int32:
		num = int64(i.(int32))
	case int64:
		num = i.(int64)
	case uint64:
		num = int64(i.(uint64))
	case float32:
		num = int64(i.(float32))
	case float64:
		num = int64(i.(float64))
	case string:
		num, err = strconv.ParseInt(i.(string), 10, 64)
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func ToFloat32(i interface{}) (num float32, err error) {
	switch i.(type) {
	case string:
		// string无法直接转换float32，只能先转换为float64，再通过float64转float32
		var num64 float64
		num64, err = strconv.ParseFloat(i.(string), 32)
		num = float32(num64)
	case int:
		num = float32(i.(int))
	case int32:
		num = float32(i.(int32))
	case int64:
		num = float32(i.(int64))
	case float32:
		num = i.(float32)
	case float64:
		// 可能造成精度丢失
		num = float32(i.(float64))
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func ToFloat64(i interface{}) (num float64, err error) {
	switch i.(type) {
	case string:
		num, err = strconv.ParseFloat(i.(string), 64)
	case int:
		num = float64(i.(int))
	case int32:
		num = float64(i.(int32))
	case int64:
		num = float64(i.(int64))
	case float32:
		num = float64(i.(float32))
	case float64:
		num = i.(float64)
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func ToByteArray(i interface{}) (b []byte) {
	switch i.(type) {
	case string:
		str := i.(string)
		return *(*[]byte)(unsafe.Pointer(&str))
	default:
		panic(common.ErrorTypeSupported)
	}
	return
}

func MapToQueryStr(param map[string]any) string {
	if param == nil {
		return ""
	}
	var keys []string
	for key := range param {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var build strings.Builder
	for i, v := range keys {
		build.WriteString(v)
		build.WriteString("=")
		build.WriteString(fmt.Sprintf("%v", param[v]))
		if i != len(keys)-1 {
			build.WriteString("&")
		}
	}
	return build.String()
}

// Play: https://go.dev/play/p/pFqMkM40w9z
func StructToUrlValues(targetStruct any) (url.Values, error) {
	result := url.Values{}
	s := StructToMap(targetStruct)
	for k, v := range s {
		result.Add(k, fmt.Sprintf("%v", v))
	}
	return result, nil
}

func JsonToStruct(value string, t any) error {
	if IsBlank(value) {
		return errorx.NewCodeErrorMsg("json velue is not empty~")
	}
	jsoniter.Unmarshal([]byte(value), t)
	return nil
}

func StructToJson(t any) (string, error) {
	if t == nil {
		return "", errorx.NewCodeErrorMsg("Struct velue is not empty~")
	}
	marshal, _ := jsoniter.Marshal(t)
	return string(marshal), nil
}

func JsonToMap(b []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := jsoniter.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func JsonToMaps(b []byte, result map[string]interface{}) (map[string]interface{}, error) {
	err := jsoniter.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func TomlToMap(b []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := toml.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func TomlToMaps(b []byte, result map[string]interface{}) (map[string]interface{}, error) {
	err := toml.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func YamlToMap(b []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := yaml.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func YamlToMaps(b []byte, result map[string]interface{}) (map[string]interface{}, error) {
	err := yaml.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
