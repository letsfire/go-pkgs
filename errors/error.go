package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	IsOk    bool   // 是否成功
	ErrCode string // 错误编号
	ErrDesc string // 错误描述
	PrevErr *Error // 上级错误
}

func (e *Error) Use(err error) {
	if err == nil {
		*e = NoneError("OK")
	} else {
		switch err.(type) {
		case Error:
			*e = err.(Error)
		case *Error:
			*e = *(err.(*Error))
		case validator.ValidationErrors:
			*e = ClientError(err.Error())
		default:
			*e = ServerError(err.Error())
		}
	}
}

func (e Error) Panic() {
	if !e.IsOk {
		panic(e)
	}
}

func (e Error) AddPrev(err Error) Error {
	if err.IsOk {
		return e
	}
	e.PrevErr = &err
	return e
}

func (e Error) Error() string {
	if e.IsOk {
		return ""
	}
	var prev string
	if e.PrevErr != nil {
		prev = ", " + e.PrevErr.Error()
	}
	return fmt.Sprintf("%s - %s%s", e.ErrCode, e.ErrDesc, prev)
}

func NoneError(desc string) Error {
	return Error{IsOk: true, ErrCode: "none", ErrDesc: desc}
}

func NewError(code, desc string) Error {
	return Error{IsOk: false, ErrCode: code, ErrDesc: desc}
}

func NewErrorf(code, format string, args ...interface{}) Error {
	return NewError(code, fmt.Sprintf(format, args...))
}
