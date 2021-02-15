package aq

import (
	"path/filepath"
	"strings"
)

type loadOptionSet struct {
	parseTestCode bool
	ignoreSuffix  []string
}

type LoadOption func(o *loadOptionSet)

func ParseTestCode(o *loadOptionSet) {
	o.parseTestCode = true
}

func IngoreSuffix(suffix ...string) LoadOption {
	return func(o *loadOptionSet) {
		o.ignoreSuffix = append(o.ignoreSuffix, suffix...)
	}
}

func (o loadOptionSet) filterFiles(filenames []string) []string {
	filtered := []string{}
L:
	for _, filename := range filenames {
		basename := filepath.Base(filename)

		if o.parseTestCode != strings.HasSuffix(basename, "_test") {
			continue
		}

		for _, is := range o.ignoreSuffix {
			if strings.HasSuffix(basename, is) {
				continue L
			}
		}

		filtered = append(filtered, filename)
	}

	return filenames
}
