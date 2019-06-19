package httphead

import (
	"bytes"
	"testing"
)

func TestScannerSkipEscaped(t *testing.T) {
	for _, test := range []struct {
		in  []byte
		c   byte
		pos int
	}{
		{
			in:  []byte(`foo,bar`),
			c:   ',',
			pos: 4,
		},
		{
			in:  []byte(`foo\,bar,baz`),
			c:   ',',
			pos: 9,
		},
	} {
		s := NewScanner(test.in)
		s.SkipEscaped(test.c)
		if act, exp := s.pos, test.pos; act != exp {
			t.Errorf("unexpected scanner pos: %v; want %v", act, exp)
		}
	}
}

type readCase struct {
	label string
	in    []byte
	out   []byte
	err   bool
}

var quotedStringCases = []readCase{
	{
		label: "nonterm",
		in:    []byte(`"`),
		out:   []byte(``),
		err:   true,
	},
	{
		label: "empty",
		in:    []byte(`""`),
		out:   []byte(``),
	},
	{
		label: "simple",
		in:    []byte(`"hello, world!"`),
		out:   []byte(`hello, world!`),
	},
	{
		label: "quoted",
		in:    []byte(`"hello, \"world\"!"`),
		out:   []byte(`hello, "world"!`),
	},
	{
		label: "quoted",
		in:    []byte(`"\"hello\", \"world\"!"`),
		out:   []byte(`"hello", "world"!`),
	},
}

var commentCases = []readCase{
	{
		label: "nonterm",
		in:    []byte(`(hello`),
		out:   []byte(``),
		err:   true,
	},
	{
		label: "empty",
		in:    []byte(`()`),
		out:   []byte(``),
	},
	{
		label: "simple",
		in:    []byte(`(hello)`),
		out:   []byte(`hello`),
	},
	{
		label: "quoted",
		in:    []byte(`(hello\)\(world)`),
		out:   []byte(`hello)(world`),
	},
	{
		label: "nested",
		in:    []byte(`(hello(world))`),
		out:   []byte(`hello(world)`),
	},
}

type readTest struct {
	label string
	cases []readCase
	fn    func(*Scanner) bool
}

var readTests = []readTest{
	{
		"ReadString",
		quotedStringCases,
		(*Scanner).fetchQuotedString,
	},
	{
		"ReadComment",
		commentCases,
		(*Scanner).fetchComment,
	},
}

func TestScannerRead(t *testing.T) {
	for _, bunch := range readTests {
		for _, test := range bunch.cases {
			t.Run(bunch.label+" "+test.label, func(t *testing.T) {
				l := &Scanner{data: []byte(test.in)}
				if ok := bunch.fn(l); ok != !test.err {
					t.Errorf("l.%s() = %v; want %v", bunch.label, ok, !test.err)
					return
				}
				if !bytes.Equal(test.out, l.itemBytes) {
					t.Errorf("l.%s() = %s; want %s", bunch.label, string(l.itemBytes), string(test.out))
				}
			})
		}

	}
}

func BenchmarkScannerReadString(b *testing.B) {
	for _, bench := range quotedStringCases {
		b.Run(bench.label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := &Scanner{data: []byte(bench.in)}
				_ = l.fetchQuotedString()
			}
		})
	}
}

func BenchmarkScannerReadComment(b *testing.B) {
	for _, bench := range commentCases {
		b.Run(bench.label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := &Scanner{data: []byte(bench.in)}
				_ = l.fetchComment()
			}
		})
	}
}
