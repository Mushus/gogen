package aq

import (
	"go/ast"
)

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
