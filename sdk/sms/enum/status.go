package enum

import "github.com/aarioai/airis/aa/aenum"

type SendStatus int8

const (
	SendFailed  = SendStatus(aenum.Failed)
	SendUnknown = 0
	SendOK      = SendStatus(aenum.Passed)
)
