package aq

import (
	"go/ast"
)

type Func struct {
	decl *ast.FuncDecl `getter:"-"`
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
	if f == nil {
		return ""
	}

	return safeIdentName(f.decl.Name)
}

func (f *Func) Recv() *Field {
	if f == nil {
		return nil
	}

	fl := f.decl.Recv
	if fl == nil {
		return nil
	}

	return NewFields(fl).First()
}

func (f *Func) Type() *FuncType {
	if f == nil {
		return nil
	}
	return NewFuncType(f.decl.Type)
}

func (f *Func) Params() Fields {
	return f.Type().Params()
}

func (f *Func) Results() Fields {
	return f.Type().Results()
}

type FuncType struct {
	typ *ast.FuncType
}

func (t *FuncType) Params() Fields {
	if t == nil || t.typ == nil || t.typ.Params == nil {
		return nil
	}
	return NewFields(t.typ.Params)
}

func (t *FuncType) Results() Fields {
	if t == nil || t.typ == nil || t.typ.Results == nil {
		return nil
	}
	return NewFields(t.typ.Results)
}
