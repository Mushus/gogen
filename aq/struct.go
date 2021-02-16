package aq

import (
	"go/ast"
)

type Struct struct {
	instance   *AQ
	file       *File
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
}

func createStruct(file *File, typeSpec *ast.TypeSpec, structType *ast.StructType) *Struct {
	return &Struct{
		file:       file,
		typeSpec:   typeSpec,
		structType: structType,
	}
}

func (s Struct) File() *File {
	return s.file
}

func (s *Struct) Exists() bool {
	return s != nil
}

func (s *Struct) Name() string {
	if !s.Exists() {
		return ""
	}

	return safeIdentName(s.typeSpec.Name)
}

func (s *Struct) Fields() Fields {
	if s.structType.Fields == nil {
		return nil
	}

	fields := make(Fields, 0, len(s.structType.Fields.List))
	for _, field := range s.structType.Fields.List {
		fields = append(fields, NewField(field))
	}

	return fields
}

type HasStructFunc func(i int, s *Struct) bool

func StructNameIs(name string) HasStructFunc {
	return func(i int, s *Struct) bool {
		return s.Name() == name
	}
}
