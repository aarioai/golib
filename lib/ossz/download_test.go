package ossz_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/aarioai/golib/lib/ossz"
	"io"
	"testing"
)

func TestDownload(t *testing.T) {
	url := "http://www.law51.net/images/tu34.jpg"
	ctx := context.Background()
	buf, _ := ossz.Download(ctx, url)

	hasher := md5.New()
	hasher.Write(buf)
	checksum1 := hex.EncodeToString(hasher.Sum(nil))
	t.Log(len(buf), "===>", checksum1)
	f, err := ossz.DownloadTmp(ctx, url)
	if err != nil {
		t.Log(err.Text())
		return
	}
	defer f.Close()

	f.Seek(0, 0) // 这里是第一个取，但是为了保险起见，还是把光标移到最开始
	h := md5.New()
	io.Copy(h, f)
	checksum2 := hex.EncodeToString(h.Sum(nil))

	f.Seek(0, 0)
	fi, _ := f.Stat()
	t.Log(fi.Size(), "===>", checksum2)
}
