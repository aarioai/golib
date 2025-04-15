package stdfmt

import (
	"errors"
	"maps"
	"strings"
)

// 全角字符
const (
	FullWidthNumbers    = "０１２３４５６７８９"
	FullWidthLowercases = "ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ"
	FullWidthUppercases = "ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ"
)

// 禁止全局共享slice！！！不可避免其他人参与编程，新手失误导致修改slice。
// 即使使用局部共享slice，对slice的操作，必须要增加测试用例，判断使用前后是否发生变化。
var (
	// @test CloneFullWidthNumbersMap
	fullWidthNumbersMap = createRuneMap(FullWidthNumbers, numbersSlice)
	// @test CloneFullWidthLowercasesMap
	fullWidthLowercasesMap = createRuneMap(FullWidthLowercases, lowercasesSlice)
	// @test CloneFullWidthUppercasesMap
	fullWidthUppercasesMap = createRuneMap(FullWidthUppercases, uppercasesSlice)
	// @test CloneFullWidthSymbolsMap
	fullWidthSymbolsMap = map[rune]rune{
		// 特殊中文符号
		'—': '-',
		'　': ' ',
		'“': '"',
		'”': '"',
		'‘': '\'',
		'’': '\'',
		'。': '.',
		'【': '[',
		'、': ',',
		'】': ']',
		'「': '[',
		'」': ']',
		'﹝': '[',
		'﹞': ']',

		// 33-47
		'！': '!',
		'＂': '"', // 这里是英文输入法下，全角双引号。不同于中文输入法的双引号“”
		'＃': '#',
		'＄': '$',
		'％': '%',
		'＆': '&',
		'＇': '\'',
		'（': '(',
		'）': ')',
		'＊': '*',
		'＋': '+',
		'，': ',',
		'－': '-',
		'．': '.',
		'／': '/',

		// 58-64
		'：': ':',
		'；': ';',
		'＜': '<',
		'＝': '=',
		'＞': '>',
		'？': '?',
		'＠': '@',

		// 91-96
		'［': '[',
		'＼': '\\',
		'］': ']',
		'＾': '^',
		'＿': '_',
		'｀': '`',

		// 123-126
		//'｛': '{',
		'｜': '|',
		//'｝': '}',
		'～': '~',
	}
	// @test CloneFullWidthBreaks
	fullWidthBreaks = map[string]string{
		"‘": " '",
		"’": "' ",
		"“": " \"",
		"”": "\" ",
		"，": ", ",
		"。": ". ",
		"．": ". ",
		"：": ": ",
		"；": "; ",
		"、": ", ",
		"？": "? ",
		"！": "! ",
	}
)

// ReplaceToStdASCII 全部只有ASCII字符或其全角模式字符，并替换全角字符为半角字符。若有其他字符，则报错
func ReplaceToStdASCII(s string, withBreakSpace bool) (string, error) {
	return ReplaceWithStdASCIIFunc(s, withBreakSpace, stdASCIIExtRuneHandler)
}

// ReplaceToStdNumbers 全部只有数字+-或其全角模式字符，并替换全角字符为半角字符。若有其他字符，则报错
func ReplaceToStdNumbers(s string, withBreakSpace bool) (string, error) {
	return ReplaceWithStdNumbersFunc(s, stdNumbersExtRuneHandler)
}

func stdASCIIExtRuneHandler(r rune) (rune, error) {
	if r > 255 {
		return r, errors.New("invalid ASCII byte: " + string(r))
	}
	return r, nil
}
func stdNumbersExtRuneHandler(r rune) (rune, error) {
	if (r >= '0' && r <= '9') || r == '-' || r == '+' {
		return r, nil
	}
	return r, errors.New("invalid number byte: " + string(r))
}

// ReplaceWithStdASCII 替换全角字符为半角字符，不替换其他字符
// to their half-width equivalents.
// @test CloneFullWidthNumbersMap CloneFullWidthLowercasesMap CloneFullWidthUppercasesMap
// @test CloneFullWidthSymbolsMap CloneFullWidthBreaks
func ReplaceWithStdASCII(s string, withBreakSpace bool) string {
	result, _ := ReplaceWithStdASCIIFunc(s, withBreakSpace, nil)
	return result
}
func ReplaceWithStdASCIIFunc(s string, withBreakSpace bool, extRuneHandler func(rune) (rune, error)) (string, error) {
	if s == "" {
		return s, nil
	}
	s = strings.ReplaceAll(s, "——", "-")
	s = strings.ReplaceAll(s, "……", "...")
	s = strings.ReplaceAll(s, "…", "...")
	if withBreakSpace {
		for from, to := range fullWidthBreaks {
			s = strings.ReplaceAll(s, from, to)
		}
	}
	var err error
	result := strings.Map(func(r rune) rune {
		if half, ok := fullWidthNumbersMap[r]; ok {
			return half
		}
		if half, ok := fullWidthLowercasesMap[r]; ok {
			return half
		}
		if half, ok := fullWidthUppercasesMap[r]; ok {
			return half
		}

		if half, ok := fullWidthSymbolsMap[r]; ok {
			return half
		}
		if extRuneHandler != nil {
			var newErr error
			if r, newErr = extRuneHandler(r); newErr != nil && err == nil {
				err = newErr
			}
		}
		return r
	}, s)
	return result, err
}

