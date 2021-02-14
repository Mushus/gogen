package aq

import (
	"go/ast"
)

type Struct struct {
	instance   *Instance
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
		fields = append(fields, createField(s, field))
	}

	return fields
}

type HasStructFunc func(s *Struct) bool

func StructNameIs(name string) HasStructFunc {
	return func(s *Struct) bool {
		return s.Name() == name
	}
}

func (s *Struct) Has(funcs ...HasStructFunc) bool {
	if !s.Exists() {
		return false
	}

	for _, f := range funcs {
		if !f(s) {
			return false
		}
	}
	return true
}

type Structs []*Struct

func (s Structs) Exists() bool {
	return s != nil || s.Count() == 0
}

func (s Structs) Find(funcs ...HasStructFunc) Structs {
	sl := make(Structs, 0)

	for _, s := range s {
		if s.Has(funcs...) {
			sl = append(sl, s)
		}
	}

	return sl
}

func (s Structs) FindOne(funcs ...HasStructFunc) *Struct {
	for _, s := range s {
		if s.Has(funcs...) {
			return s
		}
	}

	return nil
}

func (s Structs) Has(funcs ...HasStructFunc) bool {
	for _, s := range s {
		if s.Has(funcs...) {
			return true
		}
	}
	return false
}

func (s Structs) Count() int {
	return len(s)
}

func (s Structs) Get(index int) *Struct {
	if !s.Exists() {
		return nil
	}

	if index < 0 || s.Count() <= index {
		return nil
	}

	return s[index]
}
