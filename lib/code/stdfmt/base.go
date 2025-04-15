package stdfmt

import (
	"bytes"
	"github.com/aarioai/airis/pkg/types"
)

const (
	Numbers             = "0123456789"
	Lowercases          = "abcdefghijklmnopqrstuvwxyz"
	Uppercases          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SymbolsWithoutBlank = "_~+/-=<{[()]}>#:@!|$%^&*\\;'\",.?`" // 不含空格
	Base36Digits        = Numbers + Lowercases
	Base62Digits        = Base36Digits + Uppercases
	Base64Digits        = types.Base64Digits

	ReadableAsciiCodes = Base62Digits + SymbolsWithoutBlank + " "
)

// 禁止全局共享slice！！！不可避免其他人参与编程，新手失误导致修改slice。
// 即使使用局部共享slice，对slice的操作，必须要增加测试用例，判断使用前后是否发生变化。
var (
	// @test CloneNumberRunes
	numbersSlice = []byte(Numbers)
	// @test CloneLowercaseRunes
	lowercasesSlice = []byte(Lowercases)
	// @test CloneUppercaseRunes
	uppercasesSlice = []byte(Uppercases)
	// @test CloneBase36DigitsSlice
	base36DigitsSlice = []byte(Base36Digits)
	// @test CloneBase62DigitsSlice
	base62DigitsSlice = []byte(Base62Digits)
	// @test CloneAsciiCodeRunesRunes
	readableAsciiCodesSlice = []byte(ReadableAsciiCodes)
)

func CloneNumbersSlice() []byte {
	return bytes.Clone(numbersSlice)
}
func CloneLowercasesSlice() []byte {
	return bytes.Clone(lowercasesSlice)
}
func CloneUppercasesSlice() []byte {
	return bytes.Clone(uppercasesSlice)
}
func CloneBase36DigitsSlice() []byte {
	return bytes.Clone(base36DigitsSlice)
}
func CloneBase62DigitsSlice() []byte {
	return bytes.Clone(base62DigitsSlice)
}
func CloneReadableAsciiCodeRunesSlice() []byte {
	return bytes.Clone(readableAsciiCodesSlice)
}
