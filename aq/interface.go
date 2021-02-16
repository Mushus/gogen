package aq

import "go/ast"

type Interface struct {
	i   *AQ                `getter:"-"`
	typ *ast.InterfaceType `getter:"-"`
}

func (i *Interface) Name() string {
	return ""
}
