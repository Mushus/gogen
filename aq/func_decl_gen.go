// Code generated by structgen; DO NOT EDIT.

package aq

import (
	"go/ast"
)

func NewFuncDecl(
	decl *ast.FuncDecl,
) *FuncDecl {
	c := &FuncDecl{
		decl: decl,
	}

	return c
}

type FuncDecls []*FuncDecl

func (r FuncDecls) Chunk(size int) []FuncDecls {
	list := []FuncDecls{}
	chunk := FuncDecls{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = FuncDecls{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r FuncDecls) Compact() FuncDecls {
	l := FuncDecls{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r FuncDecls) Concat(l FuncDecls) FuncDecls {
	return append(append(FuncDecls{}, r...), l...)
}

func (r FuncDecls) Copy() FuncDecls {
	dist := make(FuncDecls, len(r))
	copy(dist, r)
	return dist
}

func (r FuncDecls) Count() int {
	return len(r)
}

func (r FuncDecls) Each(f func(i int, v *FuncDecl)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r FuncDecls) Exists() bool {
	return r != nil && len(r) > 0
}

func (r FuncDecls) Every(f func(i int, v *FuncDecl) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r FuncDecls) Filter(funcs ...func(i int, v *FuncDecl) bool) FuncDecls {
	list := FuncDecls{}
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

func (r FuncDecls) Find(funcs ...func(i int, v *FuncDecl) bool) *FuncDecl {
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

func (r FuncDecls) FindIndex(funcs ...func(i int, v *FuncDecl) bool) int {
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

func (r FuncDecls) First() *FuncDecl {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r FuncDecls) ForPage(pageNo int, size int) FuncDecls {
	rLen := len(r)
	list := make(FuncDecls, 0, size)
	for i, k := pageNo*size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r FuncDecls) Get(i int) *FuncDecl {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r FuncDecls) Has(f func(i int, v *FuncDecl) bool) bool {
	return r.Some(f)
}

func (r FuncDecls) IsEmpty() bool {
	return len(r) == 0
}

func (r FuncDecls) IsNotEmpty() bool {
	return len(r) > 0
}

func (r FuncDecls) Last() *FuncDecl {
	if len(r) == 0 {
		return nil
	}
	return r[len(r)-1]
}

func (r FuncDecls) Reverse() FuncDecls {
	list := make(FuncDecls, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r FuncDecls) Some(f func(i int, v *FuncDecl) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r FuncDecls) Take(size int) FuncDecls {
	if len(r) > size {
		return r
	}
	return r[:size]
}

func (r FuncDecls) Names() []string {
	l := make([]string, 0, len(r))
	for _, r := range r {
		l = append(l, r.Name())
	}
	return l
}

func (r FuncDecls) Recvs() []*Field {
	l := make([]*Field, 0, len(r))
	for _, r := range r {
		l = append(l, r.Recv())
	}
	return l
}

func (r FuncDecls) Types() []*FuncType {
	l := make([]*FuncType, 0, len(r))
	for _, r := range r {
		l = append(l, r.Type())
	}
	return l
}

func (r FuncDecls) ParamsList() []Fields {
	l := make([]Fields, 0, len(r))
	for _, r := range r {
		l = append(l, r.Params())
	}
	return l
}

func (r FuncDecls) ResultsList() []Fields {
	l := make([]Fields, 0, len(r))
	for _, r := range r {
		l = append(l, r.Results())
	}
	return l
}
