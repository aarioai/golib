package coding

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyCipherKey      = errors.New("empty cipher key")
	ErrInvalidCipherKeyLen = errors.New("invalid cipher key length")

	ErrInvalidPKCS7Padding = errors.New("invalid PKCS7 padding")
	ErrInvalidBlockSize    = errors.New("invalid block size")
)

func ErrCipherKeyMissChar(char byte) error {
	return fmt.Errorf("cipher key miss char %c", char)
}
