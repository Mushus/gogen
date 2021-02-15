package aq

import (
	"go/ast"
)

type Field struct {
	field *ast.Field `getter:"-"`
}

func (f *Field) Exists() bool {
	return f != nil
}

func (f *Field) Name() string {
	if !f.Exists() {
		return ""
	}
	return safeIdentsName(f.field.Names)
}

func (f *Field) Type() *Type {
	if !f.Exists() {
		return nil
	}
	return createType(f.field.Type)
}

func (f *Field) Tag() *Tag {
	if !f.Exists() {
		return nil
	}
	return createTag(f.field.Tag)
}

func (f *Field) IsExported() bool {
	if !f.Exists() {
		return false
	}

	name := f.Name()
	if name == "" {
		return false
	}
	// TODO: embed

	return ast.IsExported(name)
}

func NewFields(f *ast.FieldList) Fields {
	if f == nil {
		return nil
	}
	l := make(Fields, 0, f.NumFields())
	for _, v := range f.List {
		l = append(l, NewField(v))
	}

	return l
}

func (f Fields) FindByName(name string) *Field {
	return f.Find(func(i int, v *Field) bool { return v.Name() == name })
}
