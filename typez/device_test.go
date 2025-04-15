package typez_test

import (
	"github.com/aarioai/golib/typez"
	"testing"
)

func TestDeviceInfo(t *testing.T) {
	_, err := typez.DecodeDeviceInfo("AVDMpRjYpBVd2ETYwQTJmZDOtUXdwJTOw1TNygjZ5IzN3YWNykTJkBnPuMWP0ASZ1dGM9JzczZGJk1XNh8zPyETYzhmdpJGZ4MDYtRDOiEzMycjY4gzO1YDNmJnbvdWM0UTMwYDcocDN4")
	if err != nil {
		t.Errorf(`decode device info: ` + err.Error())
	}

}
