package test

// test struct
type testStruct struct {
	typ   string
	hoge  int `getter:"-"`
	Width int
	name  string `setter:"Rename"`
	star  *int
}

func (t *testStruct) construct() {
}
