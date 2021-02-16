package aq

import (
	"go/ast"
)

type TypeSpec struct {
	instance *AQ           `getter:"-"`
	file     *File         `getter:"-"`
	spec     *ast.TypeSpec `getter:"-"`
}

func (t *TypeSpec) Interface() *Interface {
	i, ok := t.spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil
	}

	return NewInterface(t.instance, i)
}

func (t *TypeSpec) Struct() *Struct {
	s, ok := t.spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	return NewStruct(t.instance, t.file, t.spec, s)
}

func (t TypeSpecs) Interfaces() Interfaces {
	if t == nil {
		return nil
	}

	l := Interfaces{}
	for _, td := range t {
		i := td.Interface()
		if i != nil {
			l = append(l, i)
		}
	}
	return l
}

func (t TypeSpecs) Structs() Structs {
	if t == nil {
		return nil
	}

	l := Structs{}
	for _, td := range t {
		i := td.Struct()
		if i != nil {
			l = append(l, i)
		}
	}
	return l
}
