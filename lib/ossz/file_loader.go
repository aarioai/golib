package ossz

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/kataras/iris/v12"
	"path/filepath"
)

type File struct {
	Checksum string `json:"checksum"` // md5(content)
	Size     int    `json:"size"`     // in bytes
	Ext      string `json:"ext"`
}

func LoadFileForm(ictx iris.Context, maxSize int64) (ReadCloser, File, *ae.Error) {
	//maxSize := int64(10 << 20) // 10M
	err := ictx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		return nil, File{}, ae.New(400, err.Error())
	}
	form := ictx.Request().MultipartForm
	files := form.File["file"]
	if len(files) != 1 {
		return nil, File{}, ae.New(400, "param `file` should be and only be once")
	}
	fileHeader := files[0]

	f, err := fileHeader.Open()
	if err != nil {
		return nil, File{}, ae.NewError(err)
	}
	sz, err := f.Seek(0, 2)
	if err != nil {
		f.Close()
		return nil, File{}, ae.NewError(err)
	}

	f.Seek(0, 0)
	checksum, err := ReaderMd5(f)
	if err != nil {
		f.Close()
		return nil, File{}, ae.NewError(err)
	}

	fi := File{
		Checksum: checksum,
		Size:     int(sz),
		Ext:      filepath.Ext(fileHeader.Filename),
	}
	return f, fi, nil

	//不要使用下面 ictx.UploadFormFiles，莫名其妙会出bug，而且还不好调试，使用原生方法
	// ch := make(chan bool)
	// ictx.UploadFormFiles(dir, func(ictx iris.Context, h *multipart.FileHeader) {
	// 	done := false

	// 	defer func() {
	// 		if !done {
	// 			h.Filename = "/tmp/iris.img.tmp.deletable"
	// 		}
	// 		go func(ch chan bool) {
	// 			ch <- true
	// 		}(ch)
	// 	}()

	// 	var (
	// 		file multipart.File
	// 		buf  []byte
	// 		ext  string
	// 	)
	// 	size = h.Size
	// 	file, err = h.Open()
	// 	defer file.Close()
	// 	if err != nil {
	// 		return
	// 	}

	// 	if buf, err = ioutil.ReadAll(file); err != nil {
	// 		return
	// 	}
	// 	md5b := md5.Sum(buf)
	// 	md5s = hex.EncodeToString(md5b[:])

	// 	ext, err = oas.ParseImgExt(buf)
	// 	if err != nil {
	// 		ext = ".jpg"
	// 	}
	// 	fullPath, shortPath, err = s.genFilePath(ictx.Request().Context(), workpath, ext)
	// 	if err != nil {
	// 		return
	// 	}

	// 	os.MkdirAll(path.Dir(fullPath), os.ModePerm)
	// 	h.Filename = shortPath
	// 	done = true
	// })

	// select {
	// case <-ch:
	// case <-time.After(time.Second * 10):
	// }
}
