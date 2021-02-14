package main

import (
	"flag"
	"log"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type generateTargets struct {
	m map[string]int
	l []*targetStruct
}

func newGenerateTarget() *generateTargets {
	return &generateTargets{
		m: map[string]int{},
		l: []*targetStruct{},
	}
}

func (g *generateTargets) get(name string) *targetStruct {
	if idx, ok := g.m[name]; ok {
		return g.l[idx]
	}

	s := &targetStruct{
		Name: name,
	}
	g.m[name] = len(g.l)
	g.l = append(g.l, s)
	return s
}

func (g *generateTargets) list() []*targetStruct {
	return g.l
}

type targetStruct struct {
	Name      string
	Construct bool
	Getter    bool
	Setter    bool
	List      bool
}

type fileGenerator interface {
	collectParams() error
	generate() error
}

type GenerateStruct struct {
	Construct         bool
	Getter            bool
	Setter            bool
	Package           string
	Imports           []Import
	Name              string
	NameUpper         string
	NameSnake         string
	Type              string
	Fields            []StructField
	ConstructTemplate string
}

type Import struct {
	Name        string
	Path        string
	PathLiteral string
}

type StructField struct {
	Name         string
	NameUpper    string
	Type         string
	IsPublic     bool
	Tag          reflect.StructTag
	ConstructTag ConstructTag
}

type ConstructTag struct {
	Ignore bool
}

var (
	path      = flag.String("path", ".", "directory")
	construct = flag.String("construct", "", "struct type")
	getter    = flag.String("getter", "", "struct type")
	setter    = flag.String("setter", "", "struct type")
	list      = flag.String("list", "", "struct type")
)

func main() {
	gt := mapFlagValue()
	gens := createGenerators(gt)

	if err := collectParams(gens); err != nil {
		log.Fatalln("failed collect params: ", err)
	}

	if err := generate(gens); err != nil {
		log.Fatalln("failed generate params:", err)
	}
}

func mapFlagValue() *generateTargets {
	flag.Parse()

	gt := newGenerateTarget()

	for _, target := range splitFlag(construct) {
		gt.get(target).Construct = true
	}

	for _, target := range splitFlag(getter) {
		gt.get(target).Getter = true
	}

	for _, target := range splitFlag(setter) {
		gt.get(target).Setter = true
	}

	for _, target := range splitFlag(list) {
		gt.get(target).List = true
	}

	return gt
}

func createGenerators(gt *generateTargets) []fileGenerator {
	gens := []fileGenerator{}
	for _, ts := range gt.list() {
		gen := newCombineGenerator(*ts, *path)
		gens = append(gens, gen)
	}
	return gens
}

func collectParams(gens []fileGenerator) error {
	for _, gen := range gens {
		if err := gen.collectParams(); err != nil {
			return err
		}
	}
	return nil
}

func generate(gens []fileGenerator) error {
	for _, gen := range gens {

		if err := gen.generate(); err != nil {
			return errors.Wrap(err, "cannot generate files")
		}
	}

	return nil
}

func splitFlag(flag *string) []string {
	cleanTexts := []string{}

	if flag == nil {
		return cleanTexts
	}

	texts := strings.Split(*flag, ",")
	for _, text := range texts {
		cleanText := strings.TrimSpace(text)
		if cleanText == "" {
			continue
		}

		cleanTexts = append(cleanTexts, cleanText)
	}

	return cleanTexts
}
