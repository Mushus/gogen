package main

import "strings"

type tagValues []string

func newTagValue(tag string) tagValues {
	values := strings.Split(tag, ",")

	tagValues := tagValues{}
	for _, value := range values {
		if value == "" {
			continue
		}
		tagValues = append(tagValues, value)
	}

	return tagValues
}

func (t tagValues) contains(value string) bool {
	for _, v := range t {
		if v == value {
			return true
		}
	}

	return false
}

func (t tagValues) find(findFunc func(string) bool) string {
	for _, v := range t {
		if findFunc(v) {
			return v
		}
	}
	return ""
}
