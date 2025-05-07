package ossz

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"
	"os"
)

// LoadImageBufReader
// @warn 不要忘了使用完后执行  f.Close()
func LoadImageBufReader(buf []byte) (ReadCloser, ImageConfig, error) {
	hasher := md5.New()
	hasher.Write(buf)
	checksum := hex.EncodeToString(hasher.Sum(nil))
	br := bytes.NewReader(buf)
	cf, err := DecodeImageConfig(br, checksum, len(buf))
	if err != nil {
		return nil, cf, err
	}
	r := NewBytesReader(br)
	return r, cf, nil
}

// DownloadImageReader
// @warn 不要忘了使用完后执行  f.Close()
func DownloadImageReader(ctx context.Context, url string) (ReadCloser, ImageConfig, error) {
	buf, err := Download(ctx, url)
	if err != nil {
		return nil, ImageConfig{}, err
	}
	return LoadImageBufReader(buf)
}

// LoadImageReader
// @warn 不要忘了使用完后执行  f.Close()
func LoadImageReader(src string) (ReadCloser, ImageConfig, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, ImageConfig{}, err
	}
	cf, err := DecodeImageFileConfig(f)
	return f, cf, err
}

func LoadImageForm(fileHeader *multipart.FileHeader) (ReadCloser, ImageConfig, error) {
	// 从数据库判断是否存在，不存在就保存
	f, err := fileHeader.Open()
	if err != nil {
		return nil, ImageConfig{}, err
	}
	cf, err := DecodeImageFileConfig(f)
	return f, cf, err
}
