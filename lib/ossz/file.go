package ossz

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func ReaderMd5(f ReadCloser) (string, error) {
	f.Seek(0, 0) // 这里是第一个取，但是为了保险起见，还是把光标移到最开始
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	checksum := hex.EncodeToString(h.Sum(nil))
	return checksum, nil
}
