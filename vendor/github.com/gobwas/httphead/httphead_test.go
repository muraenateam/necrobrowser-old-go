package httphead

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func ExampleScanTokens() {
	var values []string

	ScanTokens([]byte(`a,b,c`), func(v []byte) bool {
		values = append(values, string(v))
		return v[0] != 'b'
	})

	fmt.Println(values)
	// Output: [a b]
}

func ExampleScanOptions() {
	foo := map[string]string{}

	ScanOptions([]byte(`foo;bar=1;baz`), func(index int, key, param, value []byte) Control {
		foo[string(param)] = string(value)
		return ControlContinue
	})

	fmt.Printf("bar:%s baz:%s", foo["bar"], foo["baz"])
	// Output: bar:1 baz:
}

func ExampleParseOptions() {
	options, ok := ParseOptions([]byte(`foo;bar=1,baz`), nil)
	fmt.Println(options, ok)
	// Output: [{foo [bar:1]} {baz []}] true
}

func ExampleParseOptionsLifetime() {
	data := []byte(`foo;bar=1,baz`)
	options, ok := ParseOptions(data, nil)
	copy(data, []byte(`xxx;yyy=0,zzz`))
	fmt.Println(options, ok)
	// Output: [{xxx [yyy:0]} {zzz []}] true
}

var listCases = []struct {
	label string
	in    []byte
	ok    bool
	exp   [][]byte
}{
	{
		label: "simple",
		in:    []byte(`a,b,c`),
		ok:    true,
		exp: [][]byte{
			[]byte(`a`),
			[]byte(`b`),
			[]byte(`c`),
		},
	},
	{
		label: "simple",
		in:    []byte(`a,b,,c`),
		ok:    true,
		exp: [][]byte{
			[]byte(`a`),
			[]byte(`b`),
			[]byte(`c`),
		},
	},
	{
		label: "simple",
		in:    []byte(`a,b;c`),
		ok:    false,
		exp: [][]byte{
			[]byte(`a`),
			[]byte(`b`),
		},
	},
}

func TestScanTokens(t *testing.T) {
	for _, test := range listCases {
		t.Run(test.label, func(t *testing.T) {
			var act [][]byte
			ok := ScanTokens(test.in, func(v []byte) bool {
				act = append(act, v)
				return true
			})
			if ok != test.ok {
				t.Errorf("unexpected result: %v; want %v", ok, test.ok)
			}
			if an, en := len(act), len(test.exp); an != en {
				t.Errorf("unexpected length of result: %d; want %d", an, en)
			} else {
				for i, ev := range test.exp {
					if av := act[i]; !bytes.Equal(av, ev) {
						t.Errorf("unexpected %d-th value: %#q; want %#q", i, string(av), string(ev))
					}
				}
			}
		})
	}
}

func BenchmarkScanTokens(b *testing.B) {
	for _, bench := range listCases {
		b.Run(bench.label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ScanTokens(bench.in, func(v []byte) bool { return true })
			}
		})
	}
}

func randASCII(dst []byte) {
	for i := 0; i < len(dst); i++ {
		dst[i] = byte(rand.Intn('z'-'a')) + 'a'
	}
}

type tuple struct {
	index                    int
	option, attribute, value []byte
}

