package errorx

import "github.com/AISHU-Technology/kw-go-core/common"

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewCodeErrorMsg(msg string) error {
	return &CodeError{Code: 500, Msg: msg}
}

func NewError() error {
	return &CodeError{Code: 500, Msg: common.FailEn}
}

func (e *CodeError) Error() string {
	return e.Msg
}
