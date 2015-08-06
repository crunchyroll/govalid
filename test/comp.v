package foo

// Test import here to see what happens if the generated code needs more
// imports.
import "log"
import (
	"fmt"
	"strings"
)

func init() {
	log.Println("log use")
	fmt.Println("fmt use")
	fmt.Println(strings.ToUpper("strings use"))
}

type fooTypeSet struct {
	s string

	bo bool

	f32 float32
	f64 float64

	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64

	u   uint
	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64

  // No rune or byte.
}

type PublicStruct struct {
	x int
}

type mixed struct {
	x int
	bad1 fooTypeSet
}
