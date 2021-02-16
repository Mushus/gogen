package aq

import "go/ast"

type File struct {
	aq   *AQ       `getter:"-"`
	file *ast.File `getter:"-"`
}

func (f *File) Exists() bool {
	return f != nil
}

func (f *File) Imports() ImportSpecs {
	if f == nil {
		return nil
	}

	l := make(ImportSpecs, 0, len(f.file.Imports))
	for _, i := range f.file.Imports {
		l = append(l, NewImportSpec(i))
	}

	return l
}

func (f *File) Package() string {
	if f == nil {
		return ""
	}

	return safeIdentName(f.file.Name)
}

func (f *File) Types() TypeSpecs {
	if f == nil {
		return nil
	}

	l := make(TypeSpecs, 0, len(f.file.Scope.Objects))
	for _, o := range f.file.Scope.Objects {
		ts, ok := o.Decl.(*ast.TypeSpec)
		if ok {
			l = append(l, NewTypeSpec(f.aq, f, ts))
		}
	}

	return l
}

func (f *File) Structs() Structs {
	return f.Types().Structs()
}