var parametersCases = []struct {
	label string
	in    []byte
	ok    bool
	exp   []tuple
}{
	{
		label: "simple",
		in:    []byte(`a,b,c`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`a`)},
			{index: 1, option: []byte(`b`)},
			{index: 2, option: []byte(`c`)},
		},
	},
	{
		label: "simple",
		in:    []byte(`a,b,c;foo=1;bar=2`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`a`)},
			{index: 1, option: []byte(`b`)},
			{index: 2, option: []byte(`c`), attribute: []byte(`foo`), value: []byte(`1`)},
			{index: 2, option: []byte(`c`), attribute: []byte(`bar`), value: []byte(`2`)},
		},
	},
	{
		label: "simple",
		in:    []byte(`c;foo;bar=2`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`c`), attribute: []byte(`foo`)},
			{index: 0, option: []byte(`c`), attribute: []byte(`bar`), value: []byte(`2`)},
		},
	},
	{
		label: "simple",
		in:    []byte(`foo;bar=1;baz`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`foo`), attribute: []byte(`bar`), value: []byte(`1`)},
			{index: 0, option: []byte(`foo`), attribute: []byte(`baz`)},
		},
	},
	{
		label: "simple_quoted",
		in:    []byte(`c;bar="2"`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`c`), attribute: []byte(`bar`), value: []byte(`2`)},
		},
	},
	{
		label: "simple_dup",
		in:    []byte(`c;bar=1,c;bar=2`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`c`), attribute: []byte(`bar`), value: []byte(`1`)},
			{index: 1, option: []byte(`c`), attribute: []byte(`bar`), value: []byte(`2`)},
		},
	},
	{
		label: "all",
		in:    []byte(`foo;a=1;b=2;c=3,bar;z,baz`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`foo`), attribute: []byte(`a`), value: []byte(`1`)},
			{index: 0, option: []byte(`foo`), attribute: []byte(`b`), value: []byte(`2`)},
			{index: 0, option: []byte(`foo`), attribute: []byte(`c`), value: []byte(`3`)},
			{index: 1, option: []byte(`bar`), attribute: []byte(`z`)},
			{index: 2, option: []byte(`baz`)},
		},
	},
	{
		label: "comma",
		in:    []byte(`foo;a=1,, , ,bar;b=2`),
		ok:    true,
		exp: []tuple{
			{index: 0, option: []byte(`foo`), attribute: []byte(`a`), value: []byte(`1`)},
			{index: 1, option: []byte(`bar`), attribute: []byte(`b`), value: []byte(`2`)},
		},
	},
}

func TestParameters(t *testing.T) {
	for _, test := range parametersCases {
		t.Run(test.label, func(t *testing.T) {
			var act []tuple

			ok := ScanOptions(test.in, func(index int, key, param, value []byte) Control {
				act = append(act, tuple{index, key, param, value})
				return ControlContinue
			})

			if ok != test.ok {
				t.Errorf("unexpected result: %v; want %v", ok, test.ok)
			}
			if an, en := len(act), len(test.exp); an != en {
				t.Errorf("unexpected length of result: %d; want %d", an, en)
				return
			}

			for i, e := range test.exp {
				a := act[i]

				if a.index != e.index || !bytes.Equal(a.option, e.option) || !bytes.Equal(a.attribute, e.attribute) || !bytes.Equal(a.value, e.value) {
					t.Errorf(
						"unexpected %d-th tuple: #%d %#q[%#q = %#q]; want #%d %#q[%#q = %#q]",
						i,
						a.index, string(a.option), string(a.attribute), string(a.value),
						e.index, string(e.option), string(e.attribute), string(e.value),
					)
				}
			}
		})
	}
}

func BenchmarkParameters(b *testing.B) {
	for _, bench := range parametersCases {
		b.Run(bench.label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ScanOptions(bench.in, func(_ int, _, _, _ []byte) Control { return ControlContinue })
			}
		})
	}
}

