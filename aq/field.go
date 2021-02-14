package aq

import (
	"go/ast"
)

type Field struct {
	aqStruct *Struct
	field    *ast.Field
}

func createField(aqStruct *Struct, field *ast.Field) *Field {
	return &Field{
		aqStruct: aqStruct,
		field:    field,
	}
}

func (f *Field) Exists() bool {
	return f != nil
}

func (f *Field) Name() string {
	return safeIdentsName(f.field.Names)
}

func (f *Field) Type() *Type {
	return createType(f.field.Type)
}

func (f *Field) Tag() *Tag {
	return createTag(f.field.Tag)
}

func (f *Field) IsPublic() bool {
	if !f.Exists() {
		return false
	}

	name := f.Name()
	if name == "" {
		return false
	}
	// TODO: embed

	nameRune := []rune(name)
	firstLetter := nameRune[0]
	return 'A' <= firstLetter && firstLetter <= 'Z'
}

type Fields []*Field

func (f Fields) Exists() bool {
	return f != nil || f.Count() == 0
}

func (f Fields) Count() int {
	return len(f)
}

func (f Fields) FindByName(name string) *Field {
	for _, field := range f {
		if field.Name() == name {
			return field
		}
	}
	return nil
}
