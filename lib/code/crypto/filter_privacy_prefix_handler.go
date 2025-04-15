package crypto

import (
	"github.com/aarioai/golib/data"
	"github.com/gobwas/glob/util/runes"
	"slices"
	"strings"
)

const minPrivacyLen = 2

// 国家前缀一般不用做隐私词，而地址、公司名（可能做隐私词）前面会带地区
func trimNation(ww *[]rune, prefix *strings.Builder) {
	if string((*ww)[0:2]) == "中国" {
		prefix.WriteString("中国")
		clear((*ww)[:2])
		*ww = (*ww)[2:]
	}
	// 去掉  xxx共和国 前缀
	x := runes.LastIndex(*ww, []rune("共和国"))
	if x > 0 && x < 6 && x < len(*ww)-minPrivacyLen-3 {
		x += 3
		prefix.WriteString(string((*ww)[:x]))
		clear((*ww)[:x])
		*ww = (*ww)[x:]
	}
}
func trimAutonomous(ww *[]rune, prefix *strings.Builder) {
	// 广西壮族自治区xxx自治州xxx市xxx区
	x := runes.LastIndex(*ww, []rune("自治"))
	if x > 0 && x < len(*ww)-minPrivacyLen-3 && strings.Index("区旗州县", string((*ww)[x+3])) > 0 {
		x += 3
		prefix.WriteString(string((*ww)[:x]))
		clear((*ww)[:x])
		*ww = (*ww)[x:]
	} else {
		x = runes.Index(*ww, []rune("特别行政区"))
		if x > 0 && x < len(*ww)-minPrivacyLen-5 {
			x += 5
			prefix.WriteString(string((*ww)[:x]))
			clear((*ww)[:x])
			*ww = (*ww)[x:]
		}
	}

}
func trimDist(ww *[]rune, prefix *strings.Builder, list []string, units ...string) int {
	if len(list) == 0 {
		return 0
	}
	j := len(list[0])

	if len(*ww) < j+minPrivacyLen {
		return 0
	}
	if slices.Contains(list, string((*ww)[:j])) {
		if len(units) > 0 {
			for _, unit := range units {
				m := len(unit)
				if m > 0 && len(*ww) > m+j && string((*ww)[j:m+j]) == unit {
					j += m
					break
				}
			}
		}

		prefix.WriteString(string((*ww)[:j]))
		clear((*ww)[:j])
		*ww = (*ww)[j:]
		return 1
	}
	return 0
}
func trimProvinces(ww *[]rune, prefix *strings.Builder) int {
	var h int
	h += trimDist(ww, prefix, data.Autonos3[:])
	h += trimDist(ww, prefix, data.Autonos2[:], "省") // 广西省
	h += trimDist(ww, prefix, data.Provinces3[:], "省")
	h += trimDist(ww, prefix, data.Provinces2[:], "省")
	h += trimDist(ww, prefix, data.Municities2[:], "市")
	h += trimDist(ww, prefix, data.SpecialRegions2[:], "特区")
	return h
}
func trimCities(ww *[]rune, prefix *strings.Builder, round int) int {
	h := round
	h += trimDist(ww, prefix, data.Cities4[:], "市")
	h += trimDist(ww, prefix, data.Cities3[:], "市")

	h += trimDist(ww, prefix, data.Leagues4[:], "盟", "市")
	h += trimDist(ww, prefix, data.Leagues3[:], "盟", "市")
	h += trimDist(ww, prefix, data.Leagues2[:], "盟", "市")

	h += trimDist(ww, prefix, data.Prefecture8[:], "州", "市")
	h += trimDist(ww, prefix, data.Prefecture4[:], "州", "市")
	h += trimDist(ww, prefix, data.Prefecture3[:], "州", "市")
	h += trimDist(ww, prefix, data.Prefecture2[:], "州", "市")

	if h > 1 || len(*ww) < minPrivacyLen+3 {
		return h
	}

	// 2个字以上的城市，在上面处理过了。这里只处理2个字的城市
	u := string((*ww)[2])
	if u == "市" {
		prefix.WriteString(string((*ww)[:3]))
		clear((*ww)[:3])
		*ww = (*ww)[3:]
		h++
	}

	// 处理过1次，还可以再处理一次（县级市）
	if h == 1 && round == 0 {
		return trimCities(ww, prefix, 1)
	}

	return h
}
func trimCounties(ww *[]rune, prefix *strings.Builder) int {
	if len(*ww) < 2+minPrivacyLen {
		return 0
	}
	if (*ww)[1] == '县' && strings.Index(data.Counties1, string((*ww)[0])) > -1 {
		prefix.WriteString(string((*ww)[:2]))
		clear((*ww)[:2])
		*ww = (*ww)[2:]
		return 1
	}
	var h int
	h += trimDist(ww, prefix, data.AutonoBanners4[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.AutonoBanners3[:], "旗", "县", "市")

	h += trimDist(ww, prefix, data.Banners8[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.Banners6[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.Banners5[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.Banners4[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.Banners3[:], "旗", "县", "市")
	h += trimDist(ww, prefix, data.Banners2[:], "旗", "县", "市")

	h += trimDist(ww, prefix, data.AutonoCounties5[:], "县", "区", "市")
	h += trimDist(ww, prefix, data.AutonoCounties4[:], "县", "区", "市")
	h += trimDist(ww, prefix, data.AutonoCounties3[:], "县", "区", "市")
	h += trimDist(ww, prefix, data.AutonoCounties2[:], "县", "区", "市")

	h += trimDist(ww, prefix, data.Counties4[:], "县", "区", "市")
	h += trimDist(ww, prefix, data.Counties3[:], "县", "区", "市")

	h += trimDist(ww, prefix, data.Districts5[:], "区")
	h += trimDist(ww, prefix, data.Districts4[:], "区")
	h += trimDist(ww, prefix, data.Districts3[:], "区")

	return h

}
func DefaultPrivacyPrefixHandler(ww []rune, typ PrivacySuffixType) ([]rune, string) {
	l := len(ww)
	if l <= minPrivacyLen {
		return ww, ""
	}
	var prefix strings.Builder
	// 必须用指针，否则由于地址不同，len(ww) 不同（即使  string(ww) 显示都一样
	trimNation(&ww, &prefix)
	trimAutonomous(&ww, &prefix)
	trimProvinces(&ww, &prefix)
	trimCities(&ww, &prefix, 0)
	trimCounties(&ww, &prefix)
	return ww, prefix.String()
}
