// Code generated by structgen; DO NOT EDIT.

package aq

import (
	"go/ast"
)

func NewTypeSpec(
	aq *AQ,
	file *File,
	spec *ast.TypeSpec,
) *TypeSpec {
	c := &TypeSpec{
		aq:   aq,
		file: file,
		spec: spec,
	}

	return c
}

type TypeSpecs []*TypeSpec

func (r TypeSpecs) Chunk(size int) []TypeSpecs {
	list := []TypeSpecs{}
	chunk := TypeSpecs{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = TypeSpecs{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r TypeSpecs) Compact() TypeSpecs {
	l := TypeSpecs{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r TypeSpecs) Concat(l TypeSpecs) TypeSpecs {
	return append(append(TypeSpecs{}, r...), l...)
}

func (r TypeSpecs) Copy() TypeSpecs {
	dist := make(TypeSpecs, len(r))
	copy(dist, r)
	return dist
}

func (r TypeSpecs) Count() int {
	return len(r)
}

func (r TypeSpecs) Each(f func(i int, v *TypeSpec)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r TypeSpecs) Exists() bool {
	return r != nil && len(r) > 0
}

func (r TypeSpecs) Every(f func(i int, v *TypeSpec) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r TypeSpecs) Filter(funcs ...func(i int, v *TypeSpec) bool) TypeSpecs {
	list := TypeSpecs{}
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		list = append(list, v)
	}
	return list
}

func (r TypeSpecs) Find(funcs ...func(i int, v *TypeSpec) bool) *TypeSpec {
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		return v
	}
	return nil
}

func (r TypeSpecs) FindIndex(funcs ...func(i int, v *TypeSpec) bool) int {
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		return i
	}
	return -1
}

func (r TypeSpecs) First() *TypeSpec {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r TypeSpecs) ForPage(pageNo int, size int) TypeSpecs {
	rLen := len(r)
	list := make(TypeSpecs, 0, size)
	for i, k := pageNo*size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r TypeSpecs) Get(i int) *TypeSpec {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r TypeSpecs) Has(f func(i int, v *TypeSpec) bool) bool {
	return r.Some(f)
}

func (r TypeSpecs) IsEmpty() bool {
	return len(r) == 0
}

func (r TypeSpecs) IsNotEmpty() bool {
	return len(r) > 0
}

func (r TypeSpecs) Last() *TypeSpec {
	if len(r) == 0 {
		return nil
	}
	return r[len(r)-1]
}

func (r TypeSpecs) Reverse() TypeSpecs {
	list := make(TypeSpecs, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r TypeSpecs) Some(f func(i int, v *TypeSpec) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r TypeSpecs) Take(size int) TypeSpecs {
	if len(r) > size {
		return r
	}
	return r[:size]
}
