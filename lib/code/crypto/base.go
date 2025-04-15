package crypto

import (
	"errors"
)

const (
	base10 = 10
	base36 = 36
	base64 = 64
)

var (
	ErrEncryptFailed = errors.New("encrypt failed")
	ErrDecryptFailed = errors.New("decrypt failed")
)
 