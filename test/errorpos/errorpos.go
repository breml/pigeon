package errorpos

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 5, col: 1, offset: 22},
			expr: &seqExpr{
				pos: position{line: 5, col: 9, offset: 30},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 5, col: 9, offset: 30},
						name: "_",
					},
					&choiceExpr{
						pos: position{line: 5, col: 12, offset: 33},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 5, col: 12, offset: 33},
								name: "case01",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 21, offset: 42},
								name: "case02",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 30, offset: 51},
								name: "case03",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 39, offset: 60},
								name: "case04",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 48, offset: 69},
								name: "case05",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 57, offset: 78},
								name: "case06",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 66, offset: 87},
								name: "case07",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 75, offset: 96},
								name: "case08",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 84, offset: 105},
								name: "case09",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 93, offset: 114},
								name: "case10",
							},
							&ruleRefExpr{
								pos:  position{line: 5, col: 102, offset: 123},
								name: "case11",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 5, col: 110, offset: 131},
						name: "EOF",
					},
				},
			},
		},
		{
			name: "case01",
			pos:  position{line: 7, col: 1, offset: 136},
			expr: &seqExpr{
				pos: position{line: 7, col: 10, offset: 145},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 7, col: 10, offset: 145},
						val:        "case01",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 7, col: 19, offset: 154},
						name: "_",
					},
					&oneOrMoreExpr{
						pos: position{line: 7, col: 21, offset: 156},
						expr: &seqExpr{
							pos: position{line: 7, col: 22, offset: 157},
							exprs: []interface{}{
								&choiceExpr{
									pos: position{line: 7, col: 23, offset: 158},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 7, col: 23, offset: 158},
											name: "increment",
										},
										&ruleRefExpr{
											pos:  position{line: 7, col: 35, offset: 170},
											name: "decrement",
										},
										&ruleRefExpr{
											pos:  position{line: 7, col: 47, offset: 182},
											name: "zero",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 7, col: 53, offset: 188},
									name: "_",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "case02",
			pos:  position{line: 8, col: 1, offset: 193},
			expr: &seqExpr{
				pos: position{line: 8, col: 10, offset: 202},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 8, col: 10, offset: 202},
						val:        "case02",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 8, col: 19, offset: 211},
						name: "__",
					},
					&oneOrMoreExpr{
						pos: position{line: 8, col: 22, offset: 214},
						expr: &charClassMatcher{
							pos:        position{line: 8, col: 22, offset: 214},
							val:        "[^abc]",
							chars:      []rune{'a', 'b', 'c'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "case03",
			pos:  position{line: 9, col: 1, offset: 222},
			expr: &seqExpr{
				pos: position{line: 9, col: 10, offset: 231},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 9, col: 10, offset: 231},
						val:        "case03",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 9, col: 19, offset: 240},
						name: "__",
					},
					&zeroOrOneExpr{
						pos: position{line: 9, col: 22, offset: 243},
						expr: &litMatcher{
							pos:        position{line: 9, col: 22, offset: 243},
							val:        "x",
							ignoreCase: false,
						},
					},
					&charClassMatcher{
						pos:        position{line: 9, col: 27, offset: 248},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "case04",
			pos:  position{line: 10, col: 1, offset: 254},
			expr: &seqExpr{
				pos: position{line: 10, col: 10, offset: 263},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 10, col: 10, offset: 263},
						val:        "case04",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 10, col: 19, offset: 272},
						name: "__",
					},
					&charClassMatcher{
						pos:        position{line: 10, col: 22, offset: 275},
						val:        "[\\x30-\\x39]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 10, col: 34, offset: 287},
						val:        "[^\\x30-\\x39]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   true,
					},
					&charClassMatcher{
						pos:        position{line: 10, col: 47, offset: 300},
						val:        "[\\pN]",
						classes:    []*unicode.RangeTable{rangeTable("N")},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 10, col: 53, offset: 306},
						val:        "[^\\pN]",
						classes:    []*unicode.RangeTable{rangeTable("N")},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "case05",
			pos:  position{line: 11, col: 1, offset: 313},
			expr: &seqExpr{
				pos: position{line: 11, col: 10, offset: 322},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 11, col: 10, offset: 322},
						val:        "case05",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 11, col: 19, offset: 331},
						name: "__",
					},
					&notExpr{
						pos: position{line: 11, col: 22, offset: 334},
						expr: &litMatcher{
							pos:        position{line: 11, col: 23, offset: 335},
							val:        "not",
							ignoreCase: false,
						},
					},
					&litMatcher{
						pos:        position{line: 11, col: 29, offset: 341},
						val:        "yes",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "case06",
			pos:  position{line: 12, col: 1, offset: 347},
			expr: &seqExpr{
				pos: position{line: 12, col: 10, offset: 356},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 12, col: 10, offset: 356},
						val:        "case06",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 12, col: 19, offset: 365},
						name: "__",
					},
					&notExpr{
						pos: position{line: 12, col: 22, offset: 368},
						expr: &charClassMatcher{
							pos:        position{line: 12, col: 23, offset: 369},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&litMatcher{
						pos:        position{line: 12, col: 29, offset: 375},
						val:        "x",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "case07",
			pos:  position{line: 13, col: 1, offset: 379},
			expr: &seqExpr{
				pos: position{line: 13, col: 10, offset: 388},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 13, col: 10, offset: 388},
						val:        "case07",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 13, col: 19, offset: 397},
						name: "__",
					},
					&notExpr{
						pos: position{line: 13, col: 22, offset: 400},
						expr: &choiceExpr{
							pos: position{line: 13, col: 24, offset: 402},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 13, col: 24, offset: 402},
									val:        "abc",
									ignoreCase: true,
								},
								&charClassMatcher{
									pos:        position{line: 13, col: 33, offset: 411},
									val:        "[a-c]i",
									ranges:     []rune{'a', 'c'},
									ignoreCase: true,
									inverted:   false,
								},
								&charClassMatcher{
									pos:        position{line: 13, col: 42, offset: 420},
									val:        "[\\pL]",
									classes:    []*unicode.RangeTable{rangeTable("L")},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&charClassMatcher{
						pos:        position{line: 13, col: 49, offset: 427},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "case08",
			pos:  position{line: 14, col: 1, offset: 433},
			expr: &seqExpr{
				pos: position{line: 14, col: 10, offset: 442},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 14, col: 10, offset: 442},
						val:        "case08",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 14, col: 19, offset: 451},
						name: "__",
					},
					&andExpr{
						pos: position{line: 14, col: 22, offset: 454},
						expr: &litMatcher{
							pos:        position{line: 14, col: 23, offset: 455},
							val:        "a",
							ignoreCase: true,
						},
					},
					&anyMatcher{
						line: 14, col: 28, offset: 460,
					},
				},
			},
		},
		{
			name: "case09",
			pos:  position{line: 15, col: 1, offset: 462},
			expr: &seqExpr{
				pos: position{line: 15, col: 10, offset: 471},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 15, col: 10, offset: 471},
						val:        "case09",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 15, col: 19, offset: 480},
						name: "__",
					},
					&andExpr{
						pos: position{line: 15, col: 22, offset: 483},
						expr: &charClassMatcher{
							pos:        position{line: 15, col: 23, offset: 484},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&anyMatcher{
						line: 15, col: 29, offset: 490,
					},
				},
			},
		},
		{
			name: "case10",
			pos:  position{line: 16, col: 1, offset: 492},
			expr: &seqExpr{
				pos: position{line: 16, col: 10, offset: 501},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 16, col: 10, offset: 501},
						val:        "case10",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 16, col: 19, offset: 510},
						name: "__",
					},
					&andExpr{
						pos: position{line: 16, col: 22, offset: 513},
						expr: &choiceExpr{
							pos: position{line: 16, col: 24, offset: 515},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 16, col: 24, offset: 515},
									val:        "0",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 16, col: 30, offset: 521},
									val:        "[012]",
									chars:      []rune{'0', '1', '2'},
									ignoreCase: false,
									inverted:   false,
								},
								&charClassMatcher{
									pos:        position{line: 16, col: 38, offset: 529},
									val:        "[3-9]",
									ranges:     []rune{'3', '9'},
									ignoreCase: false,
									inverted:   false,
								},
								&charClassMatcher{
									pos:        position{line: 16, col: 46, offset: 537},
									val:        "[\\pN]",
									classes:    []*unicode.RangeTable{rangeTable("N")},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
					&anyMatcher{
						line: 16, col: 53, offset: 544,
					},
				},
			},
		},
		{
			name: "case11",
			pos:  position{line: 17, col: 1, offset: 546},
			expr: &seqExpr{
				pos: position{line: 17, col: 10, offset: 555},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 17, col: 10, offset: 555},
						val:        "case11",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 17, col: 19, offset: 564},
						name: "__",
					},
					&notExpr{
						pos: position{line: 17, col: 22, offset: 567},
						expr: &notExpr{
							pos: position{line: 17, col: 24, offset: 569},
							expr: &litMatcher{
								pos:        position{line: 17, col: 26, offset: 571},
								val:        "a",
								ignoreCase: false,
							},
						},
					},
					&litMatcher{
						pos:        position{line: 17, col: 32, offset: 577},
						val:        "a",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "increment",
			pos:  position{line: 19, col: 1, offset: 582},
			expr: &litMatcher{
				pos:        position{line: 19, col: 13, offset: 594},
				val:        "inc",
				ignoreCase: false,
			},
		},
		{
			name: "decrement",
			pos:  position{line: 20, col: 1, offset: 600},
			expr: &litMatcher{
				pos:        position{line: 20, col: 13, offset: 612},
				val:        "dec",
				ignoreCase: false,
			},
		},
		{
			name: "zero",
			pos:  position{line: 21, col: 1, offset: 618},
			expr: &litMatcher{
				pos:        position{line: 21, col: 8, offset: 625},
				val:        "zero",
				ignoreCase: false,
			},
		},
		{
			name: "oneOrMore",
			pos:  position{line: 22, col: 1, offset: 632},
			expr: &litMatcher{
				pos:        position{line: 22, col: 13, offset: 644},
				val:        "oneOrMore",
				ignoreCase: false,
			},
		},
		{
			name: "_",
			pos:  position{line: 23, col: 1, offset: 656},
			expr: &zeroOrMoreExpr{
				pos: position{line: 23, col: 5, offset: 660},
				expr: &charClassMatcher{
					pos:        position{line: 23, col: 5, offset: 660},
					val:        "[ \\t\\n\\r]",
					chars:      []rune{' ', '\t', '\n', '\r'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 24, col: 1, offset: 671},
			expr: &charClassMatcher{
				pos:        position{line: 24, col: 6, offset: 676},
				val:        "[ ]",
				chars:      []rune{' '},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 25, col: 1, offset: 680},
			expr: &notExpr{
				pos: position{line: 25, col: 7, offset: 686},
				expr: &anyMatcher{
					line: 25, col: 8, offset: 687,
				},
			},
		},
	},
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
}

type current struct {
	pos  position // start position of the match
	text []rune   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename:        filename,
		errs:            new(errList),
		data:            []rune(string(b)),
		pt:              savepoint{position: position{offset: -1, line: 1}},
		recover:         true,
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make(map[string]struct{}),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []rune
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       map[string]struct{}
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = make(map[string]struct{})
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected[want] = struct{}{}
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset++
	if p.pt.offset >= len(p.data) {
		p.pt.rn = utf8.RuneError
		return
	}
	rn := p.data[p.pt.offset]
	p.pt.rn = rn
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		p.addErr(errInvalidEncoding)
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []rune {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			expected := make([]string, 0, len(p.maxFailExpected))
			eof := false
			if _, ok := p.maxFailExpected["!."]; ok {
				delete(p.maxFailExpected, "!.")
				eof = true
			}
			for k := range p.maxFailExpected {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt
	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
