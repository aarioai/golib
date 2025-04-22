package auth

import (
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/configz"
)

func (s *Service) SelfTest() {
	s.testConfig()
}

func (s *Service) testConfig() {
	panicOnEmpty("UserTokenCryptMd5Key", configz.UserTokenCryptMd5Key)
	panicOnEmpty("UserTokenShuffleBase", configz.UserTokenShuffleBase)
	coding.ValidateShuffleEncryptKeys(configz.UserTokenShuffleBase)
}
