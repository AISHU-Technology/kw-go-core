package response

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/date"
	"time"
)

type JsonTime time.Time

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+date.GetTimeFormat(date.YYYY_MM_DD_HH_MM_SS)+`"`, string(data), time.Local)
	*t = JsonTime(now)
	return
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(date.GetTimeFormat(date.YYYY_MM_DD_HH_MM_SS)))
	var tmp = []byte(stamp)
	return tmp, nil
}
