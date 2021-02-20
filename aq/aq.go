package aq

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
)

//go:generate go run github.com/Mushus/gogen/structgen -construct Field,File,FuncDecl,FuncType,ImportSpec,InterfaceSpec,InterfaceType,StructSpec,TypeSpec -list Field,File,FuncDecl,FuncType,ImportSpec,InterfaceSpec,InterfaceType,StructSpec,TypeSpec

type AQ struct {
	fileSet *token.FileSet
	files   []*ast.File
}

func New() *AQ {
	fileSet := token.NewFileSet()

	return &AQ{
		fileSet: fileSet,
	}
}

func LoadFile(filename string) (*AQ, error) {
	a := New()
	err := a.LoadFile(filename)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func MustLoadFile(filename string) *AQ {
	a, err := LoadFile(filename)
	if err != nil {
		panic(err)
	}

	return a
}

func (a *AQ) MustLoadFile(filename string) *AQ {
	err := a.LoadDir(filename)
	if err != nil {
		panic(err)
	}
	return a
}

func (a *AQ) LoadFile(filename string) error {
	file, err := parser.ParseFile(a.fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	a.files = append(a.files, file)
	return nil
}

func MustLoadDir(dir string, opts ...LoadOption) *AQ {
	a, err := LoadDir(dir, opts...)
	if err != nil {
		panic(err)
	}

	return a
}

func LoadDir(dir string, opts ...LoadOption) (*AQ, error) {
	a := New()
	err := a.LoadDir(dir, opts...)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *AQ) MustLoadDir(dir string, opts ...LoadOption) *AQ {
	err := a.LoadDir(dir, opts...)
	if err != nil {
		panic(err)
	}
	return a
}

func (a *AQ) LoadDir(dir string, opts ...LoadOption) error {
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
		if err := a.LoadFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func MustImport(targetPackage string, path string) *AQ {
	a, err := Import(targetPackage, path)
	if err != nil {
		panic(err)
	}
	return a
}

func Import(targetPackage string, path string) (*AQ, error) {
	a := New()

	if err := a.Import(targetPackage, path); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AQ) MustImport(targetPackage string, path string) error {
	if err := a.Import(targetPackage, path); err != nil {
		panic(err)
	}
	return nil
}

func (a *AQ) Import(targetPackage string, path string) error {
	pkg, err := build.Import(targetPackage, path, 0)
	if err != nil {
		return err
	}

	dir := pkg.Dir

	files := []string{}
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)

	for _, f := range files {
		if err := a.LoadFile(filepath.Join(dir, f)); err != nil {
			return err
		}
	}

	return nil
}

func (a *AQ) MustLoadFromSource(source []byte) *AQ {
	err := a.LoadFromSource(source)
	if err != nil {
		panic(err)
	}
	return a
}

func (a *AQ) LoadFromSource(source []byte) error {
	file, err := parser.ParseFile(a.fileSet, "", source, parser.ParseComments)
	if err != nil {
		return err
	}
	a.files = append(a.files, file)
	return nil
}

func (a *AQ) Package() string {
	for _, file := range a.files {
		if file.Name == nil {
			continue
		}

		return file.Name.Name
	}
	return ""
}

func (i AQ) Packages() []string {
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

func (a *AQ) Files() Files {
	l := make(Files, 0, len(a.files))
	for _, f := range a.files {
		l = append(l, NewFile(a, f))
	}

	return l
}

func (a *AQ) File() *File {
	for _, f := range a.files {
		return NewFile(a, f)
	}

	return nil
}

func (a *AQ) Types() TypeSpecs {
	l := TypeSpecs{}
	for _, f := range a.Files() {
		l = append(l, f.Types()...)
	}
	return l
}

func (a *AQ) Interfaces() InterfaceSpecs {
	return a.Types().Interfaces()
}

func (a *AQ) Structs() StructSpecs {
	return a.Types().Structs()
}

func (a *AQ) Funcs() FuncDecls {
	l := make(FuncDecls, 0)
	for _, f := range a.files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			l = append(l, NewFuncDecl(fd))
		}
	}
	return l
}
