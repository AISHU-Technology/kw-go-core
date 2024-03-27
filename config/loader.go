package config

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/common"
	"github.com/AISHU-Technology/kw-go-core/errorx"
	"github.com/AISHU-Technology/kw-go-core/filex"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"github.com/BurntSushi/toml"
	"github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

func ReadLoad(path string, t any) error {
	suffix := strings.ToLower(filex.GetFileSuffix(path))
	ok := utils.IsConfSuffix(suffix)
	if !ok {
		return fmt.Errorf("%s=ReadLoad=unrecognized file type: %s", common.ErrorTypeSupported, path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return loader(b, t, suffix)
}

func loader(b []byte, t any, suffix string) error {
	switch suffix {
	case common.JsonType:
		return jsoniter.Unmarshal(b, t)
	case common.TomlType:
		return toml.Unmarshal(b, t)
	case common.YamlType:
		return yaml.Unmarshal(b, t)
	case common.YmlType:
		return yaml.Unmarshal(b, t)
	}
	return errorx.NewCodeErrorMsg(common.ErrorTypeSupported)
}

func ReadLoadMap(path string) (map[string]interface{}, error) {
	suffix := strings.ToLower(filex.GetFileSuffix(path))
	ok := utils.IsConfSuffix(suffix)
	if !ok {
		return nil, fmt.Errorf("%s=ReadLoadMap=unrecognized file path=: %s", common.ErrorTypeSupported, path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loaderMap(b, suffix)
}

func loaderMap(b []byte, suffix string) (map[string]interface{}, error) {
	switch suffix {
	case common.JsonType:
		return utils.JsonToMap(b)
	case common.TomlType:
		return utils.TomlToMap(b)
	case common.YamlType:
		return utils.YamlToMap(b)
	case common.YmlType:
		return utils.YamlToMap(b)
	}
	return nil, errorx.NewCodeErrorMsg(common.ErrorTypeSupported)
}

func ReadLoadMaps(path string, maps map[string]interface{}) (map[string]interface{}, error) {
	suffix := strings.ToLower(filex.GetFileSuffix(path))
	ok := utils.IsConfSuffix(suffix)
	if !ok {
		str := fmt.Sprintf("=ReadLoadMaps=unrecognized file suffix=: %s", suffix)
		log.Fatalf(str)
		return nil, errorx.NewCodeErrorMsg(str)
	}
	b, err := os.ReadFile(path)
	if err != nil || len(b) <= 0 {
		str := fmt.Sprintf("error: file is empty path=%s, %s", path, err)
		log.Fatalf(str)
		return nil, errorx.NewCodeErrorMsg(str)
	}
	return loaderMaps(b, suffix, maps)
}

func loaderMaps(b []byte, suffix string, maps map[string]interface{}) (map[string]interface{}, error) {
	switch suffix {
	case common.JsonType:
		return utils.JsonToMaps(b, maps)
	case common.TomlType:
		return utils.TomlToMaps(b, maps)
	case common.YamlType:
		return utils.YamlToMaps(b, maps)
	case common.YmlType:
		return utils.YamlToMaps(b, maps)
	}
	return nil, errorx.NewCodeErrorMsg(common.ErrorTypeSupported)
}
