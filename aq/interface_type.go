package aq

import "go/ast"

type InterfaceType struct {
	i   *AQ                `getter:"-"`
	typ *ast.InterfaceType `getter:"-"`
}

func (i *InterfaceType) Name() string {
	return ""
}
