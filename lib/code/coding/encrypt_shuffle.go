package coding

import (
	"bytes"
	"github.com/aarioai/airis/pkg/arrmap"
	"math"
	"strings"
)

// ValidateShuffleEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
// 可控的人为错误，应该 panic
func ValidateShuffleEncryptKeys[T string | []byte](keys ...T) {
	for _, item := range keys {
		key := []byte(item)
		// 查看是否有重复字符
		newKey := arrmap.Compact(key, false)
		if len(newKey) != len(key) {
			panic("shuffle-encrypt key contains repeated characters")
		}
	}
}

// 混淆要求被混淆的字符，必须全部在 key里面。

func k2s(v int, shift int, base []byte) (s [2]byte) {
	l := len(base)
	shift %= l
	a := (v/l + shift) % l
	b := (v + shift) % l
	s[0] = base[a]
	s[1] = base[b]
	return
}
func s2k(s [2]byte, shift int, base []byte) int {
	l := len(base)
	shift %= l

	var (
		a, b int
	)
	for i, h := range base {
		if h == s[0] {
			a = i
		}
		if h == s[1] {
			b = i
		}
	}
	// 当发现解析后，跟之前结果不同。注意排查是不是 base 设置的不一样！
	x := ((a + l - shift) % l) * l
	y := (b + l - shift) % l
	return x + y
}

// 计算分组参数，返回每组大小m和组数groupNum
func calculateGroupParams(length, seed int) (m, groupNum int) {
	if length < 4 { // 至少4个字节才需要打散
		return 2, 1
	}

	// 计算每组大小，保证为偶数
	m = (seed%(length/2) + 1) & ^1 // 位操作确保偶数
	if m < 2 {
		m = 2
	}

	groupNum = length / m
	// 调整每组大小，确保组数大于等于每组大小
	for groupNum < m {
		m -= 2 // 保持偶数
		if m < 2 {
			m = 2
		}
		groupNum = length / m
	}

	return m, groupNum
}

// Scatter 将key打散，key最短8字节，最长应该不超过32位
// key 一般都很短，因此这里使用新内存，避免副作用
func Scatter(key []byte, seed int) []byte {
	if len(key) < 4 {
		return bytes.Clone(key)
	}
	m, groupNum := calculateGroupParams(len(key), seed)
	newKey := bytes.Clone(key)
	for i := 0; i < groupNum; i++ {
		im := i * m
		// 组内进行对称交换。已知m是偶数。而且这里i 一定能到达 m（因为上面是满组情况）
		// 因此可以安全使用 j<m/2
		for j := 0; j < m/2; j++ {
			newKey[im+j], newKey[im+m-1-j] = newKey[im+m-1-j], newKey[im+j]
		}
	}

	// 处理剩余部分
	pos := groupNum * m
	rest := len(newKey) - pos
	// 处理最后一组边界情况。依次往每组前面元素替换
	// rest 一定小于m，而 m 小于等于 groupNum
	for i := 0; i < rest; i++ {
		im := i * m
		newKey[pos+i], newKey[im] = newKey[im], newKey[pos+i]
	}
	return newKey
}

// Unscatter
// key 一般都很短，因此这里使用新内存，避免副作用
func Unscatter(scatteredKey []byte, seed int) []byte {
	if len(scatteredKey) < 4 {
		return bytes.Clone(scatteredKey)
	}
	m, groupNum := calculateGroupParams(len(scatteredKey), seed)
	key := bytes.Clone(scatteredKey)

	// 先处理剩余部分
	pos := groupNum * m
	rest := len(key) - pos

	// 处理最后一组边界情况。依次往每组前面元素替换
	// rest 一定小于m，而 m 小于等于 groupNum
	for i := 0; i < rest; i++ {
		im := i * m
		key[pos+i], key[im] = key[im], key[pos+i]
	}

	for i := 0; i < groupNum; i++ {
		im := i * m
		// 组内进行对称交换。已知m是偶数。而且这里i 一定能到达 m（因为上面是满组情况）
		// 因此可以安全使用 j<m/2
		for j := 0; j < m/2; j++ {
			key[im+j], key[im+m-1-j] = key[im+m-1-j], key[im+j]
		}
	}
	return key
}

