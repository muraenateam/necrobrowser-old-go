package httphead

import (
	"bytes"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	for _, test := range []struct {
		name string
		line string
		exp  RequestLine
		fail bool
	}{
		{
			line: "",
			fail: true,
		},
		{
			line: "GET",
			fail: true,
		},
		{
			line: "GET ",
			fail: true,
		},
		{
			line: "GET  ",
			fail: true,
		},
		{
			line: "GET   ",
			fail: true,
		},
		{
			line: "GET / HTTP/1.1",
			exp: RequestLine{
				Method:  []byte("GET"),
				URI:     []byte("/"),
				Version: Version{1, 1},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, ok := ParseRequestLine([]byte(test.line))
			if test.fail && ok {
				t.Fatalf("unexpected successful parsing")
			}
			if !test.fail && !ok {
				t.Fatalf("unexpected parse error")
			}
			if test.fail {
				return
			}
			if act, exp := r.Method, test.exp.Method; !bytes.Equal(act, exp) {
				t.Errorf("unexpected parsed method: %q; want %q", act, exp)
			}
			if act, exp := r.URI, test.exp.URI; !bytes.Equal(act, exp) {
				t.Errorf("unexpected parsed uri: %q; want %q", act, exp)
			}
			if act, exp := r.Version, test.exp.Version; act != exp {
				t.Errorf("unexpected parsed version: %+v; want %+v", act, exp)
			}
		})
	}
}

func TestParseResponseLine(t *testing.T) {
	for _, test := range []struct {
		name string
		line string
		exp  ResponseLine
		fail bool
	}{
		{
			line: "",
			fail: true,
		},
		{
			line: "HTTP/1.1",
			fail: true,
		},
		{
			line: "HTTP/1.1 ",
			fail: true,
		},
		{
			line: "HTTP/1.1  ",
			fail: true,
		},
		{
			line: "HTTP/1.1   ",
			fail: true,
		},
		{
			line: "HTTP/1.1 200 OK",
			exp: ResponseLine{
				Version: Version{1, 1},
				Status:  200,
				Reason:  []byte("OK"),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, ok := ParseResponseLine([]byte(test.line))
			if test.fail && ok {
				t.Fatalf("unexpected successful parsing")
			}
			if !test.fail && !ok {
				t.Fatalf("unexpected parse error")
			}
			if test.fail {
				return
			}
			if act, exp := r.Version, test.exp.Version; act != exp {
				t.Errorf("unexpected parsed version: %+v; want %+v", act, exp)
			}
			if act, exp := r.Status, test.exp.Status; act != exp {
				t.Errorf("unexpected parsed status: %d; want %d", act, exp)
			}
			if act, exp := r.Reason, test.exp.Reason; !bytes.Equal(act, exp) {
				t.Errorf("unexpected parsed reason: %q; want %q", act, exp)
			}
		})
	}
}

var versionCases = []struct {
	in    []byte
	major int
	minor int
	ok    bool
}{
	{[]byte("HTTP/1.1"), 1, 1, true},
	{[]byte("HTTP/1.0"), 1, 0, true},
	{[]byte("HTTP/1.2"), 1, 2, true},
	{[]byte("HTTP/42.1092"), 42, 1092, true},
}

func TestParseHttpVersion(t *testing.T) {
	for _, c := range versionCases {
		t.Run(string(c.in), func(t *testing.T) {
			major, minor, ok := ParseVersion(c.in)
			if major != c.major || minor != c.minor || ok != c.ok {
				t.Errorf(
					"parseHttpVersion([]byte(%q)) = %v, %v, %v; want %v, %v, %v",
					string(c.in), major, minor, ok, c.major, c.minor, c.ok,
				)
			}
		})
	}
}

func BenchmarkParseHttpVersion(b *testing.B) {
	for _, c := range versionCases {
		b.Run(string(c.in), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = ParseVersion(c.in)
			}
		})
	}
}
