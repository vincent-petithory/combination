package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

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
	var buf bytes.Buffer
	if s.Name == "" {
		return nil, ErrSetInvalidName
	}
	if _, err := io.WriteString(&buf, s.Name+": "); err != nil {
		return nil, err
	}
	for _, val := range s.Values {
		if _, err := io.WriteString(&buf, strconv.Quote(val)+" "); err != nil {
			return nil, err
		}
	}
	b := buf.Bytes()
	return b[:len(b)-1], nil
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
	defer func() {
		if s.Values == nil {
			s.Values = []string{}
		}
	}()
	runne, _, err := r.ReadRune()
	if err == io.EOF {
		// We don't have any values
		return nil
	}
	if err != nil {
		return err
	}
	if runne != ' ' {
		if err := r.UnreadRune(); err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(SetValuesSplitFn)
	for scanner.Scan() {
		s.Values = append(s.Values, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return err
}

// SetValuesSplitFn is a scanner func to split values of a set.
var SetValuesSplitFn = bufio.SplitFunc(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// eat all leading spaces or "," if any
	for _, b := range data {
		if b != ' ' {
			break
		}
		advance++
	}
	start := advance
	if len(data) > start && data[start] == '"' {
		var ignoreNextQuote bool
		closingQuoteIdx := -1
		for i, b := range data[start+1:] {
			if b == '\\' && !ignoreNextQuote {
				ignoreNextQuote = true
				continue
			}
			if b == '\\' && ignoreNextQuote {
				ignoreNextQuote = false
				continue
			}
			if b != '"' {
				ignoreNextQuote = false
				continue
			}
			if ignoreNextQuote {
				ignoreNextQuote = false
				continue
			}
			closingQuoteIdx = i
			break
		}
		if closingQuoteIdx == -1 {
			if atEOF {
				err = ErrValueNoClosingQuote
				return
			}
			// ask more data to get the closing quote
			return 0, nil, nil
		}
		advance += len(data[start : (start+1)+(closingQuoteIdx+1)])
		var s string
		s, err = strconv.Unquote(string(data[start : (start+1)+(closingQuoteIdx+1)]))
		token = []byte(s)
		return
	}
	// no opening quote. let's read until ' ' or EOF
	spaceIdx := bytes.IndexByte(data[start:], ' ')
	if spaceIdx == -1 && !atEOF {
		// Not at EOF and no comma, so we have an incomplete token
		return 0, nil, nil
	} else if spaceIdx == -1 && atEOF {
		// that's our complete last token
		advance += len(data[start:])
		token = data[start:]
		return
	}
	advance += len(data[start : start+spaceIdx])
	token = data[start : start+spaceIdx]
	return
})
