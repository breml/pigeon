package json

import "sort"

// ErrorLister is the public interface to access the inner errors
// included in a errList
type ErrorLister interface {
	Errors() []error
}

func (e errList) Errors() []error {
	return e
}

// ParserError is the public interface to errors of type parserError
type ParserError interface {
	Error() string
	InnerError() error
	Pos() (int, int, int)
	Expected() []string
	ExpectedRules() []*rule
}

func (p *parserError) InnerError() error {
	return p.Inner
}

func (p *parserError) Pos() (line, col, offset int) {
	return p.pos.line, p.pos.col, p.pos.offset
}

func (p *parserError) Expected() []string {
	expected := make([]string, 0, len(p.expected))
	var eof *rule
	var ok bool
	if eof, ok = p.expected["!."]; ok {
		delete(p.expected, "!.")
	}
	for k := range p.expected {
		expected = append(expected, k)
	}
	sort.Strings(expected)
	if eof != nil {
		expected = append(expected, "EOF")
	}
	return expected
}

func (p *parserError) ExpectedRules() []*rule {
	expected := make([]*rule, 0, len(p.expected))
	// TODO: not finished
	return expected
}
