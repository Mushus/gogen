package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"sort"

	"github.com/mattn/natural"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

type generator struct {
	cmdParams cmdParams
}

func (g *generator) generate() error {
	p := g.cmdParams
	fname := p.gofile
	typeName := p.implType

	m, err := loadInterfaceMethods(p.ifacePkg, p.ifaceType)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	var file *ast.File
	if _, err := os.Stat(fname); err == nil {
		file, err = parser.ParseFile(fset, fname, nil, 0)
		if err != nil {
			return err
		}
	} else {
		file = &ast.File{
			Name: &ast.Ident{Name: "main"},
		}
	}

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if c.Name() == "Decls" {
			if fd, ok := c.Node().(*ast.FuncDecl); ok {
				name := GetFuncName(fd)
				srcFT := m.Get(name)
				if srcFT == nil {
					return false
				}
				if GetRecvTypeName(fd) != typeName {
					return false
				}
				m.MarkAsExists(name)
				dstFT := fd.Type
				c.Replace(&ast.FuncDecl{
					Doc:  fd.Doc,
					Recv: fd.Recv,
					Name: fd.Name,
					Type: OverwriteMethod(dstFT, srcFT),
					Body: fd.Body,
				})
			}
		}
		return true
	}, nil)

	for _, ms := range m.GetDeficients() {
		file.Decls = append(file.Decls, generateMethods(typeName, ms.Name, ms.Type))
	}

	buf := new(bytes.Buffer)
	if err := format.Node(buf, fset, file); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fname, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func GetRecvTypeName(fd *ast.FuncDecl) string {
	if fd == nil || fd.Recv == nil {
		return ""
	}

	for _, f := range fd.Recv.List {
		// Reciver
		if ident, ok := f.Type.(*ast.Ident); ok {
			return ident.Name
		}
		// *Reciver
		if si, ok := f.Type.(*ast.StarExpr); ok {
			if ident, ok := si.X.(*ast.Ident); ok {
				return ident.Name
			}
		}
	}
	return ""
}

func OverwriteMethod(dst *ast.FuncType, src *ast.FuncType) *ast.FuncType {
	var pos token.Pos
	var params *ast.FieldList
	var results *ast.FieldList
	if dst != nil {
		pos = dst.Func
		params = dst.Params
		results = dst.Results
	}
	newFunc := &ast.FuncType{
		Func:    pos,
		Params:  overwriteMethodParams(params, src.Params),
		Results: overwriteMethodResults(results, src.Results),
	}
	return newFunc
}

func overwriteMethodParams(dst *ast.FieldList, src *ast.FieldList) *ast.FieldList {
	var opening token.Pos
	var closing token.Pos
	if dst != nil {
		opening = dst.Opening
		closing = dst.Closing
	}
	newFL := &ast.FieldList{
		Opening: opening,
		Closing: closing,
	}
	if src != nil {
		for i, srcF := range src.List {
			var doc *ast.CommentGroup
			var comment *ast.CommentGroup
			if dst != nil && i < len(dst.List) {
				dstF := dst.List[i]
				doc = dstF.Doc
				comment = dstF.Comment
			}

			newFL.List = append(newFL.List, &ast.Field{
				Doc:     doc,
				Names:   srcF.Names,
				Type:    srcF.Type,
				Comment: comment,
			})
		}
	}
	return newFL
}

func overwriteMethodResults(dst *ast.FieldList, src *ast.FieldList) *ast.FieldList {
	var opening token.Pos
	var closing token.Pos
	if dst != nil {
		opening = dst.Opening
		closing = dst.Closing
	}
	newFL := &ast.FieldList{
		Opening: opening,
		Closing: closing,
	}
	if src != nil {
		for i, srcF := range src.List {
			var doc *ast.CommentGroup
			var comment *ast.CommentGroup
			if dst != nil && i < len(dst.List) {
				dstF := dst.List[i]
				doc = dstF.Doc
				comment = dstF.Comment
			}

			newFL.List = append(newFL.List, &ast.Field{
				Doc:     doc,
				Names:   srcF.Names,
				Type:    srcF.Type,
				Comment: comment,
			})
		}
	}
	return newFL
}

func GetFuncName(fd *ast.FuncDecl) string {
	if fd == nil {
		return ""
	}
	return GetIdentName(fd.Name)
}

func GetTypeName(ts *ast.TypeSpec) string {
	if ts == nil {
		return ""
	}
	return GetIdentName(ts.Name)
}

func GetIdentName(i *ast.Ident) string {
	if i == nil {
		return ""
	}
	return i.Name
}

