package main

import (
	"bytes"
	"fmt"
)

func main() {

}

type FieldValues struct {
	Name   string
	Values []string
}

func MakeTestTable(fieldValuesList ...FieldValues) string {
	buf := new(bytes.Buffer)
	for _, fv := range fieldValuesList {
		for _, v := range fv.Values {
			fmt.Fprintf(buf, "{%s: %s},\n", fv.Name, v)
		}
	}
	return buf.String()
}
