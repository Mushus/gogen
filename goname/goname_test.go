package goname_test

import (
	"testing"

	"github.com/Mushus/gogen/goname"
)

func TestUpperCamelCase(t *testing.T) {
	cases := []struct {
		in     string
		expect string
	}{
		{
			in:     "hello World!",
			expect: "HelloWorld",
		},
		{
			in:     "foo/bar.baz",
			expect: "FooBarBaz",
		},
		{
			in:     "thisIsAPen",
			expect: "ThisIsAPen",
		},
		{
			in:     "hello123world456!",
			expect: "Hello123World456",
		},
		{
			in:     "utf8mb4",
			expect: "UTF8Mb4",
		},
	}
	for _, c := range cases {
		actual := goname.UpperCamelCase(c.in)
		if c.expect != actual {
			t.Fatalf("expect %#v got %#v", c.expect, actual)
		}
	}
}

func TestLowerCamelCase(t *testing.T) {
	cases := []struct {
		in     string
		expect string
	}{
		{
			in:     "hello World!",
			expect: "helloWorld",
		},
		{
			in:     "foo/bar.baz",
			expect: "fooBarBaz",
		},
		{
			in:     "thisIsAPen",
			expect: "thisIsAPen",
		},
		{
			in:     "hello123world456!",
			expect: "hello123World456",
		},
		{
			in:     "utf8mb4",
			expect: "utf8Mb4",
		},
	}
	for _, c := range cases {
		actual := goname.LowerCamelCase(c.in)
		if c.expect != actual {
			t.Fatalf("expect %#v got %#v", c.expect, actual)
		}
	}
}
