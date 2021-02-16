package main

import (
	"fmt"

	"github.com/Mushus/gogen/aq"
)

type generator struct {
	cmdParams cmdParams
}

func (g *generator) collectParams() error {
	genDir := "."
	ifacePkg := g.cmdParams.ifacePkg

	aqi := aq.New()
	if err := aqi.Import(ifacePkg, genDir); err != nil {
		return err
	}

	ifaces := aqi.Decls().Types().Interfaces()
	fmt.Printf("%#v", ifaces)

	return nil
}

func (g *generator) generate() error {
	return nil
}
