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
	expr := t.typ
	if star, ok := t.typ.(*ast.StarExpr); ok {
		expr = star.X
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return safeIdentName(ident)
	}
	return ""
}

func (t *Type) GoCode() string {
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), t.typ)
	return buf.String()
}

func (t *Type) IsPtr() bool {
	_, ok := t.typ.(*ast.StarExpr)
	return ok
}
