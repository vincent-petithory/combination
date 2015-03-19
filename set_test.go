package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"unicode"
)

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
	tests := []struct {
		input string
		sets  []Set
	}{
		{
			input: `card: Heart Tile Clover Pike
figure: Jack Queen King`,
			sets: []Set{
				{Name: "card", Values: []string{"Heart", "Tile", "Clover", "Pike"}},
				{Name: "figure", Values: []string{"Jack", "Queen", "King"}},
			},
		},
		{
			input: `card: "Heart Rouge" "Tile Rouge" "Clover Noir" "Pike Noir"
figure: Jack Queen King`,
			sets: []Set{
				{Name: "card", Values: []string{"Heart Rouge", "Tile Rouge", "Clover Noir", "Pike Noir"}},
				{Name: "figure", Values: []string{"Jack", "Queen", "King"}},
			},
		},
		{
			input: `card: "\"Heart Rouge\"" "Tile Rouge" "Clover Noir" "Pike Noir"
figure: Jack Queen King`,
			sets: []Set{
				{Name: "card", Values: []string{"\"Heart Rouge\"", "Tile Rouge", "Clover Noir", "Pike Noir"}},
				{Name: "figure", Values: []string{"Jack", "Queen", "King"}},
			},
		},
		{
			input: `card: "\"Heart\nRouge\"" "Tile Rouge" "Clover\nNoir" "Pike Noir"
figure: Jack Queen King`,
			sets: []Set{
				{Name: "card", Values: []string{"\"Heart\nRouge\"", "Tile Rouge", "Clover\nNoir", "Pike Noir"}},
				{Name: "figure", Values: []string{"Jack", "Queen", "King"}},
			},
		},
		{
			input: `card:
figure: Jack Queen King`,
			sets: []Set{
				{Name: "card", Values: []string{}},
				{Name: "figure", Values: []string{"Jack", "Queen", "King"}},
			},
		},
	}
	for _, test := range tests {
		sets, err := parseSets(strings.NewReader(test.input))
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(sets, test.sets) {
			t.Errorf("expected:\n%#v\ngot:\n%#v", test.sets, sets)
		}
	}
}
