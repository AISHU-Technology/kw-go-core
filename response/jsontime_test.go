package response

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/date"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

type DateTime struct {
	StartTime JsonTime `json:"StartTime"`
	EndTime   JsonTime `json:"StartTime"`
}

type Pig struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type Dog struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	CreatedTime JsonTime `json:"created_time"`
	UpdatedTime JsonTime `json:"updated_time"`
}

func (t Pig) MarshalJSON() ([]byte, error) {
	type TmpJSON Pig
	return jsoniter.Marshal(&struct {
		TmpJSON
		CreatedTime JsonTime `json:"created_time"`
		UpdatedTime JsonTime `json:"updated_time"`
	}{
		TmpJSON:     (TmpJSON)(t),
		CreatedTime: JsonTime(t.CreatedTime),
		UpdatedTime: JsonTime(t.UpdatedTime),
	})
}

func TestJsonPigTimesTructFunc(t *testing.T) {
	good := Pig{123, "pigpig", time.Now(), time.Now()}
	bytes, _ := jsoniter.Marshal(good)
	str := utils.ToString(bytes)
	fmt.Printf("%s", str)
}

func TestJsonDogTimeFunc(t *testing.T) {
	good := Dog{
		ID: 123, Name: "dogdog",
		CreatedTime: JsonTime(time.Now()),
		UpdatedTime: JsonTime(time.Date(2019, 1, 1, 20, 8, 23, 28, time.Local)),
	}
	bytes, err := jsoniter.Marshal(good)
	if err != nil {
		println(err)
	}
	str := utils.ToString(bytes)
	fmt.Printf("%s", str)
}

func TestJSONTimeFunc(t *testing.T) {
	p := &DateTime{}
	//p := new(DateTime)
	str := `{"start":"2023-09-10 18:12:49","end":"2019-12-10 18:12:49"}`
	err := jsoniter.Unmarshal([]byte(str), p)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("=====", p, str)
	fmt.Println("====sss=", p.EndTime, p.StartTime)
	stra := date.TimeToStr(time.Time(p.EndTime), date.YYYY_MM_DD_HH_MM_SS)
	fmt.Println("=====ssssssss", stra)
	bytes, _ := jsoniter.Marshal(p)
	strs := utils.ToString(bytes)
	fmt.Printf("%s", strs)
}

func TestStrToTimeFunc(t *testing.T) {
	times, _ := date.StrToTime("2023-08-10 18:12:49", date.YYYY_MM_DD_HH_MM_SS)
	fmt.Println("=====", times)
	str := date.TimeToStr(times, date.YYYY_MM_DD_HH_MM_SS)
	fmt.Println("=====ssssssss", str)
}
