package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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

  combination generates combinations for SETS and prints them to FILE.

  One SET follows the syntax:

    set    := name ": [" value *(value) "]"
    name   := string
    value  := "\"" string "\""
              ; string must be a valid JSON string.

  Examples:

    card: ["\"Hearts\"", "\"Tiles\"", "\"Clovers\"", "\"Pikes\""]
    bool: ["true", "false"]
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

func parseSets(r io.Reader) ([]Set, error) {
	var sets []Set
	bufsrc := bufio.NewScanner(r)
	bufsrc.Split(bufio.ScanLines)
	for bufsrc.Scan() {
		var set Set
		if err := set.UnmarshalText([]byte(bufsrc.Text())); err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	if err := bufsrc.Err(); err != nil {
		return nil, err
	}
	return sets, nil
}

// Set represents a named grouping of values (a set).
type Set struct {
	Name   string
	Values []string
}

// MarshalText implements TextMarshaler.
func (s Set) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := fmt.Fprint(buf, s.Name, ": "); err != nil {
		return nil, err
	}
	if err := json.NewEncoder(buf).Encode(s.Values); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalText implements TextUnmarshaler.
func (s *Set) UnmarshalText(text []byte) error {
	r := bufio.NewReader(bytes.NewReader(text))
	name, err := r.ReadString(':')
	if err != nil {
		return err
	}
	s.Name = name[:len(name)-1]
	if s.Name == "" {
		return ErrSetInvalidName
	}
	runne, _, err := r.ReadRune()
	if err != nil {
		return err
	}
	if runne != ' ' {
		if err := r.UnreadRune(); err != nil {
			return err
		}
	}
	return json.NewDecoder(r).Decode(&s.Values)
}

type elementColumn []Element

type Element struct {
	Name  string
	Value string
}

// ErrSetNoValues represents an error when a set contains no values.
var ErrSetNoValues = errors.New("set has no values")

// ErrSetInvalidName represents an erro when a set has an empty name or an invalid string value.
var ErrSetInvalidName = errors.New("set has an invalid name")

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
		for _, e := range c {
			if _, err := fmt.Fprintf(w, "%s: %s, ", e.Name, e.Value); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "},"); err != nil {
			return err
		}
	}
	return nil
}
