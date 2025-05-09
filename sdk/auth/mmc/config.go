package mmc

import (
	"github.com/aarioai/golib/lib/code/coding"
	"time"
)

// base64.RawURLEncoding 模式下的字符
const (
	fpHeaderLength           = 10
	fpBodySegmentLength      = 6
	fingerprintSeparator     = "(^_^)"
	fingerprintValidDuration = 30 * time.Minute
)

var AtomicFingerprintIdSeq uint64

func NewFingerprintId() uint64 {
	return coding.Uint64Id(&AtomicFingerprintIdSeq)
}

func (s *Service) rsaPubDERBase64() (string, error) {
	return s.app.Config.MustGetString(s.pubDERBase64KeyName)
}
func (s *Service) rsaPrivDER() (string, error) {
	return s.app.Config.MustGetString(s.privDERKeyName)
}
func (s *Service) gcmKey() (string, error) {
	return s.app.Config.MustGetString(s.gcmKeyName)
}
