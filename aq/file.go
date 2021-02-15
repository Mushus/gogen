package aq

import "go/ast"

type File struct {
	instance *Instance `getter:"-"`
	file     *ast.File `getter:"-"`
}

func createFile(file *ast.File) *File {
	return &File{
		file: file,
	}
}

func (f *File) Exists() bool {
	return f != nil
}

func (f *File) Imports() Imports {
	if !f.Exists() {
		return nil
	}

	imports := make(Imports, 0, len(f.file.Imports))

	for _, i := range f.file.Imports {
		imports = append(imports, NewImport(i))
	}

	return imports
}

func (f *File) Package() string {
	if !f.Exists() {
		return ""
	}

	return safeIdentName(f.file.Name)
}

func (f *File) Decls() Decls {
	l := make(Decls, 0, len(f.file.Decls))
	for _, decl := range f.file.Decls {
		l = append(l, NewDecl(decl))
	}

	return l
}

func (f *File) Types() {
	// return f.Decls()
}
