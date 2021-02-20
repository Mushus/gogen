package aq

import (
	"go/ast"
)

type StructSpec struct {
	instance   *AQ             `getter:"-"`
	file       *File           `getter:"-"`
	typeSpec   *ast.TypeSpec   `getter:"-"`
	structType *ast.StructType `getter:"-"`
}

func createStruct(file *File, typeSpec *ast.TypeSpec, structType *ast.StructType) *StructSpec {
	return &StructSpec{
		file:       file,
		typeSpec:   typeSpec,
		structType: structType,
	}
}

func (s StructSpec) File() *File {
	return s.file
}

func (s *StructSpec) Exists() bool {
	return s != nil
}

func (s *StructSpec) Name() string {
	if !s.Exists() {
		return ""
	}

	return safeIdentName(s.typeSpec.Name)
}

func (s *StructSpec) Fields() Fields {
	if s.structType.Fields == nil {
		return nil
	}

	fields := make(Fields, 0, len(s.structType.Fields.List))
	for _, field := range s.structType.Fields.List {
		fields = append(fields, NewField(field))
	}

	return fields
}

type HasStructFunc func(i int, s *StructSpec) bool

func StructNameIs(name string) HasStructFunc {
	return func(i int, s *StructSpec) bool {
		return s.Name() == name
	}
}