func GetBasicLitStr(bl *ast.BasicLit) string {
	if bl.Kind != token.STRING {
		return ""
	}

	tv, _ := types.Eval(token.NewFileSet(), nil, token.NoPos, bl.Value)
	return constant.StringVal(tv.Value)
}

func CollectInterfaceMethods(pkgs []*packages.Package, ifPkg string, ifName string) (map[string]*ast.FuncType, error) {
	pkg := FindPkg(pkgs, ifPkg)
	if pkg == nil {
		return nil, fmt.Errorf("undefined package: %s", ifPkg)
	}

	dict := map[string]*ast.FuncType{}
	if err := collectInterfaceMethods(dict, pkg, ifName); err != nil {
		return nil, err
	}
	return dict, nil
}

func collectInterfaceMethods(dict map[string]*ast.FuncType, pkg *packages.Package, ifName string) error {
	pkgDict, err := CollectImports(pkg)
	if err != nil {
		return err
	}

	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, s := range gd.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}
				tn := GetTypeName(ts)
				if tn != ifName {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				if it.Methods == nil {
					continue
				}

				for _, m := range it.Methods.List {
					if ft, ok := m.Type.(*ast.FuncType); ok {
						for _, il := range m.Names {
							name := GetIdentName(il)
							dict[name] = ft
						}
					} else if ident, ok := m.Type.(*ast.Ident); ok {
						embedName := GetIdentName(ident)
						collectInterfaceMethods(dict, pkg, embedName)
					} else if se, ok := m.Type.(*ast.SelectorExpr); ok {
						xi, ok := se.X.(*ast.Ident)
						if !ok {
							return errors.New("bad ast")
						}

						targetIfaceName := GetIdentName(se.Sel)
						targetPkgName := GetIdentName(xi)
						npkg, _ := pkgDict[targetPkgName]

						collectInterfaceMethods(dict, npkg, targetIfaceName)
					}
				}
				return nil
			}
		}
	}
	return nil
}

func CollectImports(pkg *packages.Package) (map[string]*packages.Package, error) {
	dict := map[string]*packages.Package{}
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, s := range gd.Specs {
				is, ok := s.(*ast.ImportSpec)
				if !ok {
					continue
				}

				pkgPath := GetBasicLitStr(is.Path)
				npkg, _ := pkg.Imports[pkgPath]
				if npkg == nil {
					return nil, fmt.Errorf("undefined package: %s", pkgPath)
				}

				name := GetIdentName(is.Name)
				if name == "" {
					name = npkg.Name
				}

				dict[name] = npkg
			}
		}
	}

	return dict, nil
}

func FindPkg(pkgs []*packages.Package, pkgPath string) *packages.Package {
	for _, pkg := range pkgs {
		if pkg.PkgPath == pkgPath {
			return pkg
		}
	}
	return nil
}

func loadInterfaceMethods(pkg string, typ string) (*Methods, error) {
	cfg := &packages.Config{
		Mode: packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.LoadImports,
	}

	pkgs, err := packages.Load(cfg, "pattern="+pkg)
	if err != nil {
		return nil, fmt.Errorf("failed load package: %w", err)
	}

	fl, err := CollectInterfaceMethods(pkgs, pkg, typ)
	if err != nil {
		return nil, fmt.Errorf("cannot find interface: %w", err)
	}

	return newMethods(fl), nil
}

type Methods struct {
	methods map[string]*ast.FuncType
	exists  map[string]bool
}

func newMethods(methods map[string]*ast.FuncType) *Methods {
	return &Methods{
		methods: methods,
		exists:  map[string]bool{},
	}
}

func (m Methods) Get(methodName string) *ast.FuncType {
	return m.methods[methodName]
}

func (m Methods) MarkAsExists(methodName string) {
	m.exists[methodName] = true
}

func (m Methods) GetDeficients() []MethodSet {
	l := []MethodSet{}
	for name, typ := range m.methods {
		if m.exists[name] {
			continue
		}
		l = append(l, MethodSet{
			Name: name,
			Type: typ,
		})
	}

	sort.Slice(l, func(i, j int) bool {
		return natural.NaturalComp(l[i].Name, l[j].Name) < 0
	})

	return l
}

type MethodSet struct {
	Name string
	Type *ast.FuncType
}

func generateMethods(recvType string, name string, typ *ast.FuncType) *ast.FuncDecl {
	recvVar := string([]rune(recvType)[0])
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{{Name: recvVar}},
				Type:  &ast.Ident{Name: recvType},
			}},
		},
		Name: &ast.Ident{Name: name},
		Type: OverwriteMethod(&ast.FuncType{}, typ),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: "panic"},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf("%#v", "not implemented"),
							},
						},
					},
				},
			},
		},
	}
}
