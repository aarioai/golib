package cachez

import "github.com/aarioai/airis/aa/ae"

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE("libsdk_cachez: "+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE("libsdk_cachez: " + err.Error())
}
