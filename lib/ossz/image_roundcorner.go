package ossz

import (
	"image"
	"image/color"
	"math"
)

// Settable Settable
type Settable interface {
	Set(x, y int, c color.Color)
}

var empty = color.RGBA{R: 255, G: 255, B: 255, A: 0}

// RoundCorner rate  0.5 就是 50%，也就是圆形
// 这个是可以实现，但是会有毛边，需要优化
func RoundCorner(m *image.Image, rate float64) {
	b := (*m).Bounds()
	w, h := b.Dx(), b.Dy()
	r := (float64(minValue(w, h)) / 2) * rate
	sm, ok := (*m).(Settable)
	if !ok {
		// Check if image is YCbCr format.
		ym, ok := (*m).(*image.YCbCr)
		if !ok {
			return
		}
		*m = yCbCrToRGBA(ym)
		sm = (*m).(Settable)
	}
	// Parallelize?
	for y := 0.0; y <= r; y++ {
		l := math.Round(r - math.Sqrt(2*y*r-y*y))
		for x := 0; x <= int(l); x++ {
			sm.Set(x-1, int(y)-1, empty)
		}
		for x := 0; x <= int(l); x++ {
			sm.Set(w-x, int(y)-1, empty)
		}
		for x := 0; x <= int(l); x++ {
			sm.Set(x-1, h-int(y), empty)
		}
		for x := 0; x <= int(l); x++ {
			sm.Set(w-x, h-int(y), empty)
		}
	}
}

func minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func yCbCrToRGBA(m image.Image) image.Image {
	b := m.Bounds()
	nm := image.NewRGBA(b)
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			nm.Set(x, y, m.At(x, y))
		}
	}
	return nm
}
