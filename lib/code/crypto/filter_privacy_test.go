package crypto_test

import (
	"github.com/aarioai/golib/lib/code/crypto"
	"golang.org/x/exp/slices"
	"html/template"
	"testing"
)

func TestPrivacy(t *testing.T) {
	content := `<privacy><p>我爱你</p><Privacy>上海市某某科技有限公司</privacy><p>哈哈</p><privacy>北京市王德发</privacy><privacy>律需网</privacy><p>北京市王德发</p>律需网<privacy>上海市某某科技有限公司</privacy>北京市王德发`
	words := crypto.ParsePrivacy(content, nil)
	if !slices.Equal(words, []string{"上海市某某科技有限公司", "北京市王德发", "律需网"}) {
		t.Errorf("crypto.ParsePrivacy() not passed %v", words)
		return
	}
}
func TestPrivacyReplacer(t *testing.T) {
	content := `<privacy><p>我爱你</p><Privacy>上海市某某科技有限公司</privacy><p>哈哈</p><privacy>北京市王德发</privacy><privacy>律需网</privacy><p>北京市王德发</p>律需网<privacy>上海市某某科技有限公司</privacy>北京市王德发`

	words := crypto.ParsePrivacy(content, nil)
	replacer := crypto.PrivacyReplacer(words, nil)
	if replacer == nil {
		t.Errorf("crypto.PrivacyReplacer nil")
		return
	}

	html := crypto.ReplaceHtml(replacer, template.HTML(content))
	suppose := `<privacy><p>我爱你</p><Privacy><abbr data-privacy-key="0">上海市某※科技有限公司</abbr></privacy><p>哈哈</p><privacy><abbr data-privacy-key="1">北京市王※※</abbr></privacy><privacy><abbr data-privacy-key="2">律※※</abbr></privacy><p><abbr data-privacy-key="1">北京市王※※</abbr></p><abbr data-privacy-key="2">律※※</abbr><privacy><abbr data-privacy-key="0">上海市某※科技有限公司</abbr></privacy><abbr data-privacy-key="1">北京市王※※</abbr>`
	if string(html) != suppose {
		t.Errorf("crypto.PrivacyReplacer error `%s`", html)
		return
	}

	// 接管前缀，不再过滤地区
	handler := func(word []rune, typ crypto.PrivacySuffixType) ([]rune, string) {
		var prefix string

		return word, prefix
	}
	replacer = crypto.PrivacyReplacer(words, handler)
	html = crypto.ReplaceHtml(replacer, template.HTML(content))
	suppose = `<privacy><p>我爱你</p><Privacy><abbr data-privacy-key="0">上※※※※科技有限公司</abbr></privacy><p>哈哈</p><privacy><abbr data-privacy-key="1">北京※※※※</abbr></privacy><privacy><abbr data-privacy-key="2">律※※</abbr></privacy><p><abbr data-privacy-key="1">北京※※※※</abbr></p><abbr data-privacy-key="2">律※※</abbr><privacy><abbr data-privacy-key="0">上※※※※科技有限公司</abbr></privacy><abbr data-privacy-key="1">北京※※※※</abbr>`
	if string(html) != suppose {
		t.Errorf("crypto.PrivacyReplacer with handler error `%s`", html)
		return
	}

	testPrivacyUnreplacer(t, words, string(html))
}

func testPrivacyUnreplacer(t *testing.T, words []string, encContent string) {
	// 移除掉开头无效 <privacy>
	rawContent := `<p>我爱你</p><Privacy><privacy>上海市某某科技有限公司</privacy><p>哈哈</p><privacy>北京市王德发</privacy><privacy>律需网</privacy><p><privacy>北京市王德发</privacy></p><privacy>律需网</privacy><privacy>上海市某某科技有限公司</privacy><privacy>北京市王德发</privacy>`
	content, err := crypto.NoPrivacy(encContent, words)
	if content != rawContent {
		t.Errorf("crypto.privacyUnreplacer fail `%s` %v", content, err)
		return
	}
}
