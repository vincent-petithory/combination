// Combination is a tool to generate combinations from a list of grouping data (sets) and prints them to stdout or a file.
//
// It takes the sets, one per line, on stdin or a file and prints the combinations to stdout or a file. The output generated is valid Go syntax, as long as the set name and values are.
//
// A set follows the syntax:
//
//     name: value value ...
//
// The value list is much like a list of arguments in a shell:
// it is space-separated, and "non-safe" strings must be quoted.
//
// For example, the sets:
//
//     card: "Heart Red" Tile Clover "Pike Black"
//     figure: Jack Queen King
//
// would generate the following test table:
//
//     {card: "Heart Red", figure: "Jack"},
//     {card: "Heart Red", figure: "Queen"},
//     {card: "Heart Red", figure: "King"},
//     {card: "Tile", figure: "Jack"},
//     {card: "Tile", figure: "Queen"},
//     {card: "Tile", figure: "King"},
//     {card: "Clover", figure: "Jack"},
//     {card: "Clover", figure: "Queen"},
//     {card: "Clover", figure: "King"},
//     {card: "Pike Black", figure: "Jack"},
//     {card: "Pike Black", figure: "Queen"},
//     {card: "Pike Black", figure: "King"},
//
// In a shell:
//
//     cat << EOF | combination
//     card: "Heart Red" Tile Clover "Pike Black"
//     figure: Jack Queen King
//     EOF
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	srcp  string
	destp string
)

func init() {
	flag.StringVar(&srcp, "sets", "-", "read sets from this file, or stdin if -")
	flag.StringVar(&destp, "o", "-", "write combinations to this file, or stdout if -")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `combination [flags]

  combination is a tool to generate combinations from a list of grouping data (sets)
  It takes the sets, one per line, on stdin or a file and prints the combinations to stdout or a file.
  The output generated is valid Go syntax, as long as the set name and values are.

 A set follows the syntax:

     name: value value ...

 The value list is much like a list of arguments in a shell:
 it is space-separated, and "non-safe" strings must be quoted.

 For example, the sets:

     card: "Heart Red" Tile Clover "Pike Black"
     figure: Jack Queen King

 In a shell:

     cat << EOF | combination
     card: "Heart Red" Tile Clover "Pike Black"
     figure: Jack Queen King
     EOF

`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	var (
		src  io.Reader
		dest io.Writer
	)
	if srcp == "-" {
		src = os.Stdin
	} else {
		f, err := os.Open(srcp)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = f.Close()
		}()
		src = f
	}
	if destp == "-" {
		dest = os.Stdout
	} else {
		f, err := os.Create(destp)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = f.Close()
		}()
		dest = f
	}

	sets, err := parseSets(src)
	if err != nil {
		log.Fatal(err)
	}

	combinations, err := New(sets)
	if err != nil {
		log.Fatal(err)
	}
	if err := WriteCombinations(dest, combinations); err != nil {
		log.Fatal(err)
	}
}

type elementColumn []Element

type Element struct {
	Name  string
	Value string
}

// ErrSetNoValues represents an error when a set contains no values.
var ErrSetNoValues = errors.New("set has no values")

// ErrSetInvalidName represents an error when a set has an empty name or an invalid string value.
var ErrSetInvalidName = errors.New("set has an invalid name")

// ErrValueNoClosingQuote represents an error when a quoted set's value has no closing quote.
var ErrValueNoClosingQuote = errors.New("set value has no closing quote")

// New creates all combinations from sets.
//
// It returns the combinations or an error, ErrSetNoValues if one of the sets provided has no values.
func New(sets []Set) ([]Combination, error) {
	nSets := len(sets)

	// reverse list
	revSets := make([]Set, nSets)
	for i := 0; i < nSets; i++ {
		revSets[nSets-1-i] = sets[i]
	}

	numCombinations := 1
	for _, fv := range sets {
		l := len(fv.Values)
		if l == 0 {
			return nil, ErrSetNoValues
		}
		numCombinations *= l
	}

	var ecs []elementColumn
	numLinesForOneValue := 1
	for _, set := range revSets {
		n := numCombinations / (len(set.Values) * numLinesForOneValue)
		var ec elementColumn
		for k := 0; k < n; k++ {
			for _, val := range set.Values {
				for j := 0; j < numLinesForOneValue; j++ {
					ec = append(ec, Element{Name: set.Name, Value: val})
				}
			}
		}
		ecs = append(ecs, ec)
		numLinesForOneValue *= len(set.Values)
	}
	X := len(ecs)
	Y := len(ecs[0])

	combinations := make([]Combination, 0, numCombinations)
	for y := 0; y < Y; y++ {
		var c Combination
		// read them reverse again
		for x := X - 1; x > -1; x-- {
			c = append(c, ecs[x][y])
		}
		combinations = append(combinations, c)
	}
	return combinations, nil
}

// Combination represents a single combination created from one or more sets.
type Combination []Element

// WriteCombinations writes all combinations to w.
//
// It returns an error if an error occurs when writing to w.
func WriteCombinations(w io.Writer, combinations []Combination) error {
	for _, c := range combinations {
		if _, err := fmt.Fprint(w, "{"); err != nil {
			return err
		}
		for _, e := range c[:len(c)-1] {
			if _, err := fmt.Fprintf(w, "%s: %s, ", e.Name, e.Value); err != nil {
				return err
			}
		}
		e := c[len(c)-1]
		if _, err := fmt.Fprintf(w, "%s: %s", e.Name, e.Value); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "},"); err != nil {
			return err
		}
	}
	return nil
}
