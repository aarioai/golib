package strs_test

import (
	"github.com/aarioai/golib/lib/code/strs"
	"testing"
)

const text = `
对于～ ， 。这样的中文标点符号，如何通过Go进行检测？
我尝试使用包unicode的范围表，就像下面的代码一样，但是Han没有包含那些标点符号字符。
你能告诉我在这个任务中我应该使用哪一个值域表吗？(请不要使用regex，因为它的性能很低。)
mediumint 一个中等大小整数，有符号的范围是-8388608到8388607，无符号的范围是0到16777215。 一位大小为3个字节。
How China's award-winning EUV breakthrough sidesteps US chip ban?
It's -- No. 1.
M&B is a famous brand.
`

func TestWordNumber(t *testing.T) {
	words, hanTotal, engTotal, numTotal := strs.HanSplit([]rune(text), '的')
	if len(words) != (hanTotal + engTotal + numTotal) {
		t.Log(words)
		t.Errorf("word number is incorrect, %d != %d + %d + %d", len(words), hanTotal, engTotal, numTotal)
	}
}
