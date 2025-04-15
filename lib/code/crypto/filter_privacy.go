package crypto

import (
	"fmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type PrivacySuffixType uint8

const (
	PrivacySuffixCompany PrivacySuffixType = 1 // 公司后缀
)

func ParsePrivacy(content string, words []string) []string {
	r := regexp.MustCompile(`(?i)<privacy>([^<]+)</privacy>`) // (?i)  忽略大小写
	matches := r.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return words
	}
	if words == nil {
		words = make([]string, 0, len(matches))
	} else if len(words) > 1 {
		words = arrmap.Compact(words, false)
	}

	for _, v := range matches {
		var exists bool
		for _, word := range words {
			if word == v[1] {
				exists = true
			}
		}
		if !exists {
			words = append(words, v[1])
		}
	}

	sort.Strings(words) // 必须要排序，这样才能正确替换
	// 移除空格
	var j = -1
	for i, word := range words {
		if word != "" {
			j = i
			break
		}
	}
	if j == -1 {
		return nil
	} else if j > -1 {
		words = words[j:]
	}
	return words
}

// 由于仲裁文书、文章等可能会分段，因此返回 replacer 方便多次替换
// @warn 使用前，一定要对 words 进行 sort，否则可能出现替换不一致的情况
func PrivacyReplacer(words []string, prefixHandler func([]rune, PrivacySuffixType) ([]rune, string)) *strings.Replacer {
	if len(words) == 0 {
		return nil
	}
	corps2 := []string{"公司", "企业", "商行", "集团", "中心"}
	corps4 := []string{"有限公司", "合伙企业"}
	corps6 := []string{"有限责任公司", "股份有限公司", "股份合作公司", "科技有限公司", "网络有限公司", "信息有限公司", "技术有限公司", "智能有限公司", "商务有限公司", "电商有限公司", "传媒有限公司", "教育有限公司"}

	pws := make([]string, len(words)*2)
	for i, w := range words {

		ww := []rune(w)
		l := len(ww)
		if l <= 1 {
			pws[i*2] = w
			// <u> 标签 HTML5 变更 https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/u
			// 有可能外面有 <privacy> ，也可能没有（如仲裁文书当事人信息自动设为隐私词）因此这里使用 <abbr>xxx<fuzzy>x:x</fuzzy></abbr> 的格式
			pws[i*2+1] = fmt.Sprintf(`<abbr data-privacy-key="%d">%s</abbr>`, i, strings.Repeat("※", l))
			continue
		}

		// 处理后缀
		var suffix string
		var m int
		if l > 6 && slices.Contains(corps6, string(ww[l-6:l])) {
			m = 6
		} else if l > 4 && slices.Contains(corps4, string(ww[l-4:l])) {
			m = 4
		} else if slices.Contains(corps2, string(ww[l-2:l])) {
			m = 2
		}

		if m > 0 {
			b := l - m
			if b > 0 {
				suffix = string(ww[b:l])
				ww = ww[0:b]
				l = len(ww)
			}
		}

		// 处理前缀
		var prefix string
		if prefixHandler == nil {
			prefixHandler = DefaultPrivacyPrefixHandler
		}

		ww, prefix = prefixHandler(ww, PrivacySuffixCompany)
		l = len(ww)
		if l == 0 {
			pws[i*2] = w
			// <u> 标签 HTML5 变更 https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/u
			// 有可能外面有 <privacy> ，也可能没有（如仲裁文书当事人信息自动设为隐私词）因此这里使用 <abbr>xxx<fuzzy>x:x</fuzzy></abbr> 的格式
			pws[i*2+1] = fmt.Sprintf(`<abbr data-privacy-key="%d">%s</abbr>`, i, strings.Repeat("※", l))
			continue
		}

		var public string
		n := l
		if l > 1 {
			// 公开1/3长度
			k := l / 3
			if k == 0 {
				k = 1
			}
			public = string(ww[0:k])
			n = l - k
		}
		pws[i*2] = w
		// <u> 标签 HTML5 变更 https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/u
		// 有可能外面有 <privacy> ，也可能没有（如仲裁文书当事人信息自动设为隐私词）因此这里使用 <abbr>xxx<fuzzy>x:x</fuzzy></abbr> 的格式
		pws[i*2+1] = fmt.Sprintf(`<abbr data-privacy-key="%d">%s%s%s%s</abbr>`, i, prefix, public, strings.Repeat("※", n), suffix)
	}
	return strings.NewReplacer(pws...)
}

// @warn 使用前，一定要对 words 进行 sort，否则可能出现替换不一致的情况
func NoPrivacy(content string, words []string) (string, error) {
	if len(words) == 0 {
		return content, nil
	}
	content = strings.ReplaceAll(content, "<privacy>", "")
	content = strings.ReplaceAll(content, "</privacy>", "")
	r := regexp.MustCompile(`(?i)<abbr\s+data-privacy-key="(\d+)">[^<]*</abbr>`) // (?i)  忽略大小写
	matches := r.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return content, nil
	}
	done := make(map[int]struct{}, len(matches))
	var err error
	for _, v := range matches {
		i, _ := strconv.Atoi(v[1])
		if _, ok := done[i]; ok {
			continue
		}
		if i > len(words) {
			err = fmt.Errorf("unreplace privacy words index %d is out of range", i)
			continue
		}
		content = strings.ReplaceAll(content, v[0], "<privacy>"+words[i]+"</privacy>")
		done[i] = struct{}{}
	}

	return content, err
}
