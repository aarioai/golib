package ossz

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/kataras/iris/v12"
	"strconv"
)

type Video struct {
	Checksum string `json:"checksum"` // md5(content)
	Size     int    `json:"size"`     // in bytes
	Width    int
	Height   int
	Duration int            `json:"duration"` // 时长，秒
	FileType aenum.FileType `json:"filetype"`
}

func LoadVideoForm(ictx iris.Context, maxSize int64) (ReadCloser, Video, *ae.Error) {
	//maxSize := int64(100 << 20) // 100M
	err := ictx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		return nil, Video{}, ae.New(400, err.Error())
	}
	form := ictx.Request().MultipartForm
	mimes := form.Value["mime"]
	widths := form.Value["width"]
	heights := form.Value["height"]
	durations := form.Value["duration"]
	files := form.File["file"]
	if len(files) != 1 || len(widths) != 1 || len(heights) != 1 || len(mimes) != 1 || len(durations) != 1 {
		return nil, Video{}, ae.New(400, "file param `file` `width` `height` `filetype` `duration` should be and only be once")
	}
	fileHeader := files[0]
	mime := mimes[0]
	width, _ := strconv.ParseUint(widths[0], 10, 16)
	height, _ := strconv.ParseUint(heights[0], 10, 16)
	duration, _ := strconv.ParseUint(durations[0], 10, 16)
	if width == 0 || height == 0 || duration == 0 {
		return nil, Video{}, ae.NewBadParam("width, height, duration")
	}

	filetype, ok := aenum.NewVideoType(mime)
	if !ok {
		return nil, Video{}, ae.NewBadParam("mime")
	}

	r, err := fileHeader.Open()
	if err != nil {
		return nil, Video{}, ae.NewError(err)
	}
	sz, err := r.Seek(0, 2)
	if err != nil {
		r.Close()
		return nil, Video{}, ae.NewError(err)
	}

	r.Seek(0, 0)
	checksum, err := ReaderMd5(r)
	if err != nil {
		r.Close()
		return nil, Video{}, ae.NewError(err)
	}

	t := Video{
		Checksum: checksum,
		Size:     int(sz),
		Width:    int(width),
		Height:   int(height),
		Duration: int(duration),
		FileType: filetype,
	}
	r.Seek(0, 0)
	return r, t, nil

}
