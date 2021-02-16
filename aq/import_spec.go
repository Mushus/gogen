package aq

import (
	"go/ast"
)

type ImportSpec struct {
	importSpec *ast.ImportSpec `getter:"-"`
}

func (i *ImportSpec) Exists() bool {
	return i != nil && i.importSpec != nil
}

func (i *ImportSpec) Name() string {
	if !i.Exists() {
		return ""
	}

	return safeIdentName(i.importSpec.Name)
}

func (i *ImportSpec) Path() string {
	if !i.Exists() {
		return ""
	}

	return stringLiteral(i.importSpec.Path)
}
