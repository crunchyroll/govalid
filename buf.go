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

func newBuf() *buf {
	return &buf{Buffer: new(bytes.Buffer), needsStrconv: false}
}

func (b *buf) writef(format string, a ...interface{}) {
	b.Buffer.WriteString(fmt.Sprintf(format, a...))
}
