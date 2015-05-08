package vm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ϡtheProgram is the variable that holds the program generated by the
// builder for the input PEG.
var ϡtheProgram *ϡprogram

//+pigeon: vm.go

// ϡsentinel is a type used to define sentinel values that shouldn't
// be equal to something else.
type ϡsentinel int

const (
	// ϡmatchFailed is a sentinel value used to indicate a match failure.
	ϡmatchFailed ϡsentinel = iota - 1
)

const (
	// stack IDs, used in PUSH and POP's first argument
	ϡpstackID = iota + 1
	ϡlstackID
	ϡvstackID
	ϡistackID
	ϡastackID

	// special V stack values
	ϡvValNil    uint16 = 0
	ϡvValFailed uint16 = 1
	ϡvValEmpty  uint16 = 2
)

var (
	ϡstackNm = []string{
		ϡpstackID: "P",
		ϡlstackID: "L",
		ϡvstackID: "V",
		ϡistackID: "I",
		ϡastackID: "A",
	}
)

// special values that may be pushed on the V stack.
var ϡvSpecialValues = []interface{}{
	nil,
	ϡmatchFailed,
	[]interface{}(nil),
}

// ϡmemoizedResult holds the state required to reuse a memoized result.
type ϡmemoizedResult struct {
	v  interface{}
	pt ϡsvpt
}

// ϡffp holds state to record the farthest failure position.
type ϡffp struct {
	pos      position
	ruleNmIx int
	got      []byte
	want     string
}

// err returns an error if a farthest failure position has been recorded,
// nil otherwise.
func (ffp ϡffp) err() error {
	if ffp.pos.offset >= 0 {
		return ffp
	}
	return nil
}

// Error implements the error interface.
func (ffp ϡffp) Error() string {
	return fmt.Sprintf("expected %s, got %s", ffp.want, strconv.Quote(string(ffp.got)))
}

// ϡprogram is the data structure that is generated by the builder
// based on an input PEG. It contains the program information required
// to execute the grammar using the vm.
type ϡprogram struct {
	instrs []ϡinstr

	// lists
	ms []ϡmatcher
	as []func(*ϡvm) (interface{}, error)
	bs []func(*ϡvm) (bool, error)
	ss []string
}

// String formats the program's instructions in a human-readable format.
func (pg ϡprogram) String() string {
	var buf bytes.Buffer

	for i, instr := range pg.instrs {
		buf.WriteString(fmt.Sprintf("[%3d]: %s\n", i, pg.instrToString(instr)))
	}
	return buf.String()
}

// instrToString formats an instruction in a human-readable format, in the
// context of the program.
func (pg ϡprogram) instrToString(instr ϡinstr) string {
	var buf bytes.Buffer

	rule := pg.ruleNameAt(instr.ruleNmIx)
	if rule == "" {
		rule = "<bootstrap>"
	}
	buf.WriteString(fmt.Sprintf("%s.%s %v", rule, instr.op, instr.args))
	switch instr.op {
	case ϡopPush, ϡopPop:
		buf.WriteString(" " + ϡstackNm[instr.args[0]])
	case ϡopMatch:
		buf.WriteString(fmt.Sprintf(" %s", pg.ms[instr.args[0]]))
	case ϡopStoreIfT:
		buf.WriteString(" " + pg.ss[instr.args[0]])
	}
	return buf.String()
}

// ruleNameAt returns the name of the rule that contains the instruction
// index. It returns an empty string is the instruction is not part of a
// rule (bootstrap instruction, invalid index).
func (pg ϡprogram) ruleNameAt(ix int) string {
	if ix < 0 || ix >= len(pg.ss) {
		return ""
	}
	return pg.ss[ix]
}

// ϡvm holds the state to execute a compiled grammar.
type ϡvm struct {
	// input
	filename string
	parser   *ϡparser

	// options
	debug   bool
	memoize bool
	recover bool

	// program data
	pc  uint16
	pg  *ϡprogram
	cur current

	// stacks
	p *ϡpstack
	l *ϡlstack
	v *ϡvstack
	i *ϡistack
	a *ϡastack

	// memoized results: by instruction index, then by byte offset
	memo map[uint16]map[int]ϡmemoizedResult
	ffp  ϡffp

	// stats
	matchCnt    int
	callCnt     int
	actionCnt   int
	codePredCnt int

	// error list
	errs errList
}

