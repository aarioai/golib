package stdfmt_test

import (
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

func TestGlobalTel(t *testing.T) {

	tests := map[string]string{
		"15000777962":          "15000777962",
		"15 000 777 962":       "15000777962",
		"18０ｏＯ７７７8８8":          "18000777888",
		"180 ｏｏ77 7888":        "18000777888",
		"（0755）－ 1234 567":     "0755-1234567",
		"（0755）－1234567":       "0755-1234567",
		"（0755）1234566":        "0755-1234566",
		"(0755）１23４566":        "0755-1234566",
		"(0755）- 1234566":      "0755-1234566",
		"0755 1234566":         "0755-1234566",
		"+86 7551234５67":       "+86 755-1234567",
		"ｏＯ８６ ０551234567":      "+86 551-234567",
		"(0086)755 12345678":   "+86 755-12345678",
		"+86-(0)755-12345678":  "+86 755-12345678",
		"+86-(0)755 12345678":  "+86 755-12345678",
		"+86-(755)12345678":    "+86 755-12345678",
		"+8675512345678":       "+86 755-12345678",
		"+86075512345678":      "+86 755-12345678",
		"+86(755)12345678":     "+86 755-12345678",
		"(+86)[755]12345678":   "+86 755-12345678",
		"(+86](0)755 12345678": "+86 755-12345678",
		"(0755）1234 5678":      "0755-12345678",
		"[0755）1234 5678":      "0755-12345678",
		"+86101234 5678":       "+86 10-12345678",
		"+8618０ｏＯ７７７8８8":       "+86 18000777888",
		"+860755123456789":     "+86 755-123456789",
	}
	for tel, want := range tests {
		got, err := stdfmt.ParseTel(tel)
		if err != nil || got.String() != want {
			t.Errorf("ParseTel(%q) got %s, want %q: %v", tel, got.String(), want, err)
		}
	}
}
