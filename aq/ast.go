package aq

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
)

func safeIdentName(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}

	return ident.Name
}

func safeIdentsName(itents []*ast.Ident) string {
	for _, ident := range itents {
		name := safeIdentName(ident)
		if name != "" {
			return name
		}
	}
	return ""
}

func stringLiteral(basicLit *ast.BasicLit) string {
	if basicLit == nil {
		return ""
	}

	tv, _ := types.Eval(token.NewFileSet(), nil, token.NoPos, basicLit.Value)
	return constant.StringVal(tv.Value)
}
