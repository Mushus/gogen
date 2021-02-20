package main_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	main "github.com/Mushus/gogen/implgen"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

const recvType = `
package main

func (r Hoge) hello() {}
`

func TestRecvTypeIs(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", recvType, 0)
	if err != nil {
		t.Fatal(err)
	}

	astutil.Apply(f, func(c *astutil.Cursor) bool {
		fd, ok := c.Node().(*ast.FuncDecl)
		if ok {
			if main.GetRecvTypeName(fd) != "Hoge" {
				t.Fatal("GetRecvTypeName must be Hoge")
			}
		}
		return true
	}, nil)
}

func TestCollectImports(t *testing.T) {
	exported := packagestest.Export(t, packagestest.Modules, []packagestest.Module{{
		Name: "golang.org/fake",
		Files: map[string]interface{}{
			"a/a.go": `package a; import c "golang.org/fake/b"; type Z interface { MethodZ() }; type A interface { Z; MethodA(); c.B }`,
			"b/b.go": `package b; type B interface { MethodB() }`,
		},
	}})
	exported.Config.Mode = packages.LoadSyntax
	pkgs, err := packages.Load(exported.Config, "golang.org/fake/a")
	if err != nil {
		t.Fatal(err)
	}

	pkg := main.FindPkg(pkgs, "golang.org/fake/a")
	if pkg == nil {
		t.Fatal("undefined package")
	}

	p, err := main.CollectImports(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if p["c"] == nil {
		t.Fatal("p[\"c\"] is not empty")
	}
}

func TestCollectIfMethods(t *testing.T) {
	exported := packagestest.Export(t, packagestest.Modules, []packagestest.Module{{
		Name: "golang.org/fake",
		Files: map[string]interface{}{
			"a/a.go": `package a; import c "golang.org/fake/b"; type Z interface { MethodZ() }; type A interface { Z; MethodA(); c.B }`,
			"b/b.go": `package b; type B interface { MethodB() }`,
		},
	}})
	exported.Config.Mode = packages.LoadAllSyntax
	pkgs, err := packages.Load(exported.Config, "golang.org/fake/a")
	if err != nil {
		t.Fatal(err)
	}

	methods, err := main.CollectInterfaceMethods(pkgs, "golang.org/fake/a", "A")
	t.Logf("--> %#v\n", methods)
}
