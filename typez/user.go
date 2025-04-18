package typez

type AdminLevel uint8

func (t AdminLevel) Valid() bool {
	return t > 0
}
