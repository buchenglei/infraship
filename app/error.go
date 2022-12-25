package app

import (
	"errors"
	"fmt"
)

type ErrCode uint32

func newErrCode(category, code uint16) ErrCode {
	eCode := uint32(category) << 16
	eCode |= uint32(code)
	return ErrCode(eCode)
}

func (e ErrCode) Unwrap() (category, code uint16) {
	category = uint16(e >> 16)
	code = uint16((e << 16) >> 16)
	return
}

type Error struct {
	errCode ErrCode
	msg     string
	err     error
}

func NewRootError(category, code uint16, msg string) Error {
	return Error{errCode: newErrCode(category, code), msg: msg}
}

func (e Error) WithErrorDesc(desc string) Error {
	return e.WithError(errors.New(desc))
}

func (e Error) WithError(err error) Error {
	newE := e
	if newE.err != nil {
		newE.err = fmt.Errorf("%w:\n-> %s", e.err, err.Error())
	} else {
		newE.err = err
	}

	return newE
}

func (e Error) Code() ErrCode {
	return e.errCode
}

func (e Error) GetError() error {
	cate, code := e.errCode.Unwrap()
	if e.err == nil {
		return fmt.Errorf("[%d(%d,%d)-%s]", e.errCode, cate, code, e.msg)
	}
	return fmt.Errorf("[%d(%d,%d)-%s] -> %w", e.errCode, cate, code, e.msg, e.err)
}

func (e Error) Error() string {
	return e.GetError().Error()
}
