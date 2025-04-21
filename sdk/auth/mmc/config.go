package mmc

import (
	"github.com/aarioai/golib/lib/code/coding"
	"time"
)

// base64.RawURLEncoding 模式下的字符
const (
	fpHeaderLength           = 10
	fpBodySegmentLength      = 6
	base64EncodeURL          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	fingerprintSeparator     = "(^_^)"
	fingerprintValidDuration = 30 * time.Minute
)

var AtomicFingerprintIdSeq uint64

func NewFingerprintId() uint64 {
	return coding.Uint64Id(&AtomicFingerprintIdSeq)
}

func (s *Service) MmcRSAPubkeyDERBase64() (string, error) {
	return s.app.Config.MustGetString(s.pubDERBase64KeyName)
}
func (s *Service) mmcRSAPrivkeyDER() (string, error) {
	return s.app.Config.MustGetString(s.privDERKeyName)
}
func (s *Service) mmcGCMKey() (string, error) {
	return s.app.Config.MustGetString(s.gcmKeyName)
}
