// chris 080615

package main

import (
	"os"
	"path"
	"testing"
)

func testProcess(t *testing.T, srcname, dstname string) {
	src, err := os.Open(srcname)
	if err != nil {
		t.Error(err)
	}
	defer src.Close()
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	dst, err := os.OpenFile(dstname, flag, 0666)
	if err != nil {
		t.Error(err)
	}
	defer dst.Close()
	err = process(dst, srcname, src)
	if err != nil {
		t.Error(err)
	}
}

func TestProcess(t *testing.T) {
	srcname := path.Join("test", "comp.v")
	dstname := path.Join("test", "comp.go")
	testProcess(t, srcname, dstname)
	defer os.Remove(dstname)
}
