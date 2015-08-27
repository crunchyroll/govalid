package whatever

type sliceBytes struct {
	x []byte `valid:"max:64,def:[]byte(\"abc\")"`
	y int
	z []int `valid:"min:2"`
}
