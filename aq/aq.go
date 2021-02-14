package aq

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

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

func (i *Instance) MustLoadDir(dir string, opts ...LoadOption) *Instance {
	err := i.LoadDir(dir, opts...)
	if err != nil {
		panic(err)
	}
	return i
}

func (i *Instance) LoadDir(dir string, opts ...LoadOption) error {
	optSet := &loadOptionSet{
		parseTestCode: false,
	}
	for _, opt := range opts {
		opt(optSet)
	}

	filenames, err := optSet.getGoFiles(dir)
	if err != nil {
		return err
	}

LoadFile:
	for _, filename := range filenames {
		// ignore suffix
		basename := filepath.Base(filename)
		for _, suffix := range optSet.ignoreSuffix {
			if strings.HasSuffix(basename, suffix) {
				continue LoadFile
			}
		}

		file, err := parser.ParseFile(i.fileSet, filename, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		i.files = append(i.files, file)
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

func (i *Instance) Files() []*File {
	files := make([]*File, 0, len(i.files))
	for _, file := range i.files {
		files = append(files, createFile(file))
	}

	return files
}

func (i *Instance) File() *File {
	for _, file := range i.files {
		return createFile(file)
	}

	return nil
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

func (i *Instance) Funcs() FuncList {
	funcs := make(FuncList, 0)
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
