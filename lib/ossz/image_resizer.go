package ossz

import (
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

// ResizeImageReaderTo
// @warn 不要尝试在图片中加像素点的方式，来加密图片或者判断图片是否修改过。因为目前微信或其他软件传图片时候，都会压缩图片，纯色都会压缩为相近色
func ResizeImageReaderTo(r ReadCloser, cf ImageConfig, width, height int, anchor imaging.Anchor, dstFn func(ImageConfig) (string, error)) (ImageConfig, bool, error) {
	// 这里无法探知图片的md5，filename 就是动态的，不安全的， 不可用FileExists()判断文件是否存在！！！
	r.Seek(0, 0)

	// imaging.AutoOrientation 自动识别EXIF图片方向，并旋转
	im, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return cf, false, fmt.Errorf("imaging.Decode: " + err.Error())
	}

	//
	//switch im.(type) {
	//case *image.Paletted:
	//	cf.FileType = aenum.Gif // gif 动图
	//default:
	//	cf.FileType = aenum.Jpeg
	//}

	dst, err := dstFn(cf)
	if err != nil {
		return cf, false, err
	}

	if cf.Width != width || cf.Height != height {
		// dst 文件名是包括 md5的，此时md5已经变了，所以不能直接存dst
		tmp := os.TempDir() + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + path.Base(dst) + cf.FileType.Ext()
		defer os.Remove(tmp)
		// imaging.Lanczos 暂用空间太大
		// imaging.CropPattern = imaging.Resize + imaging.Crop

		dstImgFill := imaging.Fill(im, int(width), int(height), anchor, imaging.NearestNeighbor)
		if err = imaging.Save(dstImgFill, tmp, imaging.JPEGQuality(80)); err != nil {
			return cf, false, fmt.Errorf("imaging.Save to %s %s ", tmp, err.Error())
		}
		if im, err = imaging.Open(tmp); err != nil {
			return cf, false, fmt.Errorf("imaging.Decode saved file: " + err.Error())
		}

		var newFile *os.File
		newFile, cf, err = DecodeImageFile(tmp)
		if err != nil {
			return cf, false, err
		}

		dst, err = dstFn(cf)
		if err != nil {
			return cf, false, err
		}
		w, err := os.Create(dst)
		if err != nil {
			return cf, false, err
		}
		defer w.Close()
		newFile.Seek(0, 0)
		_, err = io.Copy(w, newFile)
		newFile.Close()
		return cf, false, err
	}

	// dst 跟md5 size width height 等绑定，如果一致，则可以直接返回
	if _, err = os.Stat(dst); err == nil {
		return cf, true, nil
	}
	// 这里图片md5、尺寸等完全一致，只有尺寸一致，但是本地没有存储，才可以直接返回cf（因为是直接copy的）

	w, err := os.Create(dst)
	if err != nil {
		return cf, false, err
	}
	defer w.Close()
	r.Seek(0, 0) // 指针放到头部
	_, err = io.Copy(w, r)
	return cf, false, err
}
