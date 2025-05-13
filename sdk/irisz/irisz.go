package irisz

import (
	"github.com/aarioai/airis/aa/ae"
)

const prefix = "libsdk_irisz: "

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE(prefix + err.Error())
}

func PanicOnErrors(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(prefix + err.Error())
		}
	}
}