var selectOptionsCases = []struct {
	label    string
	selector OptionSelector
	in       []byte
	p        []Option
	exp      []Option
	ok       bool
}{
	{
		label: "simple",
		selector: OptionSelector{
			Flags: SelectCopy | SelectUnique,
		},
		in: []byte(`foo;a=1,foo;a=2`),
		p:  nil,
		exp: []Option{
			NewOption("foo", map[string]string{"a": "1"}),
		},
		ok: true,
	},
	{
		label: "simple",
		selector: OptionSelector{
			Flags: SelectUnique,
		},
		in: []byte(`foo;a=1,foo;a=2`),
		p:  make([]Option, 0, 2),
		exp: []Option{
			NewOption("foo", map[string]string{"a": "1"}),
		},
		ok: true,
	},
	{
		label: "multiparam_stack",
		selector: OptionSelector{
			Flags: SelectUnique,
		},
		in: []byte(`foo;a=1;b=2;c=3;d=4;e=5;f=6;g=7;h=8,bar`),
		p:  make([]Option, 0, 2),
		exp: []Option{
			NewOption("foo", map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
				"e": "5",
				"f": "6",
				"g": "7",
				"h": "8",
			}),
			NewOption("bar", nil),
		},
		ok: true,
	},
	{
		label: "multiparam_stack",
		selector: OptionSelector{
			Flags: SelectCopy | SelectUnique,
		},
		in: []byte(`foo;a=1;b=2;c=3;d=4;e=5;f=6;g=7;h=8,bar`),
		p:  make([]Option, 0, 2),
		exp: []Option{
			NewOption("foo", map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
				"e": "5",
				"f": "6",
				"g": "7",
				"h": "8",
			}),
			NewOption("bar", nil),
		},
		ok: true,
	},
	{
		label: "multiparam_heap",
		selector: OptionSelector{
			Flags: SelectCopy | SelectUnique,
		},
		in: []byte(`foo;a=1;b=2;c=3;d=4;e=5;f=6;g=7;h=8;i=9;j=10,bar`),
		p:  make([]Option, 0, 2),
		exp: []Option{
			NewOption("foo", map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
				"e": "5",
				"f": "6",
				"g": "7",
				"h": "8",
				"i": "9",
				"j": "10",
			}),
			NewOption("bar", nil),
		},
		ok: true,
	},
}

func TestSelectOptions(t *testing.T) {
	for _, test := range selectOptionsCases {
		t.Run(test.label+test.selector.Flags.String(), func(t *testing.T) {
			act, ok := test.selector.Select(test.in, test.p)
			if ok != test.ok {
				t.Errorf("SelectOptions(%q) wellformed sign is %v; want %v", string(test.in), ok, test.ok)
			}
			if !optionsEqual(act, test.exp) {
				t.Errorf("SelectOptions(%q) = %v; want %v", string(test.in), act, test.exp)
			}
		})
	}
}

func BenchmarkSelectOptions(b *testing.B) {
	for _, test := range selectOptionsCases {
		s := test.selector
		b.Run(test.label+s.Flags.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = s.Select(test.in, test.p)
			}
		})
	}
}

func optionsEqual(a, b []Option) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func TestOptionCopy(t *testing.T) {
	for i, test := range []struct {
		pairs int
	}{
		{4},
		{16},
	} {

		name := []byte(fmt.Sprintf("test:%d", i))
		n := make([]byte, len(name))
		copy(n, name)
		opt := Option{Name: n}

		pairs := make([]pair, test.pairs)
		for i := 0; i < len(pairs); i++ {
			pair := pair{make([]byte, 8), make([]byte, 8)}
			randASCII(pair.key)
			randASCII(pair.value)
			pairs[i] = pair

			k, v := make([]byte, len(pair.key)), make([]byte, len(pair.value))
			copy(k, pair.key)
			copy(v, pair.value)

			opt.Parameters.Set(k, v)
		}

		cp := opt.Copy(make([]byte, opt.Size()))

		memset(opt.Name, 'x')
		for _, p := range opt.Parameters.data() {
			memset(p.key, 'x')
			memset(p.value, 'x')
		}

		if !bytes.Equal(cp.Name, name) {
			t.Errorf("name was not copied properly: %q; want %q", string(cp.Name), string(name))
		}
		for i, p := range cp.Parameters.data() {
			exp := pairs[i]
			if !bytes.Equal(p.key, exp.key) || !bytes.Equal(p.value, exp.value) {
				t.Errorf(
					"%d-th pair was not copied properly: %q=%q; want %q=%q",
					i, string(p.key), string(p.value), string(exp.key), string(exp.value),
				)
			}
		}
	}
}

func memset(dst []byte, v byte) {
	copy(dst, bytes.Repeat([]byte{v}, len(dst)))
}
