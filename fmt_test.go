// chris 080615

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"os/exec"
)

type closeWrapper struct {
	*bytes.Buffer
}

func (cw closeWrapper) Close() error {
	return nil
}

// Make sure that generated code isn't altered by running gofmt.
func testFmt(t *testing.T, dstname, srcname string) {
	var err error
	var fmtout io.ReadCloser
	var dst2 *os.File
	var output []byte

	testProcess(t, dstname, nil, srcname, nil)
	defer os.Remove(dstname)

	cmd := exec.Command("gofmt", dstname)

	fmtout, err = cmd.StdoutPipe()
	if err != nil {
		t.Error(err)
	}

	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	dstname2 := fmt.Sprintf("%s.fmt", dstname)
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	dst2, err = os.OpenFile(dstname2, flag, 0666)
	if err != nil {
		t.Error(err)
	}
	defer dst2.Close()
	defer os.Remove(dstname2)

	io.Copy(dst2, fmtout)

	if err := cmd.Wait(); err != nil {
		t.Error(err)
	}

	cmd = exec.Command("diff", "-u", dstname, dstname2)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("gofmt output differs; %v\n", err)
		t.Logf("%s\n", output)
	}
}

func TestFmt(t *testing.T) {
	srcname := path.Join("test", "comp.v")
	dstname := path.Join("test", "comp.go")
	testFmt(t, dstname, srcname)
}
