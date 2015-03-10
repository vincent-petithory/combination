package main

import (
	"bytes"
	"testing"
)

type iotest struct {
	Input  []FieldValues
	Output string
}

func (test *iotest) Run(tb testing.TB) {
	var buf bytes.Buffer
	if err := WriteTestTable(&buf, test.Input...); err != nil {
		tb.Fatal(err)
	}
	output := buf.String()
	if test.Output != output {
		tb.Fatalf("expected:\n%#v\ngot:\n%#v", test.Output, output)
	}
}

func TestSingleField(t *testing.T) {
	test := &iotest{
		Input: []FieldValues{
			{Name: "F1", Values: []string{`"X"`, `"Y"`}},
		},
		Output: `{F1: "X", },
{F1: "Y", },
`,
	}

	test.Run(t)
}

func TestMultiField1(t *testing.T) {
	test := &iotest{
		Input: []FieldValues{
			{Name: "S1", Values: []string{`"X"`, `"Y"`}},
			{Name: "S2", Values: []string{`"µ"`, `"v"`}},
			{Name: "I3", Values: []string{"0xEDEA", "42", "0"}},
		},
		Output: `{S1: "X", S2: "µ", I3: 0xEDEA, },
{S1: "X", S2: "µ", I3: 42, },
{S1: "X", S2: "µ", I3: 0, },
{S1: "X", S2: "v", I3: 0xEDEA, },
{S1: "X", S2: "v", I3: 42, },
{S1: "X", S2: "v", I3: 0, },
{S1: "Y", S2: "µ", I3: 0xEDEA, },
{S1: "Y", S2: "µ", I3: 42, },
{S1: "Y", S2: "µ", I3: 0, },
{S1: "Y", S2: "v", I3: 0xEDEA, },
{S1: "Y", S2: "v", I3: 42, },
{S1: "Y", S2: "v", I3: 0, },
`,
	}
	test.Run(t)
}

func TestMultiField2(t *testing.T) {
	test := &iotest{
		Input: []FieldValues{
			{Name: "S1", Values: []string{`"X"`, `"Y"`, `"Z"`}},
			{Name: "S2", Values: []string{`"µ"`}},
			{Name: "I3", Values: []string{"42", "0"}},
		},
		Output: `{S1: "X", S2: "µ", I3: 42, },
{S1: "X", S2: "µ", I3: 0, },
{S1: "Y", S2: "µ", I3: 42, },
{S1: "Y", S2: "µ", I3: 0, },
{S1: "Z", S2: "µ", I3: 42, },
{S1: "Z", S2: "µ", I3: 0, },
`,
	}
	test.Run(t)
}

func TestBiggerAndImmutableField(t *testing.T) {
	test := &iotest{
		Input: []FieldValues{
			{Name: "name", Values: []string{`""`, "randString()"}},
			{Name: "gauid", Values: []string{"rga()", "ruid()", "buid()", `""`}},
			{Name: "suid", Values: []string{"rs()", "ruid()", "buid()", `""`}},
			{Name: "cuid", Values: []string{"rc()", "ruid()", "buid()", `""`}},
			{Name: "status", Values: []string{"http.StatusBadRequest"}},
		},
		Output: `{name: "", gauid: rga(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: rga(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: ruid(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: buid(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: "", suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: "", suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: "", suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: "", gauid: "", suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: "", gauid: "", suid: "", cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: rga(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: ruid(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: buid(), suid: "", cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: rs(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: rs(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: rs(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: rs(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: ruid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: ruid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: ruid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: ruid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: buid(), cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: buid(), cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: buid(), cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: buid(), cuid: "", status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: "", cuid: rc(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: "", cuid: ruid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: "", cuid: buid(), status: http.StatusBadRequest, },
{name: randString(), gauid: "", suid: "", cuid: "", status: http.StatusBadRequest, },
`,
	}
	test.Run(t)
}
