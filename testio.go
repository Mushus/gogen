package main

//go:generate go run github.com/Mushus/gogen/implgen TestWriter "io.WriteCloser"
type TestWriter struct{}

// Auto generated code below here.


func (t TestWriter) Close() error {
	panic("not implemented")
}

func (t TestWriter) Write(p []byte) (n int, err error) {
	panic("not implemented")
}

