package gormx

import (
	"github.com/AISHU-Technology/kw-go-core/filter"
	"github.com/AISHU-Technology/kw-go-core/page"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"net/url"
	"strings"
)

func BuildQueryParams(queryParams url.Values) (query url.Values, pageModel page.PageModel) {
	params := url.Values{}
	var arr []string
	var page page.PageModel
	for key, values := range queryParams {
		if utils.IsBlankArr(values) {
			continue
		}
		switch key {
		case "sort":
			params["sort"] = values
		case "select":
			params["select"] = values
		case "omit":
			params["omit"] = values
		case "gcond":
			params["gcond"] = values
		case "pageNo":
			pageNo, _ := utils.StrToInt(values[0])
			page.PageNo = pageNo
		case "pageSize":
			pageSize, _ := utils.StrToInt(values[0])
			page.PageSize = pageSize
		default:
			for _, value := range values {
				columKey := parseColumn(filter.DealNoStringType(key))
				arr = append(arr, columKey+"="+filter.DealNoStringType(filter.SqlFilter(value, false)))
			}
		}
	}
	if len(arr) > 0 {
		params["q"] = arr
	}
	return params, page
}

func parseColumn(key string) string {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return ""
	}
	var str strings.Builder
	strArr := strings.Split(key, "")
	for index, value := range strArr {
		if utils.IsInUpChar(value) {
			if index == 0 {
				value = strings.ToLower(value)
			} else {
				value = "_" + strings.ToLower(value)
			}
		}
		str.WriteString(value)
	}
	return str.String()
}
