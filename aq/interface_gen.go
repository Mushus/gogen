// Code generated by structgen; DO NOT EDIT.

package aq

import (
	"go/ast"
)

func NewInterface(
	i *AQ,
	typ *ast.InterfaceType,
) *Interface {
	c := &Interface{
		i:   i,
		typ: typ,
	}

	return c
}

type Interfaces []*Interface

func (r Interfaces) Chunk(size int) []Interfaces {
	list := []Interfaces{}
	chunk := Interfaces{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = Interfaces{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r Interfaces) Compact() Interfaces {
	l := Interfaces{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r Interfaces) Concat(l Interfaces) Interfaces {
	return append(append(Interfaces{}, r...), l...)
}

func (r Interfaces) Copy() Interfaces {
	dist := make(Interfaces, len(r))
	copy(dist, r)
	return dist
}

func (r Interfaces) Count() int {
	return len(r)
}

func (r Interfaces) Each(f func(i int, v *Interface)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r Interfaces) Exists() bool {
	return r != nil && len(r) > 0
}

func (r Interfaces) Every(f func(i int, v *Interface) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r Interfaces) Filter(funcs ...func(i int, v *Interface) bool) Interfaces {
	list := Interfaces{}
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

func (r Interfaces) Find(funcs ...func(i int, v *Interface) bool) *Interface {
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

func (r Interfaces) FindIndex(funcs ...func(i int, v *Interface) bool) int {
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

func (r Interfaces) First() *Interface {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r Interfaces) ForPage(pageNo int, size int) Interfaces {
	rLen := len(r)
	list := make(Interfaces, 0, size)
	for i, k := pageNo*size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r Interfaces) Get(i int) *Interface {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r Interfaces) Has(f func(i int, v *Interface) bool) bool {
	return r.Some(f)
}

func (r Interfaces) IsEmpty() bool {
	return len(r) == 0
}

func (r Interfaces) IsNotEmpty() bool {
	return len(r) > 0
}

func (r Interfaces) Last() *Interface {
	if len(r) == 0 {
		return nil
	}
	return r[len(r)-1]
}

func (r Interfaces) Reverse() Interfaces {
	list := make(Interfaces, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r Interfaces) Some(f func(i int, v *Interface) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r Interfaces) Take(size int) Interfaces {
	if len(r) > size {
		return r
	}
	return r[:size]
}