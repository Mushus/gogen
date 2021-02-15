package aq

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type Type struct {
	typ ast.Expr
}

func createType(typ ast.Expr) *Type {
	return &Type{
		typ: typ,
	}
}

func (t *Type) Name() string {
	if t == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), t.typ)
	return buf.String()
}

func (t *Type) GoCode() string {
	if t == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), t.typ)
	return buf.String()
}

func (t *Type) IsPtr() bool {
	if t == nil {
		return false
	}
	_, ok := t.typ.(*ast.StarExpr)
	return ok
}

func (t *Type) UnwrapPtr() *Type {
	if t == nil {
		return nil
	}
	if star, ok := t.typ.(*ast.StarExpr); ok {
		return createType(star.X)
	}
	return t
}
