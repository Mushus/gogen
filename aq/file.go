package aq

import "go/ast"

type File struct {
	instance *Instance
	file     *ast.File
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

	imports := make([]*Import, 0, len(f.file.Imports))

	for _, i := range f.file.Imports {
		imports = append(imports, createImport(i))
	}

	return imports
}

func (f *File) Package() string {
	if !f.Exists() {
		return ""
	}

	return safeIdentName(f.file.Name)
}
