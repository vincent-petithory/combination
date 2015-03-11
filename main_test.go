package main

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"unicode"
)

type iotest struct {
	Input           []Set
	Output          string
	NumCombinations int
}

func (test *iotest) Run(tb testing.TB) {
	combinations, err := New(test.Input)
	if err != nil {
		log.Fatal(err)
	}
	if n := len(combinations); test.NumCombinations != n {
		tb.Fatalf("expected %d combinations, got %d", test.NumCombinations, n)
	}

	var buf bytes.Buffer
	if err := WriteCombinations(&buf, combinations); err != nil {
		log.Fatal(err)
	}
	output := buf.String()
	if test.Output != output {
		tb.Fatalf("expected:\n%#v\ngot:\n%#v", test.Output, output)
	}
}

func TestSingleField(t *testing.T) {
	test := &iotest{
		Input: []Set{
			{Name: "F1", Values: []string{`"X"`, `"Y"`}},
		},
		Output: `{F1: "X"},
{F1: "Y"},
`,
		NumCombinations: 2,
	}

	test.Run(t)
}

func TestMultiField1(t *testing.T) {
	test := &iotest{
		Input: []Set{
			{Name: "S1", Values: []string{`"X"`, `"Y"`}},
			{Name: "S2", Values: []string{`"µ"`, `"v"`}},
			{Name: "I3", Values: []string{"0xEDEA", "42", "0"}},
		},
		Output: `{S1: "X", S2: "µ", I3: 0xEDEA},
{S1: "X", S2: "µ", I3: 42},
{S1: "X", S2: "µ", I3: 0},
{S1: "X", S2: "v", I3: 0xEDEA},
{S1: "X", S2: "v", I3: 42},
{S1: "X", S2: "v", I3: 0},
{S1: "Y", S2: "µ", I3: 0xEDEA},
{S1: "Y", S2: "µ", I3: 42},
{S1: "Y", S2: "µ", I3: 0},
{S1: "Y", S2: "v", I3: 0xEDEA},
{S1: "Y", S2: "v", I3: 42},
{S1: "Y", S2: "v", I3: 0},
`,
		NumCombinations: 12,
	}
	test.Run(t)
}

func TestMultiField2(t *testing.T) {
	test := &iotest{
		Input: []Set{
			{Name: "S1", Values: []string{`"X"`, `"Y"`, `"Z"`}},
			{Name: "S2", Values: []string{`"µ"`}},
			{Name: "I3", Values: []string{"42", "0"}},
		},
		Output: `{S1: "X", S2: "µ", I3: 42},
{S1: "X", S2: "µ", I3: 0},
{S1: "Y", S2: "µ", I3: 42},
{S1: "Y", S2: "µ", I3: 0},
{S1: "Z", S2: "µ", I3: 42},
{S1: "Z", S2: "µ", I3: 0},
`,
		NumCombinations: 6,
	}
	test.Run(t)
}

func TestBiggerAndImmutableField(t *testing.T) {
	test := &iotest{
		Input: []Set{
			{Name: "name", Values: []string{`""`, "randString()"}},
			{Name: "gauid", Values: []string{"rga()", "ruid()", "buid()", `""`}},
			{Name: "suid", Values: []string{"rs()", "ruid()", "buid()", `""`}},
			{Name: "cuid", Values: []string{"rc()", "ruid()", "buid()", `""`}},
			{Name: "status", Values: []string{"http.StatusBadRequest"}},
		},
		Output: `{name: "", gauid: rga(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: rga(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: ruid(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: buid(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: "", gauid: "", suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: "", suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: "", suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: "", gauid: "", suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: "", gauid: "", suid: "", cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: rga(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: ruid(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: buid(), suid: "", cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: rs(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: rs(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: rs(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: rs(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: ruid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: ruid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: ruid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: ruid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: buid(), cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: buid(), cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: buid(), cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: buid(), cuid: "", status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: "", cuid: rc(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: "", cuid: ruid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: "", cuid: buid(), status: http.StatusBadRequest},
{name: randString(), gauid: "", suid: "", cuid: "", status: http.StatusBadRequest},
`,
		NumCombinations: 128,
	}
	test.Run(t)
}

func TestErrNoValues(t *testing.T) {
	combinations, err := New([]Set{
		{Name: "x", Values: []string{"0", "1", "2"}},
		{Name: "y", Values: []string{}},
		{Name: "z", Values: []string{"3", "4", "5"}},
	})
	if err == nil {
		t.Error("expected a non-nil error, got nil")
	}
	if combinations != nil {
		t.Errorf("expected a nil combinations, got %#v", combinations)
	}
}

// TestMarshalUnmarshalSet tests a roundtrip {enc,dec}oding of a Set.
func TestMarshalUnmarshalSet(t *testing.T) {
	if err := quick.Check(func(name string, values []string) bool {
		if name == "" {
			set := Set{Name: name, Values: values}
			_, err := set.MarshalText()
			if err == nil {
				t.Error("expected non-nil error, got nil")
				return false
			}
			return true
		}
		// cleanup name
		var buf bytes.Buffer
		for _, ruune := range name {
			// Go identifier
			if unicode.IsLetter(ruune) || unicode.IsDigit(ruune) {
				buf.WriteRune(ruune)
			}
		}

		iset := Set{Name: buf.String(), Values: values}
		if iset.Name == "" {
			// name had no valid rune, just put one
			iset.Name = "a"
		}
		b, err := iset.MarshalText()
		if err != nil {
			t.Error(err)
			return false
		}
		var oset Set
		if err := oset.UnmarshalText(b); err != nil {
			t.Error(err)
			return false
		}
		if !reflect.DeepEqual(iset, oset) {
			t.Errorf("expected %#v, got %#v", iset, oset)
			return false
		}
		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

func TestParseSets(t *testing.T) {
	input := `card: ["\"Heart\"", "\"Tile\"", "\"Clover\"", "\"Pike\""]
figure: ["\"Jack\"", "\"Queen\"", "\"King\""]`
	expectedSets := []Set{
		{Name: "card", Values: []string{"\"Heart\"", "\"Tile\"", "\"Clover\"", "\"Pike\""}},
		{Name: "figure", Values: []string{"\"Jack\"", "\"Queen\"", "\"King\""}},
	}
	sets, err := parseSets(strings.NewReader(input))
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(sets, expectedSets) {
		t.Errorf("expected %#v, got %#v", expectedSets, sets)
	}
}
