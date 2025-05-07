package ossz

import (
	"bytes"
	"errors"
	"github.com/aarioai/airis/aa/aenum"
	"image"
	"image/color"
	"io"
	"os"
)

type ImageConfig struct {
	ColorModel color.Model
	Checksum   string
	Size       int
	Width      int
	Height     int
	FileType   aenum.FileType
}

func DecodeImageConfig(f io.Reader, checksum string, size int) (ImageConfig, error) {
	c, formatName, err := image.DecodeConfig(f)
	if err != nil {
		return ImageConfig{}, err
	}
	fileType, ok := aenum.NewImageType("." + formatName)
	if !ok {
		return ImageConfig{}, errors.New("unknown image format: " + formatName)
	}
	cf := ImageConfig{
		ColorModel: c.ColorModel,
		Checksum:   checksum,
		Size:       size,
		Width:      c.Width,
		Height:     c.Height,
		FileType:   fileType,
	}

	return cf, nil
}

func DecodeImageFileConfig(f ReadCloser) (ImageConfig, error) {
	checksum, err := ReaderMd5(f)
	if err != nil {
		return ImageConfig{}, err
	}
	f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	// 获取图片宽高
	// @note 这里需要在bootstrap注册  image.RegisterFormat  支持的格式
	//
	f.Seek(0, 0) // 必须要把光标移到最开始
	return DecodeImageConfig(f, checksum, buf.Len())
}

func DecodeImageFile(p string) (*os.File, ImageConfig, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, ImageConfig{}, err
	}

	cf, err := DecodeImageFileConfig(file)
	if err != nil {
		file.Close()
		return nil, ImageConfig{}, err
	}

	return file, cf, nil
}
