package aq

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type FuncDecl struct {
	decl *ast.FuncDecl `getter:"-"`
}

func (f *FuncDecl) Exists() bool {
	return f != nil
}

func (f *FuncDecl) Name() string {
	if f == nil {
		return ""
	}

	return safeIdentName(f.decl.Name)
}

func (f *FuncDecl) Recv() *Field {
	if f == nil {
		return nil
	}

	fl := f.decl.Recv
	if fl == nil {
		return nil
	}

	return NewFields(fl).First()
}

func (f *FuncDecl) Type() *FuncType {
	if f == nil {
		return nil
	}
	return NewFuncType(f.decl.Type)
}

func (f *FuncDecl) Params() Fields {
	return f.Type().Params()
}

func (f *FuncDecl) Results() Fields {
	return f.Type().Results()
}

func (f *FuncDecl) BodyCode() string {
	if f == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), f.decl.Body)
	return buf.String()
}
