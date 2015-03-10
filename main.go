package main

import (
	"errors"
	"fmt"
	"io"
)

func main() {

}

type FieldValues struct {
	Name   string
	Values []string
}

type Col []C

type C struct {
	Name  string
	Value string
}

var ErrZeroValues = errors.New("field has zero values in its set")

func WriteTestTable(w io.Writer, fieldValuesList ...FieldValues) error {
	nf := len(fieldValuesList)

	// reverse list
	revList := make([]FieldValues, nf)
	for i := 0; i < nf; i++ {
		revList[nf-1-i] = fieldValuesList[i]
	}

	numCombinations := 1
	for _, fv := range fieldValuesList {
		l := len(fv.Values)
		if l == 0 {
			return ErrZeroValues
		}
		numCombinations *= l
	}

	var cols []Col
	numLinesForOneValue := 1
	for _, fv := range revList {
		n := numCombinations / (len(fv.Values) * numLinesForOneValue)
		var col Col
		for k := 0; k < n; k++ {
			for _, val := range fv.Values {
				for j := 0; j < numLinesForOneValue; j++ {
					col = append(col, C{Name: fv.Name, Value: val})
				}
			}
		}
		cols = append(cols, col)
		numLinesForOneValue *= len(fv.Values)
	}
	X := len(cols)
	Y := len(cols[0])

	for y := 0; y < Y; y++ {
		if _, err := fmt.Fprint(w, "{"); err != nil {
			return err
		}
		// read them reverse again
		for x := X - 1; x > -1; x-- {
			if _, err := fmt.Fprintf(w, "%s: %s, ", cols[x][y].Name, cols[x][y].Value); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "},"); err != nil {
			return err
		}
	}
	return nil
}
