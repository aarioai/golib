package coding

import (
	"bytes"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"math/rand"
)

// 禁止全局共享slice！！！不可避免其他人参与编程，新手失误导致修改slice。
// 即使使用局部共享slice，对slice的操作，必须要增加测试用例，判断使用前后是否发生变化。
var (
	// @test CloneNumberRunes
	numberRunes = []byte(stdfmt.Numbers)
	// @test CloneLowercaseRunes
	lowercaseRunes = []byte(stdfmt.Lowercases)
	// @test CloneUppercaseRunes
	uppercaseRunes = []byte(stdfmt.Uppercases)
	// @test CloneNumLowerRunes
	lowerNumberRunes = []byte(stdfmt.Base36Digits)
	// @test CloneAlphabetRunes
	alphabetRunes = []byte(stdfmt.Base62Digits)
	// @test CloneAsciiCodeRunes
	asciiCodesRunes = []byte(stdfmt.ReadableAsciiCodes)
)

// 不可逆的随机数生成
var (
	// @test CloneShuffledNums
	shuffledNums []byte
	// @test CloneShuffledNumLowers
	shuffledLowerNumbers []byte
	// @test CloneShuffledAlphabets
	shuffledAlphabets []byte
	// @test CloneShuffledAlphabets
	shuffledASCIICodes []byte
)

func init() {
	shuffledNums = Shuffle(numberRunes)
	shuffledLowerNumbers = Shuffle(lowerNumberRunes)
	shuffledAlphabets = Shuffle(alphabetRunes)
	shuffledASCIICodes = Shuffle(asciiCodesRunes)
}

// Shuffle 使用洗牌算法将切片随机打乱
func Shuffle[T byte | rune](s []T) []T {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

// 随机从  randStart-randEnd 之间找一个字符填充
func RandPad(str string, minlen int, randStart, randEnd byte, rightPad bool) string {
	n := int(randEnd) - int(randStart) + 1
	for len(str) < minlen {
		pad := byte(rand.Intn(n)) + randStart
		if rightPad {
			str += string(pad)
		} else {
			str = string(pad) + str
		}
	}
	return str
}

func CloneNumberRunes() []byte {
	return bytes.Clone(numberRunes)
}
func CloneLowercaseRunes() []byte {
	return bytes.Clone(lowercaseRunes)
}
func CloneUppercaseRunes() []byte {
	return bytes.Clone(uppercaseRunes)
}
func CloneNumLowerRunes() []byte {
	return bytes.Clone(lowerNumberRunes)
}
func CloneAlphabetRunes() []byte {
	return bytes.Clone(alphabetRunes)
}
func CloneAsciiCodeRunes() []byte {
	return bytes.Clone(asciiCodesRunes)
}
func CloneShuffledNums() []byte {
	return bytes.Clone(shuffledNums)
}
func CloneShuffledNumLowers() []byte {
	return bytes.Clone(shuffledLowerNumbers)
}
func CloneShuffledAlphabets() []byte {
	return bytes.Clone(shuffledAlphabets)
}
