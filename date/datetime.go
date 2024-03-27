package date

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/errorx"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"strings"
	"time"
)

const (
	YYYY_MM_DD_HH_MM_SS string = "yyyy-mm-dd hh:mm:ss"
	HH_MM_SS            string = "hh:mm:ss"
	YYYYY_MM_DD         string = "yyyy-mm-dd"
)

func GetTimeFormat(format string) string {
	timeFormat := map[string]string{
		YYYY_MM_DD_HH_MM_SS:   "2006-01-02 15:04:05",
		"yyyy-mm-dd hh:mm":    "2006-01-02 15:04",
		"yyyy-mm-dd hh":       "2006-01-02 15",
		YYYYY_MM_DD:           "2006-01-02",
		"yyyy-mm":             "2006-01",
		"mm-dd":               "01-02",
		"dd-mm-yy hh:mm:ss":   "02-01-06 15:04:05",
		"yyyy/mm/dd hh:mm:ss": "2006/01/02 15:04:05",
		"yyyy/mm/dd hh:mm":    "2006/01/02 15:04",
		"yyyy/mm/dd hh":       "2006/01/02 15",
		"yyyy/mm/dd":          "2006/01/02",
		"yyyy/mm":             "2006/01",
		"mm/dd":               "01/02",
		"dd/mm/yy hh:mm:ss":   "02/01/06 15:04:05",
		"yyyymmdd":            "20060102",
		"mmddyy":              "010206",
		"yyyy":                "2006",
		"yy":                  "06",
		"mm":                  "01",
		HH_MM_SS:              "15:04:05",
		"hh:mm":               "15:04",
		"mm:ss":               "04:05",
	}
	return timeFormat[strings.ToLower(format)]
}

func GetNowDateTime() string {
	return time.Now().Format(GetTimeFormat(YYYY_MM_DD_HH_MM_SS))
}

func GetNowTime() string {
	return time.Now().Format(GetTimeFormat(HH_MM_SS))
}

func GetNowDate() string {
	return time.Now().Format(GetTimeFormat(YYYYY_MM_DD))
}

func GetAddMinute(t time.Time, minute int64) time.Time {
	return t.Add(time.Minute * time.Duration(minute))
}

func GetAddHour(t time.Time, hour int64) time.Time {
	return t.Add(time.Hour * time.Duration(hour))
}

func GetAddDay(t time.Time, day int64) time.Time {
	return t.Add(24 * time.Hour * time.Duration(day))
}

func GetNowAddMinute(minute int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(minute))
}

func GetNowAddHour(hour int64) time.Time {
	return time.Now().Add(time.Hour * time.Duration(hour))
}

func GetNowAddDay(day int64) time.Time {
	return time.Now().Add(24 * time.Hour * time.Duration(day))
}

func TimeToStr(t time.Time, format string, timezone ...string) string {
	fr := GetTimeFormat(format)
	if utils.IsBlank(fr) {
		return ""
	}

	if timezone != nil && timezone[0] != "" {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return ""
		}
		return t.In(loc).Format(fr)
	}
	return t.Format(fr)
}

func StrToTime(str, format string, timezone ...string) (time.Time, error) {
	fr := GetTimeFormat(format)
	if utils.IsBlank(fr) {
		return time.Time{}, errorx.NewCodeErrorMsg(fmt.Sprintf("format %s not support", format))
	}
	if timezone != nil && timezone[0] != "" {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
		return time.ParseInLocation(fr, str, loc)
	}
	return time.Parse(fr, str)
}

func Equal(t1 time.Time, t2 time.Time) (bool, error) {
	if utils.IsBlankTime(&t1) || utils.IsBlankTime(&t2) {
		return false, fmt.Errorf("==t1 or t2 is not empty=")
	}
	return t1.Equal(t2), nil
}