// ReplaceWithStdNumbers 替换全角数字和负号为半角，包括正负符号。不替换其他字符
// 因此可以适用于国际电话号码，如 +1 102343
// @test CloneFullWidthNumbersMap
func ReplaceWithStdNumbers(s string) string {
	result, _ := ReplaceWithStdNumbersFunc(s, nil)
	return result
}
func ReplaceWithStdNumbersFunc(s string, extRuneHandler func(rune) (rune, error)) (string, error) {
	if s == "" {
		return s, nil
	}
	var err error
	result := strings.Map(func(r rune) rune {
		// 一个是英文全角字符，一个是中文扩展符
		if r == '－' || r == '—' {
			return '-'
		}
		if r == '＋' {
			return '+'
		}
		// 特别的符号 字母oO，或者ｏＯ
		if r == 'o' || r == 'O' || r == 'ｏ' || r == 'Ｏ' {
			return '0'
		}
		if half, ok := fullWidthNumbersMap[r]; ok {
			return half
		}
		if extRuneHandler != nil {
			var newErr error
			if r, newErr = extRuneHandler(r); newErr != nil && err == nil {
				err = newErr
			}
		}
		return r
	}, s)
	return result, err
}

// ReplaceWithStdLowercases 替换全角小写字母为半角字符
// @test CloneFullWidthLowercasesMap
func ReplaceWithStdLowercases(s string) string {
	result, _ := ReplaceWithHalfWidthFunc(s, fullWidthLowercasesMap, nil)
	return result
}
func ReplaceWithStdLowercasesFunc(s string, extRuneHandler func(rune) (rune, error)) (string, error) {
	return ReplaceWithHalfWidthFunc(s, fullWidthLowercasesMap, extRuneHandler)
}

// ReplaceWithStdUppercases 替换全角大写字母为半角字符
// @test CloneFullWidthUppercasesMap
func ReplaceWithStdUppercases(s string) string {
	result, _ := ReplaceWithHalfWidthFunc(s, fullWidthUppercasesMap, nil)
	return result
}
func ReplaceWithStdUppercasesFunc(s string, extRuneHandler func(rune) (rune, error)) (string, error) {
	return ReplaceWithHalfWidthFunc(s, fullWidthUppercasesMap, extRuneHandler)
}

// ReplaceWithStdSymbols 替换全角英文符号为半角字符
// @test CloneFullWidthSymbolsMap CloneFullWidthBreaks
func ReplaceWithStdSymbolsFunc(s string, withBreakSpace bool, extRuneHandler func(rune) (rune, error)) (string, error) {
	if s == "" {
		return s, nil
	}
	s = strings.ReplaceAll(s, "——", "-")
	s = strings.ReplaceAll(s, "……", "...")
	s = strings.ReplaceAll(s, "…", "...")
	if withBreakSpace {
		for from, to := range fullWidthBreaks {
			s = strings.ReplaceAll(s, from, to)
		}
	}
	var err error
	result := strings.Map(func(r rune) rune {
		if half, ok := fullWidthSymbolsMap[r]; ok {
			return half
		}
		if extRuneHandler != nil {
			var newErr error
			if r, newErr = extRuneHandler(r); newErr != nil && err == nil {
				err = newErr
			}
		}
		return r
	}, s)
	return result, err
}

func ReplaceWithStdSymbols(s string, withBreakSpace bool) string {
	result, _ := ReplaceWithStdSymbolsFunc(s, withBreakSpace, nil)
	return result
}

func ReplaceWithHalfWidthFunc(s string, dict map[rune]rune, extRuneHandler func(rune) (rune, error)) (string, error) {
	if s == "" || len(dict) == 0 {
		return s, nil
	}
	var err error
	result := strings.Map(func(r rune) rune {
		if half, ok := dict[r]; ok {
			return half
		}
		if extRuneHandler != nil {
			var newErr error
			if r, newErr = extRuneHandler(r); newErr != nil && err == nil {
				err = newErr
			}
		}
		return r
	}, s)
	return result, err
}
func createRuneMap[T byte | rune](keys string, values []T) map[rune]rune {
	keysRunes := []rune(keys)
	if len(keysRunes) != len(values) {
		panic("keys and values must have the same length")
	}
	m := make(map[rune]rune, len(keysRunes))
	for i, k := range keysRunes {
		m[k] = rune(values[i])
	}
	return m
}

func CloneFullWidthNumbersMap() map[rune]rune {
	return maps.Clone(fullWidthNumbersMap)
}
func CloneFullWidthLowercasesMap() map[rune]rune {
	return maps.Clone(fullWidthLowercasesMap)
}
func CloneFullWidthUppercasesMap() map[rune]rune {
	return maps.Clone(fullWidthUppercasesMap)
}
func CloneFullWidthSymbolsMap() map[rune]rune {
	return maps.Clone(fullWidthSymbolsMap)
}
func CloneFullWidthBreaks() map[string]string {
	return maps.Clone(fullWidthBreaks)
}
