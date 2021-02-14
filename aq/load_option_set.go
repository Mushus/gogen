package aq

import (
	"fmt"
	"io/ioutil"
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

func (o loadOptionSet) getGoFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir %#v: %w", dir, err)
	}

	filenames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if filepath.Ext(filename) != ".go" {
			continue
		}

		basename := filepath.Base(filename)
		if o.parseTestCode != strings.HasSuffix(basename, "_test") {
			continue
		}

		filenames = append(filenames, filepath.Join(dir, filename))
	}

	return filenames, nil
}
