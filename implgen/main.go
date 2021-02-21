package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type cmdParams struct {
	recv      string
	isGen     bool
	implType  string
	iface     string
	ifaceType string
	ifacePkg  string
	dst       string
}

func createCmdParams(dst string, recv string, iface string) (cmdParams, error) {
	ptr := strings.HasPrefix(recv, "*")

	implType := recv
	if ptr {
		implType = recv[1:]
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
		recv:      recv,
		implType:  implType,
		ifacePkg:  ifacePkg,
		ifaceType: ifaceType,
		iface:     iface,
		dst:       dst,
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
		usage(err)
		os.Exit(1)
	}

	gen := createGenerators(prm)
	if err := gen.generate(); err != nil {
		usage(err)
		os.Exit(1)
	}
}

func mapArgs() (cmdParams, error) {
	dst := flag.String("dst", os.Getenv("GOFILE"), "destination go file path")

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return cmdParams{}, fmt.Errorf("invalid args: the number of args must be 2, got %d", len(args))
	}

	typ, iface := args[0], args[1]

	return createCmdParams(*dst, typ, iface)
}

func usage(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS...] ReciverType path/to/go/package.InterfaceType\n", os.Args[0])
	flag.PrintDefaults()
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
