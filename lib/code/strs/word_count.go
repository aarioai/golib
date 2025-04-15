package strs

import (
	"strings"
	"unicode"
)

// 这里 r 必须用指针
func handleNumOrEngTotal(r *strings.Builder, result []string, engTotal, digTotal int) ([]string, int, int) {
	if r.Len() > 0 {
		fs := r.String()
		if strings.Trim(fs, "-") != "" {
			result = append(result, fs)

			fc := fs[0]
			if IsDigit(fc) || (fc == '-' && IsDigit(fs[1])) {
				digTotal++
			} else {
				engTotal++
			}
		}
		r.Reset()
	}
	return result, engTotal, digTotal
}

// HanSplit 统计 英文单词数、中文汉字数
func HanSplit(s []rune, ignores ...rune) (result []string, hanTotal, engTotal, digTotal int) {
	result = make([]string, 0)
	var r strings.Builder
	l := len(s)
	for i := 0; i < l; i++ {
		c := s[i]
		if len(ignores) > 0 && RunesContains(ignores, c) {
			result, engTotal, digTotal = handleNumOrEngTotal(&r, result, engTotal, digTotal)
			continue
		}
		// 不包括标点符号
		if unicode.Is(unicode.Han, c) {
			result, engTotal, digTotal = handleNumOrEngTotal(&r, result, engTotal, digTotal)
			result = append(result, string(c))
			hanTotal++
			continue
		}

		// 处理英文连词
		if r.Len() > 0 && c == '\'' && i < l-1 {
			next := s[i+1]
			if IsAlpha(next) {
				// it's, i'm, i'd, we're, i'll, d'ng,
				if i == l-1 {
					r.WriteRune(c)
					continue
				}
				ng := s[i+2]
				if ng == ' ' {
					r.WriteRune(c)
					continue
				}
				if i < l-3 && s[i+3] == ' ' {
					if (next == 'r' && ng == 'e') || (next == 'l' && ng == 'l') || (next == 'n' && ng == 'g') {
						r.WriteRune('c')
						continue
					}
				}
			}
		}
		if c == '-' || IsAlphaDigit(c) {
			r.WriteRune(c)
			continue
		}
		// R&B, 1+1=2
		if r.Len() > 0 && Contains("&+=", c) && i < l-1 {
			next := s[i+1]
			if IsAlphaDigit(next) {
				r.WriteRune(c)
				continue
			}
		}
		result, engTotal, digTotal = handleNumOrEngTotal(&r, result, engTotal, digTotal)
	}

	// 最后可能会有剩余未处理
	result, engTotal, digTotal = handleNumOrEngTotal(&r, result, engTotal, digTotal)
	if len(result) == 0 {
		result = nil
	}
	return
}
