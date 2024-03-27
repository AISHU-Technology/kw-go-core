package response

import (
	"github.com/AISHU-Technology/kw-go-core/common"
	"github.com/AISHU-Technology/kw-go-core/i18n"
	"github.com/json-iterator/go"
	"net/http"
)

func getMsg(i, t string) string {
	i18nFormat := map[string]map[string]string{
		"zh-CN": {
			"ok":   common.SuccessfulCn,
			"fail": common.FailCn,
		},
		"en_US": {
			"ok":   common.SuccessfulEn,
			"fail": common.FailEn,
		},
	}
	return i18nFormat[i][t]
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ToString 返回 JSON 格式的错误详情，会忽略time类型格式化
func ResultToStr(code int, msg string, data interface{}) *string {
	err := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	raw, _ := jsoniter.Marshal(err)
	r := string(raw)
	return &r
}

func Result(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
func ResultData(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func Ok(r *http.Request) *Response {
	return &Response{
		Code: common.SuccessCode,
		Msg:  getMsg(i18n.GetAcceptLanguage(r), "ok"),
		Data: nil,
	}
}

func Fail(r *http.Request) *Response {
	return &Response{
		Code: common.FailCode,
		Msg:  getMsg(i18n.GetAcceptLanguage(r), "fail"),
		Data: nil,
	}
}

// OkMsg
/**
 * @Description: 返回带消息的成功
 * @param c
 * @param message
 */
func OkMsg(msg string) *Response {
	return &Response{
		Code: common.SuccessCode,
		Msg:  msg,
		Data: nil,
	}
}

// OkDataMsg
/**
 * @Description: 返回带消息和数据的成功
 * @param c
 * @param data
 * @param message
 */
func OkDataMsg(data interface{}, msg string) *Response {
	return &Response{
		Code: common.SuccessCode,
		Msg:  msg,
		Data: data,
	}
}

// FailMsg
/**
 * @Description: 返回带消息的失败
 * @param c
 * @param err
 */
func FailMsg(msg string) *Response {
	return &Response{
		Code: common.FailCode,
		Msg:  msg,
	}
}

// FailDataMsg
/**
 * @Description: 返回带消息和数据的失败
 * @param c
 * @param data
 * @param err
 */
func FailDataMsg(data interface{}, msg string) *Response {
	return &Response{
		Code: common.FailCode,
		Msg:  msg,
		Data: data,
	}
}
