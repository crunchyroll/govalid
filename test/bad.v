package badpkg

type foo struct {
	x int
}

type badTypeSet struct {
	bad1 foo
	bad2 *foo
	bad3 []foo
	bad4 map[string]foo
}

type badEmpty struct {
}
