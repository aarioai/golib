package ossz

import (
	"bytes"
	"errors"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
)

// FitQrInnerIcon
// @param cornerRate < 0 表示自动计算圆角度；
func FitQrInnerIcon(qrSize float64, icon string, padding int, bgc color.Color, cornerRate float64) (image.Image, error) {
	if padding < 0 {
		if qrSize < 100.0 {
			// f(x) = 0.02x + 1
			padding = int(math.Round(.02*qrSize)) + 1
		} else {
			// f(x) = sqrt(1.3x) / 5   // 大致公式
			padding = int(math.Round(math.Sqrt(1.3*qrSize) / 5.0))
		}
	}
	if bgc == nil {
		bgc = color.Alpha{A: 200} // 80%透明度
	}
	if cornerRate < 0 {
		// f(x) = (sqrt4(7x) + 10) / 100 // 大致公式，sqrt4() 开四次方根
		cornerRate = math.Round(math.Pow(7.0*qrSize, 0.25)+10.0) / 100.0
	}
	iw := uint(math.Ceil(qrSize / 4.32)) // 大概二维码尺寸:logo保持4.32:1

	iconr, _, err := LoadImageReader(icon)
	if err != nil {
		return nil, err
	}
	defer iconr.Close()

	var iconim image.Image
	// 不要直接用 image.Decode，要不然会报错：”image: unknown format“
	// 这里对一些不规范的JPEP图片会报 · create utac url qr img: redis error: invalid JPEG format: missing SOI marker
	// 针对这种情况，换png尝试
	iconr.Seek(0, 0)
	if iconim, err = jpeg.Decode(iconr); err != nil {
		iconr.Seek(0, 0)
		iconim, err = png.Decode(iconr)
		if err != nil {
			iconr.Seek(0, 0)
			iconim, _, err = image.Decode(iconr)
		}
	}
	if err != nil {
		return iconim, errors.New("image Decode: " + err.Error())
	}

	b := iconim.Bounds()
	if b.Dx() > int(iw) {
		iconim = resize.Resize(iw, 0, iconim, resize.Lanczos3)
	}

	b = iconim.Bounds()
	w := b.Dx()
	h := b.Dy()
	if padding > 0 {
		p := padding + padding // 两边
		w += p
		h += p
	}

	// 设置圆角
	RoundCorner(&iconim, cornerRate)

	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	if padding > 0 && bgc != color.Transparent {
		draw.Draw(m, m.Bounds(), image.NewUniform(bgc), image.Point{}, draw.Src) // logo 背景白
		g := nrgbaToImage(m, m.Bounds().Dx(), m.Bounds().Dy())
		RoundCorner(&g, cornerRate)
	}
	draw.Draw(m, b.Add(image.Pt(padding, padding)), iconim, image.Point{}, draw.Over)
	return m, nil
}
func nrgbaToImage(m *image.NRGBA, width, height int) image.Image {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.Set(x, y, m.At(x, y))
		}
	}
	return m
}

// 格子数
func QrGridNum(level qrcode.RecoveryLevel) uint16 {
	switch level {
	case qrcode.Low:
		return 29
	case qrcode.Medium:
		return 29
	case qrcode.High:
		return 33
	case qrcode.Highest:
		return 37
	}
	return 33
}
func softCeil(gridNum float64) float64 {
	n := math.Floor(gridNum)
	g := gridNum - n
	if g < 0.5 {
		return n + 0.5
	}
	return math.Ceil(gridNum)
}
func softFloor(gridNum float64) float64 {
	n := math.Floor(gridNum)
	g := gridNum - n
	if g > 0.5 {
		return n + 0.5
	}
	return n
}

/*
*
QR 用处：1、URL；   2、密钥
*/
func DrawQrCode(content string, side uint16, level qrcode.RecoveryLevel, icon image.Image, qrColor, bg color.Color) ([]byte, error) {
	sz := int(side)
	gridNum := QrGridNum(level) // 二维码每行格子数

	// There are four levels of error recovery: Low, medium, high and highest.
	// QR Codes with a higher recovery level are more robust to damage, at the cost of being physically larger.
	if qrColor == nil {
		qrColor = color.Black
	}
	if bg == nil {
		bg = color.White
	}

	q, err := qrcode.New(content, level)
	if err != nil {
		return nil, err
	}
	q.BackgroundColor = bg // 背景色 透明
	q.ForegroundColor = qrColor
	q.DisableBorder = true // 禁止白色边框（全局二维码）

	imBytes, err := q.PNG(sz)
	if err != nil {
		return nil, err
	}
	if icon == nil {
		return imBytes, nil
	}

	ib := icon.Bounds()
	x := (sz - ib.Dx()) / 2
	y := (sz - ib.Dy()) / 2
	offset := image.Pt(x, y)

	ir := bytes.NewReader(imBytes)
	ir.Seek(0, 0) // 光标必须要移到开头
	qrim, err := png.Decode(ir)
	if err != nil {
		return nil, err
	}
	b := qrim.Bounds()

	m := image.NewNRGBA(b)
	draw.Draw(m, b, qrim, image.Point{}, draw.Src) // 二维码

	// 抠掉中心部分
	// 二维码为33格，要保证抠掉的部分正好把格子扣完。否则会出现留边不好看
	gridWidth := float64(sz) / float64(gridNum)
	x0 := int(math.Floor(softFloor(float64(x)/gridWidth) * gridWidth))
	y0 := int(math.Floor(softFloor(float64(y)/gridWidth) * gridWidth))
	x1 := int(math.Ceil(softCeil(float64(x+ib.Dx())/gridWidth) * gridWidth))
	y1 := int(math.Ceil(softCeil(float64(y+ib.Dy())/gridWidth) * gridWidth))
	cb := image.Rect(x0, y0, x1, y1)
	draw.Draw(m, cb, image.Transparent, image.Point{}, draw.Src)
	// ，并边角去掉1像素，模拟圆角
	//cb := image.Rect(x0+2, y0, x1-2, y1-2)
	//draw.Draw(m, cb, image.Transparent, image.Point{}, draw.Src)
	//cbl := image.Rect(x0, y0+2, x0+1, y1-2)
	//draw.Draw(m, cbl, image.Transparent, image.Point{}, draw.Src)
	//cbr := image.Rect(x1-1, y0+2, x1, y1-2)
	//draw.Draw(m, cbr, image.Transparent, image.Point{}, draw.Src)
	//cbl2 := image.Rect(x0+1, y0+1, x0+2, y1-1)
	//draw.Draw(m, cbl2, image.Transparent, image.Point{}, draw.Src)
	//cbr2 := image.Rect(x1-2, y0+1, x1-1, y1-1)
	//draw.Draw(m, cbr2, image.Transparent, image.Point{}, draw.Src)

	// logo 圆角在 FitQrInnerIcon 里面设置的
	draw.Draw(m, icon.Bounds().Add(offset), icon, image.Point{}, draw.Over)

	buf := new(bytes.Buffer)
	if err = png.Encode(buf, m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
