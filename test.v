package foo

type fooTypeSet struct {
	s string

	bo bool

	f32 float32
	f64 float64

	i8  int8
	i16 int16
	i32 int32
	i64 int64

	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64
}

type blahTypeSet struct {
	r rune // well-defined?
	by byte // well-defined?
}

type PublicStruct struct {
	x int
}
