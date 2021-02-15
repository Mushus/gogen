package aq

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
)

//go:generate go run github.com/Mushus/gogen/structgen -construct Decl,Field,File,Func,FuncType,Import,Struct,TypeDef,Decl -list Decl,Field,File,Func,FuncType,Import,Struct,TypeDef,Decl

type Instance struct {
	fileSet *token.FileSet
	files   []*ast.File
}

func New() *Instance {
	fileSet := token.NewFileSet()

	return &Instance{
		fileSet: fileSet,
	}
}

func (i *Instance) MustLoadFile(filename string) *Instance {
	err := i.LoadDir(filename)
	if err != nil {
		panic(err)
	}
	return i
}

func (i *Instance) LoadFile(filename string) error {
	file, err := parser.ParseFile(i.fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	i.files = append(i.files, file)
	return nil
}

func (i *Instance) MustLoadDir(dir string, opts ...LoadOption) *Instance {
	err := i.LoadDir(dir, opts...)
	if err != nil {
		panic(err)
	}
	return i
}

func (i *Instance) LoadDir(dir string, opts ...LoadOption) error {
	pkg, err := build.ImportDir(dir, 0)
	if err != nil {
		return err
	}

	filenames := make([]string, 0, len(pkg.GoFiles)+len(pkg.CgoFiles))
	filenames = append(filenames, pkg.GoFiles...)
	filenames = append(filenames, pkg.CgoFiles...)

	optSet := &loadOptionSet{
		parseTestCode: false,
	}
	for _, opt := range opts {
		opt(optSet)
	}

	filenames = optSet.filterFiles(filenames)

	for _, filename := range filenames {
		if err := i.LoadFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func (i *Instance) MustLoadFromSource(source []byte) *Instance {
	err := i.LoadFromSource(source)
	if err != nil {
		panic(err)
	}
	return i
}

func (i *Instance) LoadFromSource(source []byte) error {
	file, err := parser.ParseFile(i.fileSet, "", source, parser.ParseComments)
	if err != nil {
		return err
	}
	i.files = append(i.files, file)
	return nil
}

func (i Instance) Package() string {
	for _, file := range i.files {
		if file.Name == nil {
			continue
		}

		return file.Name.Name
	}
	return ""
}

func (i Instance) Packages() []string {
	pkgMap := map[string]bool{}
	pkgList := make([]string, 0)
	for _, file := range i.files {
		if file.Name == nil {
			continue
		}

		pkg := file.Name.Name
		if pkgMap[pkg] {
			continue
		}
		pkgMap[pkg] = true
		pkgList = append(pkgList, pkg)
	}

	return pkgList
}

func (i *Instance) Files() Files {
	files := make(Files, 0, len(i.files))
	for _, file := range i.files {
		files = append(files, NewFile(i, file))
	}

	return files
}

func (i *Instance) File() *File {
	for _, file := range i.files {
		return createFile(file)
	}

	return nil
}

func (i *Instance) Decls() {
	// for _, f := range i.files {
	// 	f.Decls
	// }
}

func (i *Instance) Structs() Structs {
	structs := make(Structs, 0)
	for _, f := range i.files {
		aqFile := createFile(f)
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				structs = append(structs, createStruct(aqFile, typeSpec, structType))
			}
		}
	}

	return structs
}

func (i *Instance) Funcs() Funcs {
	funcs := make(Funcs, 0)
	for _, f := range i.files {
		for _, decl := range f.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			funcs = append(funcs, createFunc(funcDecl))
		}
	}
	return funcs
}

func (i *Instance) Types() {
	// types := make(TypeDefList, 0)
	// for _, f := range i.files {
	// 	for _, decl := range f.Decls {
	// 		TypeDecl
	// 	}
	// }
}
