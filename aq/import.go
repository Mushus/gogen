package aq

import (
	"go/ast"
)

type Import struct {
	importSpec *ast.ImportSpec
}

func createImport(importSpec *ast.ImportSpec) *Import {
	return &Import{
		importSpec: importSpec,
	}
}

func (i *Import) Exists() bool {
	return i != nil && i.importSpec != nil
}

func (i *Import) Name() string {
	if !i.Exists() {
		return ""
	}

	return safeIdentName(i.importSpec.Name)
}

func (i *Import) Path() string {
	if !i.Exists() {
		return ""
	}

	return stringLiteral(i.importSpec.Path)
}

type Imports []*Import
