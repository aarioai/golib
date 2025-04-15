package enumz

import (
	"github.com/mssola/useragent"
	"strings"
)

type plt uint8

const (
	unknownPlt   plt = 0
	pltAndroid   plt = 1
	pltMacintosh plt = 2
	pltIphone    plt = 3
	pltIpad      plt = 4
	pltWindows   plt = 5
)

func osType(ag *useragent.UserAgent) plt {
	if strings.Index(ag.OS(), "Android") > -1 {
		return pltAndroid
	}
	switch ag.Platform() {
	case "Macintosh":
		return pltMacintosh
	case "iPhone":
		return pltIphone
	case "iPad":
		return pltIpad
	case "Windows":
		return pltWindows
	}

	return unknownPlt
}
