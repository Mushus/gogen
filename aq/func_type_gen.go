// Code generated by structgen; DO NOT EDIT.

package aq

import (
	"go/ast"
)

func NewFuncType(
	typ *ast.FuncType,
) *FuncType {
	c := &FuncType{
		typ: typ,
	}

	return c
}

type FuncTypes []*FuncType

func (r FuncTypes) Chunk(size int) []FuncTypes {
	list := []FuncTypes{}
	chunk := FuncTypes{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = FuncTypes{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r FuncTypes) Compact() FuncTypes {
	l := FuncTypes{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r FuncTypes) Concat(l FuncTypes) FuncTypes {
	return append(append(FuncTypes{}, r...), l...)
}

func (r FuncTypes) Copy() FuncTypes {
	dist := make(FuncTypes, len(r))
	copy(dist, r)
	return dist
}

func (r FuncTypes) Count() int {
	return len(r)
}

func (r FuncTypes) Each(f func(i int, v *FuncType)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r FuncTypes) Exists() bool {
	return r != nil && len(r) > 0
}

func (r FuncTypes) Every(f func(i int, v *FuncType) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r FuncTypes) Filter(funcs ...func(i int, v *FuncType) bool) FuncTypes {
	list := FuncTypes{}
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

func (r FuncTypes) Find(funcs ...func(i int, v *FuncType) bool) *FuncType {
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

func (r FuncTypes) FindIndex(funcs ...func(i int, v *FuncType) bool) int {
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

func (r FuncTypes) First() *FuncType {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r FuncTypes) ForPage(pageNo int, size int) FuncTypes {
	rLen := len(r)
	list := make(FuncTypes, 0, size)
	for i, k := pageNo*size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r FuncTypes) Get(i int) *FuncType {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r FuncTypes) Has(f func(i int, v *FuncType) bool) bool {
	return r.Some(f)
}

func (r FuncTypes) IsEmpty() bool {
	return len(r) == 0
}

func (r FuncTypes) IsNotEmpty() bool {
	return len(r) > 0
}

func (r FuncTypes) Last() *FuncType {
	if len(r) == 0 {
		return nil
	}
	return r[len(r)-1]
}

func (r FuncTypes) Reverse() FuncTypes {
	list := make(FuncTypes, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r FuncTypes) Some(f func(i int, v *FuncType) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r FuncTypes) Take(size int) FuncTypes {
	if len(r) > size {
		return r
	}
	return r[:size]
}

func (r FuncTypes) Typs() []*ast.FuncType {
	l := make([]*ast.FuncType, 0, len(r))
	for _, r := range r {
		l = append(l, r.typ)
	}
	return l
}

func (r FuncTypes) ParamsList() []Fields {
	l := make([]Fields, 0, len(r))
	for _, r := range r {
		l = append(l, r.Params())
	}
	return l
}

func (r FuncTypes) ResultsList() []Fields {
	l := make([]Fields, 0, len(r))
	for _, r := range r {
		l = append(l, r.Results())
	}
	return l
}
