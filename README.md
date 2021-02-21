# gogen

## implegen

continuous generation of implementations for interfaces

### Motivation

It is very hard to keep up with the implementation of a fluid and complex interface.
This tool will add and update methods as the interface changes.

### Get Started

By using implgen you can automatically generate a set of methods defined for a specific interface type.
The following example writes the method set defined for the interface type `InterfaceName` in the import path `path/to/a/interface` to the file `update_file.go`

``` bash
go run github.com/Mushus/gogen/implgen -dts "update_file.go" ImplTypeName "path/to/a/interface.InterfaceName"
```

If you use the go generate command to embed directives as comments, you do not need the `dst` flag.

```go
//go:generate go run github.com/Mushus/gogen/implgen ImplTypeName "path/to/interface.InterfaceName"
```

When auto-generating for the io.WriteCloser interface type, the bottom half of the following code is generated.

```go
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

```

### Install as command

```
go install github.com/Mushus/gogen/implgen
```
