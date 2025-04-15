package coding

import (
	"github.com/aarioai/airis/pkg/afmt"
	"math/rand"
	"strings"
	"time"
)

// RandAESKey 生成随机length位AES密钥。密钥分为128位（16字节）、192位（24字节）、256位（32字节）
// key 支持 [
//func RandAESKey(length int, seed ...int64) []byte {
//	return []byte(Rand(length, shuffledASCIICodes, seed...))
//}

// Rand 是一个通用的随机字符生成器
// 一般随机数都很短，返回字符串最方便，且安全。若返回byte，则要处理rune，比较麻烦
func Rand[T byte | rune](length int, charSet []T, seed ...int64) string {
	firstSeed := afmt.First(seed)
	if firstSeed == 0 {
		firstSeed = time.Now().UnixNano()
	}
	r := rand.New(rand.NewSource(firstSeed))

	b := strings.Builder{}
	b.Grow(length)
	for i := 0; i < length; i++ {
		c := charSet[r.Intn(len(charSet))]
		b.WriteRune(rune(c))
	}
	return b.String()
}

// RandNum
// @warn 每次更新文件的时候，随机数就会变化；如果不更新文件，就是伪随机数（不停重复生成）
// @Description: crypto/rand是为了提供更好的随机性满足密码对随机数的要求，在linux上已经有一个实现就是/dev/urandom，crypto/rand 就是从这个地方读“真随机”数字返回，但性能比较慢
// @test CloneShuffledNums
func RandNum(length int, seed ...int64) string {
	return Rand(length, shuffledNums, seed...)
}

// RandNumLowers 生成n位 [a-z\d]（数字和小写字母）区间的随机字符
// @warn 每次更新文件的时候，随机数就会变化；如果不更新文件，就是伪随机数（不停重复生成）
// @test CloneShuffledNumLowers
func RandNumLowers(length int, seed ...int64) string {
	return Rand(length, shuffledLowerNumbers, seed...)
}

// RandAlphabets 生成n位 [a-zA-Z\d]（数字、大小写字符）区间的随机字符
// @warn 每次更新文件的时候，随机数就会变化；如果不更新文件，就是伪随机数（不停重复生成）
// @test CloneShuffledAlphabets
func RandAlphabets(length int, seed ...int64) string {
	return Rand(length, shuffledAlphabets, seed...)
}

func RandASCIICode(length int, seed ...int64) string {
	return Rand(length, shuffledASCIICodes, seed...)
}
