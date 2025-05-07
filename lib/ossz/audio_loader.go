package ossz

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"mime/multipart"
)

type Audio struct {
	Checksum string         `json:"checksum"` // md5(content)
	Size     int            `json:"size"`     // in bytes
	Duration int            `json:"duration"` // 时长，秒
	FileType aenum.FileType `json:"filetype"`
}

func LoadAudioForm(fileHeader *multipart.FileHeader, filetype aenum.FileType, duration int) (ReadCloser, Audio, *ae.Error) {
	r, err := fileHeader.Open()
	if err != nil {
		return nil, Audio{}, ae.NewE("audio file header open: " + err.Error())
	}
	sz, err := r.Seek(0, 2)
	if err != nil {
		r.Close()
		return nil, Audio{}, ae.NewE("audio file header seek(0, 2): " + err.Error())
	}

	r.Seek(0, 0)
	checksum, err := ReaderMd5(r)
	if err != nil || checksum == "" {
		r.Close()
		return nil, Audio{}, ae.NewE("audio reader md5: %v", err)
	}

	t := Audio{
		Checksum: checksum,
		Size:     int(sz),
		Duration: int(duration),
		FileType: filetype,
	}
	r.Seek(0, 0)
	return r, t, nil

}
