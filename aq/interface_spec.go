package aq

import (
	"go/ast"
)

type InterfaceSpec struct {
	aq   *AQ                `getter:"-"`
	file *File              `getter:"-"`
	spec *ast.TypeSpec      `getter:"-"`
	typ  *ast.InterfaceType `getter:"-"`
}

func (i *InterfaceSpec) Name() string {
	if i == nil {
		return ""
	}
	return safeIdentName(i.spec.Name)
}

func (i *InterfaceSpec) Methods() Fields {
	l := make(Fields, 0, len(i.typ.Methods.List))
	for _, m := range i.typ.Methods.List {
		l = append(l, NewField(m))
	}
	return l
}
