package crypto

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type EmailCipher SidCipher

const (
	EmailCipherDotCom    = byte('~') // 除了分隔符以外其他符号，都可以
	EmailCipherSeparator = byte(' ')
	EmailCipherKeyMinLen = 66 // 邮箱支持的全部字符， [\w-@.]
)

var (
	emailCipherDotCom    = string(EmailCipherDotCom)
	emailCipherSeparator = string(EmailCipherSeparator)
)

var domainMap = map[string]string{
	"A": "aol.com",
	"B": "baidu.com",
	"m": "mi.com",
	"e": "xiaomi.com",
	"b": "bytedance.com",
	"G": "gmail.com",
	"O": "outlook.com",
	"Y": "yahoo.com",
	"I": "icloud.com",
	"M": "mail.com",
	"T": "tom.com",
	"H": "hotmail.com",
	"a": "aliyun.com",
	"t": "tencent.com",
	"Q": "qq.com",
	"1": "126.com",
	"3": "163.com",
	"2": "21cn.com",
	"S": "sina.com",
	"s": "sina.com.cn",
	"u": "sohu.com",
	"L": "luexu.com",
}

// ValidateTelEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
// 支持 除分隔符（crypto.TelCipherSeparator）外，其他所有字符。但是要求避免把types.Base64Digits全部字符都包括在内
func ValidateEmailEncryptKeys[T string | []byte](keys ...T) {
	// 基于 SidCipher
	ValidateSidEncryptKeys(keys...)

	for _, item := range keys {
		key := []byte(item)
		// 不能包含分隔符
		if bytes.Contains(key, []byte{EmailCipherDotCom, EmailCipherSeparator}) {
			panic(fmt.Sprintf("email encrypt key %s is invalid", string(key)))
		}
	}
}

// ShuffleEncryptEmail
// email 用户名前1位和最后1位是明文
// @test 保证不修改 key
func ShuffleEncryptEmail(email string, key []byte, scatter bool) (EmailCipher, error) {
	ValidateKeyLen(key, EmailCipherKeyMinLen)
	if email == "" {
		return "", nil
	}
	username, domain, err := parseEmailForEncrypt(email, key)
	if err != nil {
		return "", err
	}
	sep := string(EmailCipherSeparator)
	if len(domain) == 1 {
		sep = ""
	}
	cryptogram, err := ShuffleEncryptSid(username, key, scatter)
	if err != nil {
		return "", err
	}

	var s strings.Builder
	s.Grow(1 + len(cryptogram) + len(sep) + len(domain))
	s.WriteString(cryptogram.String())
	s.WriteString(sep)
	s.WriteString(domain)
	return EmailCipher(s.String()), nil
}

func (c EmailCipher) Decrypt(key []byte) (string, error) {
	ValidateKeyLen(key, EmailCipherKeyMinLen)
	if c == "" {
		return "", nil
	}
	cryptogram, domain, err := parseEmailForDecrypt(c)
	if err != nil {
		return "", err
	}
	username, err := SidCipher(cryptogram).Decrypt(key)
	if err != nil {
		return "", err
	}
	return username + "@" + domain, nil
}

func (c EmailCipher) String() string {
	return string(c)
}
func (c EmailCipher) Desensitize(usernameWantLen ...int) (string, error) {
	if c == "" {
		return "", nil
	}
	ciphertext, domain, err := parseEmailForDecrypt(c)
	if err != nil {
		return "", err
	}
	text := SidCipher(ciphertext).Desensitize(usernameWantLen...)
	return text + "@" + domain, nil
}

func parseEmailForEncrypt(email string, key []byte) (string, string, error) {
	if email == "" || len(key) == 0 {
		return "", "", fmt.Errorf("email or key is empty")
	}

	if bytes.IndexByte(key, EmailCipherSeparator) > -1 {
		return "", "", errors.New("key contains invalid character " + emailCipherSeparator)
	}
	arr := strings.Split(email, "@")
	if len(arr) != 2 {
		return "", "", errors.New("invalid email")
	}
	username := arr[0] // 保留大小写规则
	domain := strings.ToLower(arr[1])

	for alias, domainAlias := range domainMap {
		if domain == domainAlias {
			domain = alias
			break
		}
	}
	if len(domain) > 1 {
		domain = strings.ReplaceAll(domain, ".com", emailCipherDotCom)
	}
	return username, domain, nil
}
func parseEmailForDecrypt(c EmailCipher) (string, string, error) {
	arr := strings.Split(c.String(), emailCipherSeparator)
	switch len(arr) {
	case 1:
		a := arr[0]
		last := len(a) - 1
		domain, ok := domainMap[string(a[last])]
		if !ok {
			return "", "", fmt.Errorf("invalid email cipher: %s", c)
		}
		username := a[0:last]
		return username, domain, nil
	case 2:
		username := arr[0] // 保留大小写规则
		domain := strings.ReplaceAll(arr[1], emailCipherDotCom, ".com")
		return username, domain, nil
	default:
		return "", "", fmt.Errorf("invalid email cipher: %s", string(c))
	}

}
