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

func TestMultiField(t *testing.T) {
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
