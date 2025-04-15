package crypto

import (
	"html/template"
	"strings"
)

func ReplaceHtml(rep *strings.Replacer, s template.HTML) template.HTML {
	if rep == nil {
		return s
	}
	return template.HTML(rep.Replace(string(s)))
}
