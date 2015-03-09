package main

import "testing"

func equals(tb testing.TB, expected string, actual string) {
	if expected != actual {
		tb.Errorf("expected:\n%#v\ngot:\n%#v", expected, actual)
	}
}

func TestSingleField(t *testing.T) {
	fieldValuesList := []FieldValues{
		{Name: "F1", Values: []string{`"X"`, `"Y"`}},
	}
	expectedOutput := `{F1: "X"},
{F1: "Y"},
`

	output := MakeTestTable(fieldValuesList...)
	equals(t, expectedOutput, output)
}
