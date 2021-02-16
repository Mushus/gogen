package main

import (
	"go/parser"
	"go/token"

	"github.com/k0kubun/pp"
)

const d = `
package main

type (
	st string
	it interface{}
)

const sv = ""
var iv = 0

func fn() {}

i := 0
`

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", []byte(d), parser.ParseComments)

	pp.Println(f.Scope.Objects)
}
