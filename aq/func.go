package aq

import (
	"go/ast"
)

//go:generate go run github.com/Mushus/gogen/structgen -list Func
type Func struct {
	decl *ast.FuncDecl
}

func createFunc(decl *ast.FuncDecl) *Func {
	if decl == nil {
		return nil
	}

	return &Func{
		decl: decl,
	}
}

func (f *Func) Exists() bool {
	return f != nil
}

func (f *Func) Name() string {
	if !f.Exists() {
		return ""
	}

	return safeIdentName(f.decl.Name)
}

func (f *Func) Recv() *Field {
	if !f.Exists() {
		return nil
	}

	fl := f.decl.Recv
	if fl == nil {
		return nil
	}

	for _, recv := range fl.List {
		return createField(nil, recv)
	}
	return nil
}
