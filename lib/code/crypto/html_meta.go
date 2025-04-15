package crypto

import (
	"html/template"
	"strings"
)

// Text 最大65535字符，约20000汉字

// 650个 rune字符  --> 1个汉字 len(content) 为3；而len([]rune(content))为1
func FilterPlain(content string, runeLen int) string {
	if content == "" {
		return ""
	}
	//runeLen := 650
	ct := []rune(content)
	if runeLen > 0 {
		if len(ct) > runeLen {
			ct = ct[:runeLen]
		}
	}

	return string(FilterPlainText(ct))
}
func FilterStr(content string, strlen int) string {
	if content == "" {
		return ""
	}
	if len(content) > strlen {
		content = content[:strlen]
	}
	return string(Filter([]rune(content), FilterSensitiveWords, FilterPlainText))
}

// 9999个 rune字符  --> 1个汉字 len(content) 为3；而len([]rune(content))为1
func FilterRawHtml(content template.HTML, runelen int, inject bool) template.HTML {
	s := string(content)
	if s == "" {
		return ""
	}
	if inject {
		s = strings.ReplaceAll(s, "<script ", "")
		s = strings.ReplaceAll(s, "<iframe ", "")
		s = strings.ReplaceAll(s, "<xml ", "")
		s = strings.ReplaceAll(s, "<html ", "")
	}

	if runelen > 0 {
		ct := []rune(s)
		if runelen > 0 && len(ct) > runelen {
			ct = ct[:runelen]
		}
		s = string(ct)
	}

	return template.HTML(s)
}
