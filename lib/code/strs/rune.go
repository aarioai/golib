package strs

func IsDigit[T rune | byte](r T) bool {
	return r >= '0' && r <= '9'
}
func IsAlpha[T rune | byte](r T) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
func IsAlphaDigit[T rune | byte](r T) bool {
	return IsAlpha(r) || IsDigit(r)
}
func Contains(s string, r rune) bool {
	for _, v := range s {
		if v == r {
			return true
		}
	}
	return false
}

func RunesContains(rs []rune, r rune) bool {
	for _, v := range rs {
		if v == r {
			return true
		}
	}
	return false
}
