// chris 080615

package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

// The caller is responsible for cleaning up any left over files.  If
// nil source or destination writers are passed in, thus instructing
// testProcess to open files for them based on the given names, then
// testProcess will close the opened files.
func testProcess(t *testing.T, dstname string, dst io.Writer, srcname string, src io.Reader) {
	if src == nil {
		srcfile, err := os.Open(srcname)
		if err != nil {
			t.Error(err)
		}
		defer srcfile.Close()
		src = srcfile
	}

	if dst == nil {
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		dstfile, err := os.OpenFile(dstname, flag, 0666)
		if err != nil {
			t.Error(err)
		}
		defer dstfile.Close()
		dst = dstfile
	}

	if err := process(dst, srcname, src); err != nil {
		t.Error(err)
	}
}

func TestProcess(t *testing.T) {
	for _, name := range testGoodNames {
		dstname := path.Join("test", fmt.Sprintf("%s.go", name))
		srcname := path.Join("test", fmt.Sprintf("%s.v", name))
		testProcess(t, dstname, nil, srcname, nil)
		defer os.Remove(dstname)
	}
}
