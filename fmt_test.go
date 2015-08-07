// chris 080615

package main

import (
	"bytes"
	"io"
	"path"
	"strings"
	"testing"

	"os/exec"

	"github.com/aryann/difflib"
)

// Make sure that generated code isn't altered by running gofmt.
func testFmt(t *testing.T, srcname string, src io.Reader) {
	var err error

	var fmtin    io.WriteCloser
	var fmtout   io.ReadCloser
	var validbuf *bytes.Buffer
	var fmtbuf   *bytes.Buffer

	var validbytes []byte

	// process (i.e., "govalid" output) will be buffered here.
	validbuf = new(bytes.Buffer)
	// gofmt output will be buffered here.
	fmtbuf = new(bytes.Buffer)

	testProcess(t, "-", validbuf, srcname, src)
	validbytes = validbuf.Bytes()

	// XXX Isn't there a way to gofmt without having to launch a
	// subprocess?
	cmd := exec.Command("gofmt")

	fmtin, err = cmd.StdinPipe()
	if err != nil {
		t.Error(err)
	}
	fmtout, err = cmd.StdoutPipe()
	if err != nil {
		t.Error(err)
	}

	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	sync := make(chan struct{})
	go func() {
		io.Copy(fmtin, bytes.NewBuffer(validbytes))
		fmtin.Close()
		sync <- struct{}{}
	}()
	go func() {
		io.Copy(fmtbuf, fmtout)
		sync <- struct{}{}
	}()
	<-sync
	<-sync

	if err := cmd.Wait(); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(validbytes, fmtbuf.Bytes()) {
		t.Log("gofmt output differs")
		validlines := strings.Split(string(validbytes), "\n")
		fmtlines := strings.Split(string(fmtbuf.Bytes()), "\n")
		diffrecords := difflib.Diff(validlines, fmtlines)
		t.Log("--- valid output")
		t.Log("+++ gofmt output")
		for _, dr := range diffrecords {
			if dr.Delta == difflib.Common {
				continue
			}
			t.Log(dr)
		}
		t.Fail()
	}
}

func TestFmt(t *testing.T) {
	srcname := path.Join("test", "comp.v")
	testFmt(t, srcname, nil)
}
