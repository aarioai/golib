package crypto

// 不能用string, 必须要用 []rune，否则中文切片 [:]会出不对
func Filter(content []rune, filters ...func([]rune) []rune) []rune {
	for _, filter := range filters {
		content = filter(content)
	}
	return content
}

var plainTextEncodes = map[rune][]rune{
	'<':      []rune("&lt;"),
	'>':      []rune("&gt;"),
	'"':      []rune("&#34;"),
	rune(39): []rune("&#39;"), // 单引号
}

// 不能用string, 必须要用 []rune，否则中文切片 [:]会出不对
// 允许颜文字、\n符号
func FilterPlainText(content []rune) []rune {
	// 替换  < >
	var i int
	n := len(content)
	for ; i < n; i++ {
		c := content[i]
		if ec, ok := plainTextEncodes[c]; ok {
			b := []rune(string(content[i+1:]))
			content = append(content[:i], ec...)
			content = append(content, b...)
			i += len(ec)
			n = len(content)
		}
	}
	return content
}
