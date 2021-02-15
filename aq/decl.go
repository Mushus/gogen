package aq

import "go/ast"

type Decl struct {
	decl ast.Decl `getter:"-"`
}

func (d *Decl) Types() TypeDefs {
	gen, ok := d.decl.(*ast.GenDecl)
	if !ok {
		return nil
	}

	l := TypeDefs{}
	for _, spec := range gen.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		l = append(l, NewTypeDef(typeSpec))
	}

	return l
}
