package coding_test

import (
	"bytes"
	"github.com/aarioai/golib/lib/code/coding"
	"testing"
	"time"
)

func TestEncodeHex(t *testing.T) {
	text := []byte("~" + coding.RandAlphabets(10, time.Now().UnixMicro()) + "!")
	textClone := bytes.Clone(text)

	hexText := coding.EncodeHex(text)
	if !bytes.Equal(text, textClone) {
		t.Error("EncodeHex changed text slice")
		return
	}

	got, err := coding.DecodeHex(hexText)
	if err != nil {
		t.Errorf("DecodeHex failed %v", err)
		return
	}
	if !bytes.Equal(text, got) {
		t.Errorf("DecodeHex got %s, want %s", got, text)
		return
	}
	if !bytes.Equal(text, textClone) {
		t.Error("EncodeHex changed text slice")
		return
	}

}
