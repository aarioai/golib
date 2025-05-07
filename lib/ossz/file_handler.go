package ossz

import (
	"os"
)

func ForceSave(filepath string, buf []byte) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(buf); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}
