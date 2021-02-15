package aq_test

import (
	"reflect"
	"testing"

	"github.com/Mushus/gogen/aq"
)

func TestStructs(t *testing.T) {
	source := []byte(`
	package main

	type Foo struct {}

	type Bar Foo
	`)

	structs := aq.New().
		MustLoadFromSource(source).
		Structs()

	if len(structs) != 1 {
		t.Fail()
	}

	if structs[0].Name() != "Foo" {
		t.Fail()
	}
}

func TestStructsHas(t *testing.T) {
	source := []byte(`
	package main

	type Foo struct {}

	type Bar struct {}
	`)

	structs := aq.New().
		MustLoadFromSource(source).
		Structs()

	if !structs.Has(aq.StructNameIs("Foo")) {
		t.Fail()
	}
}

func TestStructsField(t *testing.T) {
	source := []byte(`
	package main

	type Foo struct {
		Bar string
		Baz int
	}
	`)

	f := aq.New().
		MustLoadFromSource(source).
		Structs().
		Find(aq.StructNameIs("Foo")).
		Fields()

	if f.Count() != 2 {
		t.Fatalf("expect 2 got %d", f.Count())
	}
}

func TestStructTag(t *testing.T) {
	source := []byte(`
	package main

	type Foo struct {
		Bar string ` + "`construct:\"-\"`" + `
	}
	`)

	tag := aq.New().
		MustLoadFromSource(source).
		Structs().
		Find(aq.StructNameIs("Foo")).
		Fields().
		FindByName("Bar").
		Tag().Body()
	except := reflect.StructTag("construct:\"-\"")

	if tag != except {
		t.Fatalf("expect %#v got %#v", except, tag)
	}
}

func TestImport(t *testing.T) {
	source := []byte(`
	package main

	import (
		test "go/ast"
	)
	`)

	importPath := aq.New().
		MustLoadFromSource(source).
		File().
		Imports()[0].Path()
	except := "go/ast"

	if importPath != except {
		t.Fatalf("expect %#v got %#v", except, importPath)
	}
}
