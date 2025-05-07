package ossz

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"io"
	"net/http"
	"os"
)

// Download 不要返回 http.Response，那个很不方便，容易因为指针被清空；而且使用没有关闭，也会导致内存泄露
// bytes.NewReader([]byte) 就可以转换回 io>ReadCloser
// resp.Body =  ioutil.NopCloser(bytes.NewBuffer(buf))  可以返回头设置回被 io.ReadAll 清空的
func Download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s, status code: %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func DownloadTmp(ctx context.Context, url string) (*os.File, *ae.Error) {
	buf, err := Download(ctx, url)
	if err != nil {
		return nil, ae.NewError(err)
	}

	var out *os.File
	out, err = os.CreateTemp("", "down")
	if err != nil {
		return nil, ae.NewError(err).WithDetail("create temp file")
	}
	// Write the body to file
	if _, err = out.Write(buf); err != nil {
		out.Close()
		os.Remove(out.Name())
		return nil, ae.NewError(err).WithDetail("write temp file")
	}
	//out.Seek(0, 0)
	return out, nil
}
