package aq

import (
	"go/ast"
	"reflect"
)

type Tag struct {
	tag *ast.BasicLit
}

func createTag(tag *ast.BasicLit) *Tag {
	if tag == nil {
		return nil
	}

	return &Tag{
		tag: tag,
	}
}

func (t *Tag) Exists() bool {
	return t != nil
}

func (t *Tag) Body() reflect.StructTag {
	if !t.Exists() {
		return ""
	}

	if t.tag.Value == "" {
		return ""
	}

	tag := stringLiteral(t.tag)
	return reflect.StructTag(tag)
}
