package stdfmt_test

import (
	"github.com/aarioai/golib/lib/code/stdfmt"
	"maps"
	"strings"
	"testing"
)

func TestStdNumbers(t *testing.T) {
	tests := map[string]string{
		"":         "", // 测试默认空值
		"－10０":     "-100",
		"—10Ｏ０":    "-1000",
		"—1oOｏＯ0０": "-1000000",
		"＋８６　ｏ７５５－１２３４５６７８9": "+86　0755-123456789",
		"＋1 07５5-1２3４567８9":  "+1 0755-123456789",
		"0755—123456789":     "0755-123456789",
	}
	for full, want := range tests {
		numsBefore := stdfmt.CloneFullWidthNumbersMap()
		got := stdfmt.ReplaceWithStdNumbers(full)
		if got != want {
			t.Errorf("ReplaceWithStdNumbers(%q) = %q, want %q", full, got, want)
			return
		}
		numsAfter := stdfmt.CloneFullWidthNumbersMap()
		if !maps.Equal(numsBefore, numsAfter) {
			t.Errorf("ReplaceWithStdASCII changed numbers map")
			return
		}
	}
}

func TestStdLowercases(t *testing.T) {
	tests := map[string]string{
		"":         "", // 测试默认空值
		"ｈｅｌｌｏ":    "hello",
		"Ｗｏｒｌｄ123": "Ｗorld123",
	}
	for full, want := range tests {
		lowsBefore := stdfmt.CloneFullWidthLowercasesMap()
		got := stdfmt.ReplaceWithStdLowercases(full)
		if got != want {
			t.Errorf("ReplaceWithStdLowercases(%q) = %q, want %q", full, got, want)
			return
		}
		lowsAfter := stdfmt.CloneFullWidthLowercasesMap()
		if !maps.Equal(lowsBefore, lowsAfter) {
			t.Errorf("ReplaceWithStdASCII changed lowercases map")
			return
		}
	}
}

func TestStdUppercases(t *testing.T) {
	tests := map[string]string{
		"":         "", // 测试默认空值
		"ｈＥＬＬＯ":    "ｈELLO",
		"ＷＯＲＬＤ１２３": "WORLD１２３",
	}
	for full, want := range tests {
		uppersBefore := stdfmt.CloneFullWidthUppercasesMap()
		got := stdfmt.ReplaceWithStdUppercases(full)
		if got != want {
			t.Errorf("ReplaceWithStdUppercases(%q) = %q, want %q", full, got, want)
			return
		}
		uppersAfter := stdfmt.CloneFullWidthUppercasesMap()
		if !maps.Equal(uppersBefore, uppersAfter) {
			t.Errorf("ReplaceWithStdASCII changed uppercases map")
			return
		}
	}
}

func TestStdSymbols(t *testing.T) {
	tests := map[string]string{
		"":              "", // 测试默认空值
		"Ｈｅｌｌｏ，１２３。":    "Ｈｅｌｌｏ, １２３. ",
		"Ｈeｌｌo，Ａario……": "Ｈeｌｌo, Ａario...",
	}
	for full, want := range tests {
		puncsBefore := stdfmt.CloneFullWidthSymbolsMap()
		breaksBefore := stdfmt.CloneFullWidthBreaks()
		got := stdfmt.ReplaceWithStdSymbols(full, false)
		if got != strings.ReplaceAll(want, " ", "") {
			t.Errorf("ReplaceWithStdUppercases(%q, false) = %q, want %q", full, got, want)
			return
		}

		got = stdfmt.ReplaceWithStdSymbols(full, true)
		if got != want {
			t.Errorf("ReplaceWithStdUppercases(%q, true) = %q, want %q", full, got, want)
			return
		}
		puncsAfter := stdfmt.CloneFullWidthSymbolsMap()
		breaksAfter := stdfmt.CloneFullWidthBreaks()
		if !maps.Equal(puncsBefore, puncsAfter) {
			t.Errorf("ReplaceWithStdASCII changed puncs map")
			return
		}
		if !maps.Equal(breaksBefore, breaksAfter) {
			t.Errorf("ReplaceWithStdASCII changed breaks map")
			return
		}
	}
}

func TestStdASCII(t *testing.T) {
	tests := map[string]string{
		"":           "", // 测试默认空值
		"Ｈｅｌｌｏ，１２３。": "Hello, 123. ",
		"ＷＯＲＬＤ——０１２３４５６７８９": "WORLD-0123456789",
		"Ｈeｌｌo，Ａario……":     "Hello, Aario...",
	}
	for full, want := range tests {
		numsBefore := stdfmt.CloneFullWidthNumbersMap()
		lowsBefore := stdfmt.CloneFullWidthLowercasesMap()
		uppersBefore := stdfmt.CloneFullWidthUppercasesMap()
		puncsBefore := stdfmt.CloneFullWidthSymbolsMap()
		breaksBefore := stdfmt.CloneFullWidthBreaks()
		got := stdfmt.ReplaceWithStdASCII(full, true)
		if got != want {
			t.Errorf("ReplaceWithStdASCII(%q, true) = %q, want %q", full, got, want)
		}
		numsAfter := stdfmt.CloneFullWidthNumbersMap()
		lowsAfter := stdfmt.CloneFullWidthLowercasesMap()
		uppersAfter := stdfmt.CloneFullWidthUppercasesMap()
		puncsAfter := stdfmt.CloneFullWidthSymbolsMap()
		breaksAfter := stdfmt.CloneFullWidthBreaks()

		if !maps.Equal(numsBefore, numsAfter) {
			t.Errorf("ReplaceWithStdASCII changed numbers map")
			return
		}

		if !maps.Equal(lowsBefore, lowsAfter) {
			t.Errorf("ReplaceWithStdASCII changed lowercases map")
			return
		}
		if !maps.Equal(uppersBefore, uppersAfter) {
			t.Errorf("ReplaceWithStdASCII changed uppercases map")
			return
		}
		if !maps.Equal(puncsBefore, puncsAfter) {
			t.Errorf("ReplaceWithStdASCII changed puncs map")
			return
		}
		if !maps.Equal(breaksBefore, breaksAfter) {
			t.Errorf("ReplaceWithStdASCII changed breaks map")
			return
		}

	}
}
