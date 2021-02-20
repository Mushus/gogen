package aq

import (
	"go/ast"
)

type TypeSpec struct {
	aq   *AQ           `getter:"-"`
	file *File         `getter:"-"`
	spec *ast.TypeSpec `getter:"-"`
}

func (t *TypeSpec) Interface() *InterfaceSpec {
	i, ok := t.spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil
	}

	return NewInterfaceSpec(t.aq, t.file, t.spec, i)
}

func (t *TypeSpec) Struct() *StructSpec {
	s, ok := t.spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	return NewStructSpec(t.aq, t.file, t.spec, s)
}

func (t TypeSpecs) Interfaces() InterfaceSpecs {
	if t == nil {
		return nil
	}

	l := InterfaceSpecs{}
	for _, td := range t {
		i := td.Interface()
		if i != nil {
			l = append(l, i)
		}
	}
	return l
}

func (t TypeSpecs) Structs() StructSpecs {
	if t == nil {
		return nil
	}

	l := StructSpecs{}
	for _, td := range t {
		i := td.Struct()
		if i != nil {
			l = append(l, i)
		}
	}
	return l
}
