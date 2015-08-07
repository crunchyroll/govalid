// chris 080615

package main

import (
	"io"
	"os"
	"path"
	"testing"
)

func testProcess(t *testing.T, dstname string, dst io.WriteCloser, srcname string, src io.ReadCloser) {
	var err error

	if src == nil {
		src, err = os.Open(srcname)
		if err != nil {
			t.Error(err)
		}
		defer src.Close()
	}

	if dst == nil {
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		dst, err = os.OpenFile(dstname, flag, 0666)
		if err != nil {
			t.Error(err)
		}
		defer dst.Close()
	}

	err = process(dst, srcname, src)
	if err != nil {
		t.Error(err)
	}
}

func TestProcess(t *testing.T) {
	srcname := path.Join("test", "comp.v")
	dstname := path.Join("test", "comp.go")
	testProcess(t, dstname, nil, srcname, nil)
	defer os.Remove(dstname)
}