// @param shift，是 innerToken 第一个字符
// @return 结果一定是偶数，会自动填充的
func ShuffleEncryptNumber(num uint64, shift int, base []byte) string {
	var b strings.Builder
	b.Grow(10)
	var x int
	var k [2]byte
	for num > 0 {
		//用2个字符代表3个数字，倒序
		x = int(num % 1000)
		num = num / 1000
		k = k2s(x, shift, base)
		b.WriteByte(k[0])
		b.WriteByte(k[1])
	}

	return b.String()
}

func ShuffleDecryptNumber(u string, shift int, base []byte) (v uint64) {
	l := len(u)
	if l%2 != 0 {
		return 0
	}
	// 一定是偶数
	for i := 0; i < l; i += 2 {
		s := [2]byte{u[i], u[i+1]}
		k := float64(s2k(s, shift, base)) * math.Pow(1000, float64(i)/2.0)
		v += uint64(k)
	}
	return v
}

// 预计算字符在base中的位置映射表
func createBaseIndexMap(base []byte) map[byte]int {
	indexMap := make(map[byte]int, len(base))
	for i, b := range base {
		indexMap[b] = i
	}
	return indexMap
}

// ShuffleEncrypt 逐一字符混淆
func ShuffleEncrypt(str []byte, shift int, key []byte) error {
	baseLen := len(key)
	if baseLen == 0 {
		return ErrEmptyCipherKey
	}

	// 预计算base字符索引
	baseIndex := createBaseIndexMap(key)

	strLen := len(str)
	half := strLen / 2
	remaining := strLen - half

	// 计算有效shift值
	shift = int(key[shift%baseLen]) % baseLen
	if shift < half {
		shift += half
	}

	// 处理中间字符的位置
	midPos := -1
	if v, ok := baseIndex[str[half]]; ok {
		midPos = v
	}

	// 处理两侧字符
	for i := 0; i < half; i++ {
		leftPos, ok := baseIndex[str[i]]
		if !ok {
			return ErrCipherKeyMissChar(str[i])
		}
		rightPos, ok := baseIndex[str[i+remaining]]
		if !ok {
			return ErrCipherKeyMissChar(str[i+remaining])
		}
		// 计算新位置
		leftShift := (i + shift) % baseLen
		rightShift := ((i + remaining) + shift) % baseLen

		x := key[(leftPos+leftShift)%baseLen]
		y := key[(rightPos+rightShift)%baseLen]

		// 根据中间字符决定交换方式
		if str[half]+x+y%2 == 0 {
			str[i], str[i+remaining] = x, y
		} else {
			str[i], str[i+remaining] = y, x
		}
	}

	// 处理中间字符（如果存在）
	if half < remaining && midPos >= 0 {
		shift = ((half + remaining) + shift) % baseLen
		str[half] = key[(midPos+shift)%baseLen]
	}

	return nil
}

// ShuffleDecrypt 解混淆
func ShuffleDecrypt(v []byte, shift int, base []byte) error {
	baseLen := len(base)
	if baseLen == 0 {
		return ErrEmptyCipherKey
	}

	// 预计算base字符索引
	baseIndex := createBaseIndexMap(base)

	strLen := len(v)
	half := strLen / 2
	remaining := strLen - half

	// 计算有效shift值
	shift = int(base[shift%baseLen]) % baseLen
	if shift < half {
		shift += half
	}

	// 处理中间字符的位置
	midPos := -1
	if v, ok := baseIndex[v[half]]; ok {
		midPos = v
	}

	// 处理两侧字符
	for i := 0; i < half; i++ {
		leftPos, ok := baseIndex[v[i]]
		if !ok {
			return ErrCipherKeyMissChar(v[i])
		}
		rightPos, ok := baseIndex[v[i+remaining]]
		if !ok {
			return ErrCipherKeyMissChar(v[i+remaining])
		}

		// 计算原始位置
		leftShift := ((i + remaining) + shift) % baseLen
		rightShift := (i + shift) % baseLen

		x := base[(baseLen+leftPos-leftShift)%baseLen]
		y := base[(baseLen+rightPos-rightShift)%baseLen]

		// 根据中间字符决定交换方式
		if v[half]+x+y%2 == 0 {
			v[i], v[i+remaining] = x, y
		} else {
			v[i], v[i+remaining] = y, x
		}
	}

	// 处理中间字符（如果存在）
	if half < remaining && midPos >= 0 {
		shift = ((half + remaining) + shift) % baseLen
		v[half] = base[(baseLen+midPos-shift)%baseLen]
	}

	return nil
}
