// Code generated by structgen; DO NOT EDIT.

package aq

import (
	"go/ast"
)

func NewFile(
	aq *AQ,
	file *ast.File,
) *File {
	c := &File{
		aq:   aq,
		file: file,
	}

	return c
}

type Files []*File

func (r Files) Chunk(size int) []Files {
	list := []Files{}
	chunk := Files{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = Files{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r Files) Compact() Files {
	l := Files{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r Files) Concat(l Files) Files {
	return append(append(Files{}, r...), l...)
}

func (r Files) Copy() Files {
	dist := make(Files, len(r))
	copy(dist, r)
	return dist
}

func (r Files) Count() int {
	return len(r)
}

func (r Files) Each(f func(i int, v *File)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r Files) Exists() bool {
	return r != nil && len(r) > 0
}

func (r Files) Every(f func(i int, v *File) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r Files) Filter(funcs ...func(i int, v *File) bool) Files {
	list := Files{}
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

func (r Files) Find(funcs ...func(i int, v *File) bool) *File {
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

func (r Files) FindIndex(funcs ...func(i int, v *File) bool) int {
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

func (r Files) First() *File {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r Files) ForPage(pageNo int, size int) Files {
	rLen := len(r)
	list := make(Files, 0, size)
	for i, k := pageNo*size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r Files) Get(i int) *File {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r Files) Has(f func(i int, v *File) bool) bool {
	return r.Some(f)
}

func (r Files) IsEmpty() bool {
	return len(r) == 0
}

func (r Files) IsNotEmpty() bool {
	return len(r) > 0
}

func (r Files) Last() *File {
	if len(r) == 0 {
		return nil
	}
	return r[len(r)-1]
}

func (r Files) Reverse() Files {
	list := make(Files, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r Files) Some(f func(i int, v *File) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r Files) Take(size int) Files {
	if len(r) > size {
		return r
	}
	return r[:size]
}

func (r Files) ImportsList() []ImportSpecs {
	l := make([]ImportSpecs, 0, len(r))
	for _, r := range r {
		l = append(l, r.Imports())
	}
	return l
}

func (r Files) Packages() []string {
	l := make([]string, 0, len(r))
	for _, r := range r {
		l = append(l, r.Package())
	}
	return l
}

func (r Files) StructsList() []StructSpecs {
	l := make([]StructSpecs, 0, len(r))
	for _, r := range r {
		l = append(l, r.Structs())
	}
	return l
}