func (v *ϡvm) fromMemo(ix uint16, pt ϡsvpt) bool {
	if v.memo == nil {
		return false
	}
	m := v.memo[ix]
	if m == nil {
		return false
	}
	result, ok := m[pt.offset]
	if ok {
		v.parser.pt = result.pt
		v.v.push(result.v)
		return true
	}
	return false
}

func (v *ϡvm) memoizeMatch(ix uint16, pt ϡsvpt, match bool) {
	if v.memo == nil {
		v.memo = make(map[uint16]map[int]ϡmemoizedResult)
	}
	m, ok := v.memo[ix]
	if !ok {
		m = make(map[int]ϡmemoizedResult)
		v.memo[ix] = m
	}

	if match {
		m[pt.offset] = ϡmemoizedResult{v.parser.sliceFrom(pt), v.parser.pt}
		return
	}
	m[pt.offset] = ϡmemoizedResult{ϡmatchFailed, pt}
}

// setOptions applies the options in sequence on the vm. It returns the
// vm to allow for chaining calls.
func (v *ϡvm) setOptions(opts []Option) *ϡvm {
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// addErr adds the error at the current parser position, without rule name
// information.
func (v *ϡvm) addErr(err error) {
	v.addErrAt(err, -1, v.parser.pt.position)
}

// addErrAt adds the error at the specified position, for the rule name at
// ruleNmIx.
func (v *ϡvm) addErrAt(err error, ruleNmIx int, pos position) {
	var buf bytes.Buffer
	if v.filename != "" {
		buf.WriteString(v.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%s", pos))

	ruleNm := v.pg.ruleNameAt(ruleNmIx)
	if ruleNm != "" {
		buf.WriteString(": ")
		buf.WriteString("rule " + ruleNm)
	}

	pe := &parserError{Inner: err, ϡprefix: buf.String()}
	v.errs.ϡadd(pe)
}

// dumpSnapshot writes a dump of the current VM state to w.
func (v *ϡvm) dumpSnapshot(w io.Writer) {
	var buf bytes.Buffer

	if v.filename != "" {
		buf.WriteString(v.filename + ":")
	}
	buf.WriteString(fmt.Sprintf("%s: %#U\n", v.parser.pt.position, v.parser.pt.rn))

	// write the next 5 instructions
	ix := v.pc - 1
	if ix > 0 {
		ix--
	}
	stdFmt := ". [%d]: %s"
	for i := 0; i < 5; i++ {
		stdFmt := stdFmt
		if ix == v.pc-1 {
			stdFmt = ">" + stdFmt[1:]
		}
		instr := v.pg.instrs[ix]
		switch instr.op {
		case ϡopCall:
			buf.WriteString(fmt.Sprintf(stdFmt+"\n", ix, v.pg.instrToString(instr)))
			ix = v.i.pop() // continue with instructions at this index
			v.i.push(ix)
			continue
		default:
			buf.WriteString(fmt.Sprintf(stdFmt+"\n", ix, v.pg.instrToString(instr)))
		}
		ix++
		if int(ix) >= len(v.pg.instrs) {
			break
		}
	}

	// // print the stacks
	// buf.WriteString("[ P: ")
	// for i := 0; i < 3; i++ {
	// 	if len(v.p) <= i {
	// 		break
	// 	}
	// 	if i > 0 {
	// 		buf.WriteString(", ")
	// 	}
	// 	val := v.p[len(v.p)-i-1]
	// 	buf.WriteString(fmt.Sprintf("\"%v\"", val))
	// }
	// buf.WriteString(" ]\n[ V: ")
	// for i := 0; i < 3; i++ {
	// 	if len(v.v) <= i {
	// 		break
	// 	}
	// 	if i > 0 {
	// 		buf.WriteString(", ")
	// 	}
	// 	val := v.v[len(v.v)-i-1]
	// 	buf.WriteString(fmt.Sprintf("%#v", val))
	// }
	// buf.WriteString(" ]\n[ I: ")
	// for i := 0; i < 3; i++ {
	// 	if len(v.i) <= i {
	// 		break
	// 	}
	// 	if i > 0 {
	// 		buf.WriteString(", ")
	// 	}
	// 	val := v.i[len(v.i)-i-1]
	// 	buf.WriteString(fmt.Sprintf("%d", val))
	// }
	// buf.WriteString(" ]\n[ L: ")
	// for i := 0; i < 3; i++ {
	// 	if len(v.l) <= i {
	// 		break
	// 	}
	// 	if i > 0 {
	// 		buf.WriteString(", ")
	// 	}
	// 	val := v.l[len(v.l)-i-1]
	// 	buf.WriteString(fmt.Sprintf("%v", val))
	// }
	// buf.WriteString(" ]\n")
	fmt.Fprintln(w, buf.String())
}

// run executes the provided program in this VM, and returns the result.
func (v *ϡvm) run(pg *ϡprogram) (interface{}, error) {
	v.pg = pg
	v.a = newAstack(128)
	v.i = newIstack(128)
	v.v = newVstack(128)
	v.l = newLstack(128)
	v.p = newPstack(128)
	v.ffp.pos.offset = -1
	ret := v.dispatch()

	// if the match failed, translate that to a nil result and make
	// sure it returns an error
	if ret == ϡmatchFailed {
		ret = nil
		if len(v.errs) == 0 {
			if err := v.ffp.err(); err != nil {
				v.addErrAt(err, v.ffp.ruleNmIx, v.ffp.pos)
			} else {
				v.addErr(errNoMatch)
			}
		}
	}

	return ret, v.errs.ϡerr()
}

// dispatch is the proper execution method of the VM, it loops over
// the instructions and executes each opcode.
func (v *ϡvm) dispatch() interface{} {
	var instrPath []uint16
	if v.debug {
		fmt.Fprintln(os.Stderr, v.pg)
		defer func() {
			var buf bytes.Buffer

			buf.WriteString("Execution path:\n")
			for _, ix := range instrPath {
				buf.WriteString(fmt.Sprintf("[%3d]: %s\n", ix, v.pg.instrToString(v.pg.instrs[ix])))
			}
			fmt.Fprintln(os.Stderr, buf.String())
		}()
	}

	if v.recover {
		// if recover is set, recover from panics and convert to error.
		defer func() {
			if e := recover(); e != nil {
				ruleIx := -1
				if v.pc > 0 {
					ruleIx = v.pg.instrs[v.pc-1].ruleNmIx
				}
				switch e := e.(type) {
				case error:
					v.addErrAt(e, ruleIx, v.parser.pt.position)
				default:
					v.addErrAt(fmt.Errorf("%v", e), ruleIx, v.parser.pt.position)
				}
			}
		}()
	}

	// move to first rune before starting the loop
	v.parser.read()
	for {
		// fetch and decode the instruction
		instr := v.pg.instrs[v.pc]
		instrPath = append(instrPath, v.pc)

		// increment program counter
		v.pc++

		switch instr.op {
		case ϡopCall:
			if v.debug {
				v.dumpSnapshot(os.Stderr)
			}
			ix := v.i.pop()
			v.i.push(v.pc)
			v.pc = ix
			v.callCnt++

		case ϡopCallA:
			if v.debug {
				v.dumpSnapshot(os.Stderr)
			}
			v.v.pop()
			start := v.p.pop()
			v.cur.pos = start.position
			v.cur.text = v.parser.sliceFrom(start)
			if int(instr.args[0]) >= len(v.pg.as) {
				panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
			}
			fn := v.pg.as[instr.args[0]]
			val, err := fn(v)
			if err != nil {
				v.addErrAt(err, int(instr.ruleNmIx), start.position)
			}
			v.v.push(val)
			v.actionCnt++

		case ϡopCallB:
			if v.debug {
				v.dumpSnapshot(os.Stderr)
			}
			v.cur.pos = v.parser.pt.position
			v.cur.text = nil
			if int(instr.args[0]) >= len(v.pg.bs) {
				panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
			}
			fn := v.pg.bs[instr.args[0]]
			val, err := fn(v)
			if err != nil {
				v.addErrAt(err, int(instr.ruleNmIx), v.parser.pt.position)
			}
			v.codePredCnt++
			if !val {
				v.v.push(ϡmatchFailed)
				break
			}
			v.v.push(nil)

		case ϡopCumulOrF:
			va, vb := v.v.pop(), v.v.pop()
			if va == ϡmatchFailed {
				v.v.push(ϡmatchFailed)
				break
			}
			switch vb := vb.(type) {
			case []interface{}:
				vb = append(vb, va)
				v.v.push(vb)
			case ϡsentinel:
				v.v.push([]interface{}{va})
			default:
				panic(fmt.Sprintf("invalid %s value type on the V stack: %T", instr.op, vb))
			}

		case ϡopExit:
			return v.v.pop()

		case ϡopJump:
			v.pc = instr.args[0]

		case ϡopJumpIfF:
			if top := v.v.peek(); top == ϡmatchFailed {
				v.pc = instr.args[0]
			}

		case ϡopJumpIfT:
			if top := v.v.peek(); top != ϡmatchFailed {
				v.pc = instr.args[0]
			}

		case ϡopMatch:
			start := v.parser.pt
			if v.memoize && v.fromMemo(v.pc-1, start) {
				break
			}

			if int(instr.args[0]) >= len(v.pg.ms) {
				panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
			}
			m := v.pg.ms[instr.args[0]]
			ok := m.match(v.parser)
			if v.memoize {
				v.memoizeMatch(v.pc-1, start, ok)
			}
			v.matchCnt++
			if ok {
				v.v.push(v.parser.sliceFrom(start))
				break
			}
			// did not match, record ffp if this is the farthest failure
			if start.offset > v.ffp.pos.offset {
				v.ffp.pos = start.position
				v.ffp.got = v.parser.sliceFrom(start)
				v.ffp.ruleNmIx = instr.ruleNmIx
				v.ffp.want = m.toDisplayMsg()
			}
			v.v.push(ϡmatchFailed)
			v.parser.pt = start

			if v.debug {
				v.dumpSnapshot(os.Stderr)
			}

		case ϡopNilIfF:
			if top := v.v.pop(); top == ϡmatchFailed {
				v.v.push(nil)
				break
			}
			v.v.push(ϡmatchFailed)

		case ϡopNilIfT:
			if top := v.v.pop(); top != ϡmatchFailed {
				v.v.push(nil)
				break
			}
			v.v.push(ϡmatchFailed)

		case ϡopPop:
			switch instr.args[0] {
			case ϡlstackID:
				v.l.pop()
			case ϡpstackID:
				v.p.pop()
			case ϡastackID:
				v.a.pop()
			case ϡvstackID:
				v.v.pop()
			default:
				panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
			}

		case ϡopPopVJumpIfF:
			if top := v.v.peek(); top == ϡmatchFailed {
				v.v.pop()
				v.pc = instr.args[0]
			}

		case ϡopPush:
			switch instr.args[0] {
			case ϡpstackID:
				v.p.push(v.parser.pt)
			case ϡistackID:
				v.i.push(instr.args[1])
			case ϡvstackID:
				if int(instr.args[1]) >= len(ϡvSpecialValues) {
					panic(fmt.Sprintf("invalid %s V stack argument: %d", instr.op, instr.args[1]))
				}
				v.v.push(ϡvSpecialValues[instr.args[1]])
			case ϡastackID:
				v.a.push()
			case ϡlstackID:
				v.l.push(instr.args[1:])
			default:
				panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
			}

		case ϡopRestore:
			pt := v.p.pop()
			v.parser.pt = pt

		case ϡopRestoreIfF:
			pt := v.p.pop()
			if top := v.v.peek(); top == ϡmatchFailed {
				v.parser.pt = pt
			}

		case ϡopReturn:
			ix := v.i.pop()
			v.pc = ix

		case ϡopStoreIfT:
			if top := v.v.peek(); top != ϡmatchFailed {
				// get the label name
				if int(instr.args[0]) >= len(v.pg.ss) {
					panic(fmt.Sprintf("invalid %s argument: %d", instr.op, instr.args[0]))
				}
				lbl := v.pg.ss[instr.args[0]]

				// store the value
				as := v.a.peek()
				as[lbl] = top
			}

		case ϡopTakeLOrJump:
			ix := v.l.take()
			if ix < 0 {
				v.pc = instr.args[0]
				break
			}
			v.i.push(uint16(ix))

		default:
			panic(fmt.Sprintf("unknown opcode %s", instr.op))
		}
	}
}
