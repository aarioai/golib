package crypto

import (
	"github.com/gobwas/glob/util/runes"
)

// 基于简易DFA算法，匹配首位
var sensitivePairs = map[rune][][]rune{
	'八': {[]rune("学潮")},
	'裆': {[]rune("中央")},
	'党': {[]rune("中央")},
	'法': {{'轮'}, {'功'}},
	'胡': {{'涛'}},
	'共': {{'党'}, {'裆'}, {'匪'}},
	'江': {{'青', '民', '蛤', '蟆'}},
	'毛': {{'东'}, []rune("腊肉")},
	'猫': {{'东'}, []rune("腊肉")},
	'四': {[]rune("人帮")},
	'天': {[]rune("安门事件")},
	'文': {{'革'}},
	'温': {{'宝'}},
	'习': {{'平'}, []rune("包子"), []rune("明泽")},
}

const null = rune(0) // byte 0

func IsSensitive(s string) bool {
	t := []rune(s)
	x := len(t)
	for i := 0; i < x; i++ {
		c := t[i]
		for start, ends := range sensitivePairs {
			// 首字匹配
			if c == start && (x-1) != i {
				n := i + 8 // 8 个延迟
				if n > x {
					n = x
				}
				sli := t[i+1 : n]
				for _, end := range ends {

					m := runes.Index(sli, end)
					if m < 0 {
						continue
					}
					return true
				}
			}
		}
	}
	return false
}

// 这里只做简单DFA过滤
// 不能用string, 必须要用 []rune，否则中文切片 [:]会出不对
func FilterSensitiveWords(content []rune) []rune {
	x := len(content)
	for i := 0; i < x; i++ {
		c := content[i]
		for start, ends := range sensitivePairs {
			// 首字匹配
			if c == start && (x-1) != i {
				n := i + 8 // 8 个延迟
				if n > x {
					n = x
				}
				sli := content[i+1 : n]
				for _, end := range ends {
					m := runes.Index(sli, end)
					if m < 0 {
						continue
					}
					y := i + m + len(end)
					// 敏感词匹配
					for k := i; k < y+1; k++ {
						content[k] = null
					}
					i = y // -1 去掉 for循环结束后执行的 i++
					break
				}
			}
		}
	}
outer:
	for {
		j := -1
		x := len(content) - 1
		for i, c := range content {
			if c == null { // ascii 0
				if j == -1 {
					j = i
				}
				if i == x {
					content = content[:j]
					break outer
				}
			} else {
				if j > -1 {
					content[i-1] = '*'
					content = append(content[:j], content[i-1:]...)
					continue outer
				}
				if i == x {
					break outer
				}
			}
		}
	}
	return content
}
