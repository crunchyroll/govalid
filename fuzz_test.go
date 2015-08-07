// chris 080715

package main

import (
	"bytes"
	"os"
	"path"
	"testing"
	"time"

	"io/ioutil"
	"math/rand"

	"chrispennello.com/go/rebnf"

	"golang.org/x/exp/ebnf"
)

func testFuzz(t *testing.T, srcname, start string) {
	var err error
	var grammar ebnf.Grammar
	var tempfile *os.File

	grammar, err = rebnf.Parse(srcname, nil)
	if err != nil {
		t.Error(err)
	}
	dst := new(bytes.Buffer)
	ctx := rebnf.NewCtx(20, 20, " ", false)
	err = ctx.Random(dst, grammar, start)
	if err != nil {
		t.Error(err)
	}
	dstbytes := dst.Bytes()

	// Test gofmt.
	testFmt(t, "-", bytes.NewBuffer(dstbytes))

	// Test go build.
	tempfile, err = ioutil.TempFile("test", "tmp_")
	if err != nil {
		t.Error(err)
	}
	tempfile.Close()
	testBuild(t, tempfile.Name(), "-", bytes.NewBuffer(dstbytes))
}

func TestFuzz(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping fuzz tests in short mode")
	}

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	t.Logf("chose random seed %d\n", seed)

	srcname := path.Join("test", "struct.ebnf")
	for i := 0; i < 100; i++ {
		testFuzz(t, srcname, "Start")
	}
}
