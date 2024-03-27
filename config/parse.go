package config

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/utils"
)

func SqlParseMap(maps map[string]interface{}, key string) (string, error) {
	if utils.IsBlankMap(maps) {
		return "", fmt.Errorf("=SqlParse=map is empty==key: %s", key)
	}
	v := maps[key]
	if v == nil {
		return "", fmt.Errorf("=SqlParse=get value is empty plase check key!==key: %s", key)
	}
	strValue := utils.ToString(v)
	if utils.IsBlank(strValue) {
		return "", fmt.Errorf("=SqlParse==value is empty==key: %s", key)
	}
	return strValue, nil
}
