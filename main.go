package main

import (
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

func WriteTestTable(w io.Writer, fieldValuesList ...FieldValues) error {
	nf := len(fieldValuesList)

	var cols []Col
	for i, fv := range fieldValuesList {
		n := 1
		for j := i + 1; j < nf; j++ {
			l := len(fieldValuesList[j].Values)
			if l == 0 {
				continue
			}
			n *= l
		}
		repeatN := 1
		for m := 0; m < i; m++ {
			repeatN *= 2
		}

		col := make(Col, 0, n*len(fv.Values)*repeatN)
		for k := 0; k < repeatN; k++ {
			for _, val := range fv.Values {
				for j := 0; j < n; j++ {
					col = append(col, C{Name: fv.Name, Value: val})
				}
			}
		}
		cols = append(cols, col)
	}
	X := len(cols)
	Y := len(cols[0])

	for y := 0; y < Y; y++ {
		if _, err := fmt.Fprint(w, "{"); err != nil {
			return err
		}
		for x := 0; x < X; x++ {
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
