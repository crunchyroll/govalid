// chris 080615

package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"os/exec"
)

// testBuild handles cleaning up the destination file.
func testBuild(t *testing.T, dstname, srcname string, src io.Reader) ([]byte, error) {
	testProcess(t, dstname, nil, srcname, src)
	defer os.Remove(dstname)

	// XXX Isn't there a way to build without having to launch a
	// subprocess?
	cmd := exec.Command("go", "build", dstname)
	return cmd.CombinedOutput()
}

// Make sure that generated validation code from test .v file builds
// without error.
func testBuildGood(t *testing.T, dstname, srcname string) {
	output, err := testBuild(t, dstname, srcname, nil)
	if err != nil {
		t.Errorf("build failed: %v, %s", err, output)
	}
}

// Make sure that generated validation code from good test .v files
// builds without error.
func TestBuildGood(t *testing.T) {
	for _, name := range goodnames {
		dstname := path.Join("test", fmt.Sprintf("%s.go", name))
		srcname := path.Join("test", fmt.Sprintf("%s.v", name))
		testBuildGood(t, dstname, srcname)
	}
}

// Make sure that generated validation code from bad test .v file builds
// with error.
func TestBuildBad(t *testing.T) {
	srcname := path.Join("test", "bad.v")
	dstname := path.Join("test", "bad.go")
	output, err := testBuild(t, dstname, srcname, nil)
	if err == nil {
		t.Errorf("build failed to fail: %s", output)
	}
}
