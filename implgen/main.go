package main

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type cmdParams struct {
	recv       string
	isGen      bool
	structType string
	iface      string
	ifaceType  string
	ifacePkg   string
}

func createCmdParams(recv string, iface string) (cmdParams, error) {
	ptr := strings.HasPrefix(recv, "*")

	structType := recv
	if ptr {
		structType = recv[1:]
	}

	slashPos := strings.LastIndex(iface, "/")
	dotPos := strings.LastIndex(iface, ".")

	ifacePkg := ""
	ifaceType := ""

	if dotPos == -1 { // ex. "InterfaceType"
		ifaceType = iface
	} else if slashPos < dotPos { // ex. "io.Reader", "github.com/usern/pkg.InterfaceType"
		ifacePkg = iface[:dotPos]
		ifaceType = iface[dotPos+1:]
	} else {
		return cmdParams{}, errors.New("invalid argument iface")
	}

	return cmdParams{
		recv:       recv,
		structType: structType,
		ifacePkg:   ifacePkg,
		ifaceType:  ifaceType,
		iface:      iface,
	}, nil
}

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

func main() {
	prm, err := mapArgs()
	if err != nil {
		log.Fatalln("failed get params:", err)
	}
	gen := createGenerators(prm)

	if err := gen.collectParams(); err != nil {
		log.Fatalln("failed collect params:", err)
	}

	if err := gen.generate(); err != nil {
		log.Fatalln("failed generate params:", err)
	}
}

func mapArgs() (cmdParams, error) {
	args := os.Args
	if len(args) != 3 {
		return cmdParams{}, errors.New("invalid args")
	}

	typ, iface := args[1], args[2]

	return createCmdParams(typ, iface)
}

func createGenerators(prms cmdParams) *generator {
	return &generator{cmdParams: prms}
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
