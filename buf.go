// chris 071415

package main

import (
	"bytes"
	"fmt"
)

type buf struct {
	*bytes.Buffer
	needsStrconv bool
}

func newBuf(size int64) *buf {
	b := make([]byte, 0, size)
	bu := bytes.NewBuffer(b)
	return &buf{Buffer: bu, needsStrconv: false}
}

func (b *buf) writef(format string, a ...interface{}) {
	b.Buffer.WriteString(fmt.Sprintf(format, a...))
}
