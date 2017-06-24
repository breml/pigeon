package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mna/pigeon/ast"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 5, col: 1, offset: 18},
			expr: &actionExpr{
				pos: position{line: 5, col: 11, offset: 28},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 5, col: 11, offset: 28},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 5, col: 11, offset: 28},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 5, col: 14, offset: 31},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 5, col: 26, offset: 43},
								expr: &seqExpr{
									pos: position{line: 5, col: 28, offset: 45},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 28, offset: 45},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 40, offset: 57},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 5, col: 46, offset: 63},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 5, col: 52, offset: 69},
								expr: &seqExpr{
									pos: position{line: 5, col: 54, offset: 71},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 54, offset: 71},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 59, offset: 76},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 5, col: 65, offset: 82},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Initializer",
			pos:  position{line: 24, col: 1, offset: 523},
			expr: &actionExpr{
				pos: position{line: 24, col: 15, offset: 537},
				run: (*parser).callonInitializer1,
				expr: &seqExpr{
					pos: position{line: 24, col: 15, offset: 537},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 24, col: 15, offset: 537},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 20, offset: 542},
								name: "CodeBlock",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 30, offset: 552},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 28, col: 1, offset: 582},
			expr: &actionExpr{
				pos: position{line: 28, col: 8, offset: 589},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 28, col: 8, offset: 589},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 8, offset: 589},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 13, offset: 594},
								name: "IdentifierName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 28, offset: 609},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 31, offset: 612},
							label: "display",
							expr: &zeroOrOneExpr{
								pos: position{line: 28, col: 39, offset: 620},
								expr: &seqExpr{
									pos: position{line: 28, col: 41, offset: 622},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 28, col: 41, offset: 622},
											name: "StringLiteral",
										},
										&ruleRefExpr{
											pos:  position{line: 28, col: 55, offset: 636},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 61, offset: 642},
							name: "RuleDefOp",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 71, offset: 652},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 74, offset: 655},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 28, col: 79, offset: 660},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 90, offset: 671},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 41, col: 1, offset: 955},
			expr: &ruleRefExpr{
				pos:  position{line: 41, col: 14, offset: 968},
				name: "ChoiceExpr",
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 43, col: 1, offset: 980},
			expr: &actionExpr{
				pos: position{line: 43, col: 14, offset: 993},
				run: (*parser).callonChoiceExpr1,
				expr: &seqExpr{
					pos: position{line: 43, col: 14, offset: 993},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 993},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 20, offset: 999},
								name: "ActionExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 43, col: 31, offset: 1010},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 36, offset: 1015},
								expr: &seqExpr{
									pos: position{line: 43, col: 38, offset: 1017},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 38, offset: 1017},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 43, col: 41, offset: 1020},
											val:        "/",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 45, offset: 1024},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 48, offset: 1027},
											name: "ActionExpr",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ActionExpr",
			pos:  position{line: 58, col: 1, offset: 1432},
			expr: &actionExpr{
				pos: position{line: 58, col: 14, offset: 1445},
				run: (*parser).callonActionExpr1,
				expr: &seqExpr{
					pos: position{line: 58, col: 14, offset: 1445},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 58, col: 14, offset: 1445},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 58, col: 19, offset: 1450},
								name: "SeqExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 58, col: 27, offset: 1458},
							label: "code",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 32, offset: 1463},
								expr: &seqExpr{
									pos: position{line: 58, col: 34, offset: 1465},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 58, col: 34, offset: 1465},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 58, col: 37, offset: 1468},
											name: "CodeBlock",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 72, col: 1, offset: 1734},
			expr: &actionExpr{
				pos: position{line: 72, col: 11, offset: 1744},
				run: (*parser).callonSeqExpr1,
				expr: &seqExpr{
					pos: position{line: 72, col: 11, offset: 1744},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 72, col: 11, offset: 1744},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 17, offset: 1750},
								name: "LabeledExpr",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 29, offset: 1762},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 34, offset: 1767},
								expr: &seqExpr{
									pos: position{line: 72, col: 36, offset: 1769},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 72, col: 36, offset: 1769},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 72, col: 39, offset: 1772},
											name: "LabeledExpr",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 85, col: 1, offset: 2123},
			expr: &choiceExpr{
				pos: position{line: 85, col: 15, offset: 2137},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 15, offset: 2137},
						run: (*parser).callonLabeledExpr2,
						expr: &seqExpr{
							pos: position{line: 85, col: 15, offset: 2137},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 85, col: 15, offset: 2137},
									label: "label",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 21, offset: 2143},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 32, offset: 2154},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 85, col: 35, offset: 2157},
									val:        ":",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 39, offset: 2161},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 85, col: 42, offset: 2164},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 85, col: 47, offset: 2169},
										name: "PrefixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2342},
						name: "PrefixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 93, col: 1, offset: 2356},
			expr: &choiceExpr{
				pos: position{line: 93, col: 16, offset: 2371},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 93, col: 16, offset: 2371},
						run: (*parser).callonPrefixedExpr2,
						expr: &seqExpr{
							pos: position{line: 93, col: 16, offset: 2371},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 93, col: 16, offset: 2371},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 19, offset: 2374},
										name: "PrefixedOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 93, col: 30, offset: 2385},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 93, col: 33, offset: 2388},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 93, col: 38, offset: 2393},
										name: "SuffixedExpr",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 5, offset: 2675},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 106, col: 1, offset: 2689},
			expr: &actionExpr{
				pos: position{line: 106, col: 14, offset: 2702},
				run: (*parser).callonPrefixedOp1,
				expr: &choiceExpr{
					pos: position{line: 106, col: 16, offset: 2704},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 16, offset: 2704},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 106, col: 22, offset: 2710},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 110, col: 1, offset: 2752},
			expr: &choiceExpr{
				pos: position{line: 110, col: 16, offset: 2767},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 110, col: 16, offset: 2767},
						run: (*parser).callonSuffixedExpr2,
						expr: &seqExpr{
							pos: position{line: 110, col: 16, offset: 2767},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 110, col: 16, offset: 2767},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 21, offset: 2772},
										name: "PrimaryExpr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 33, offset: 2784},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 36, offset: 2787},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 39, offset: 2790},
										name: "SuffixedOp",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 129, col: 5, offset: 3320},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 131, col: 1, offset: 3334},
			expr: &actionExpr{
				pos: position{line: 131, col: 14, offset: 3347},
				run: (*parser).callonSuffixedOp1,
				expr: &choiceExpr{
					pos: position{line: 131, col: 16, offset: 3349},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 16, offset: 3349},
							val:        "?",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 22, offset: 3355},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 131, col: 28, offset: 3361},
							val:        "+",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 135, col: 1, offset: 3403},
			expr: &choiceExpr{
				pos: position{line: 135, col: 15, offset: 3417},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 135, col: 15, offset: 3417},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 28, offset: 3430},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 47, offset: 3449},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 60, offset: 3462},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 74, offset: 3476},
						name: "SemanticPredExpr",
					},
					&actionExpr{
						pos: position{line: 135, col: 93, offset: 3495},
						run: (*parser).callonPrimaryExpr7,
						expr: &seqExpr{
							pos: position{line: 135, col: 93, offset: 3495},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 93, offset: 3495},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 97, offset: 3499},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 100, offset: 3502},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 105, offset: 3507},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 116, offset: 3518},
									name: "__",
								},
								&litMatcher{
									pos:        position{line: 135, col: 119, offset: 3521},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 138, col: 1, offset: 3550},
			expr: &actionExpr{
				pos: position{line: 138, col: 15, offset: 3564},
				run: (*parser).callonRuleRefExpr1,
				expr: &seqExpr{
					pos: position{line: 138, col: 15, offset: 3564},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 138, col: 15, offset: 3564},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 20, offset: 3569},
								name: "IdentifierName",
							},
						},
						&notExpr{
							pos: position{line: 138, col: 35, offset: 3584},
							expr: &seqExpr{
								pos: position{line: 138, col: 38, offset: 3587},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 138, col: 38, offset: 3587},
										name: "__",
									},
									&zeroOrOneExpr{
										pos: position{line: 138, col: 41, offset: 3590},
										expr: &seqExpr{
											pos: position{line: 138, col: 43, offset: 3592},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 138, col: 43, offset: 3592},
													name: "StringLiteral",
												},
												&ruleRefExpr{
													pos:  position{line: 138, col: 57, offset: 3606},
													name: "__",
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 138, col: 63, offset: 3612},
										name: "RuleDefOp",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredExpr",
			pos:  position{line: 143, col: 1, offset: 3728},
			expr: &actionExpr{
				pos: position{line: 143, col: 20, offset: 3747},
				run: (*parser).callonSemanticPredExpr1,
				expr: &seqExpr{
					pos: position{line: 143, col: 20, offset: 3747},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 143, col: 20, offset: 3747},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 23, offset: 3750},
								name: "SemanticPredOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 143, col: 38, offset: 3765},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 143, col: 41, offset: 3768},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 143, col: 46, offset: 3773},
								name: "CodeBlock",
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 154, col: 1, offset: 4050},
			expr: &actionExpr{
				pos: position{line: 154, col: 18, offset: 4067},
				run: (*parser).callonSemanticPredOp1,
				expr: &choiceExpr{
					pos: position{line: 154, col: 20, offset: 4069},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 20, offset: 4069},
							val:        "&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 154, col: 26, offset: 4075},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 158, col: 1, offset: 4117},
			expr: &choiceExpr{
				pos: position{line: 158, col: 13, offset: 4129},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 158, col: 13, offset: 4129},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 19, offset: 4135},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 26, offset: 4142},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 158, col: 37, offset: 4153},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 160, col: 1, offset: 4163},
			expr: &anyMatcher{
				line: 160, col: 14, offset: 4176,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 161, col: 1, offset: 4178},
			expr: &choiceExpr{
				pos: position{line: 161, col: 11, offset: 4188},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 161, col: 11, offset: 4188},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 161, col: 30, offset: 4207},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 162, col: 1, offset: 4225},
			expr: &seqExpr{
				pos: position{line: 162, col: 20, offset: 4244},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 162, col: 20, offset: 4244},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 162, col: 25, offset: 4249},
						expr: &seqExpr{
							pos: position{line: 162, col: 27, offset: 4251},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 162, col: 27, offset: 4251},
									expr: &litMatcher{
										pos:        position{line: 162, col: 28, offset: 4252},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 162, col: 33, offset: 4257},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 162, col: 47, offset: 4271},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 163, col: 1, offset: 4276},
			expr: &seqExpr{
				pos: position{line: 163, col: 36, offset: 4311},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 163, col: 36, offset: 4311},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 163, col: 41, offset: 4316},
						expr: &seqExpr{
							pos: position{line: 163, col: 43, offset: 4318},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 163, col: 43, offset: 4318},
									expr: &choiceExpr{
										pos: position{line: 163, col: 46, offset: 4321},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 46, offset: 4321},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 53, offset: 4328},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 163, col: 59, offset: 4334},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 163, col: 73, offset: 4348},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 164, col: 1, offset: 4353},
			expr: &seqExpr{
				pos: position{line: 164, col: 21, offset: 4373},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 164, col: 21, offset: 4373},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 164, col: 26, offset: 4378},
						expr: &seqExpr{
							pos: position{line: 164, col: 28, offset: 4380},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 164, col: 28, offset: 4380},
									expr: &ruleRefExpr{
										pos:  position{line: 164, col: 29, offset: 4381},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 164, col: 33, offset: 4385},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 166, col: 1, offset: 4400},
			expr: &actionExpr{
				pos: position{line: 166, col: 14, offset: 4413},
				run: (*parser).callonIdentifier1,
				expr: &labeledExpr{
					pos:   position{line: 166, col: 14, offset: 4413},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 166, col: 20, offset: 4419},
						name: "IdentifierName",
					},
				},
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 174, col: 1, offset: 4638},
			expr: &actionExpr{
				pos: position{line: 174, col: 18, offset: 4655},
				run: (*parser).callonIdentifierName1,
				expr: &seqExpr{
					pos: position{line: 174, col: 18, offset: 4655},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 174, col: 18, offset: 4655},
							name: "IdentifierStart",
						},
						&zeroOrMoreExpr{
							pos: position{line: 174, col: 34, offset: 4671},
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 34, offset: 4671},
								name: "IdentifierPart",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 177, col: 1, offset: 4753},
			expr: &charClassMatcher{
				pos:        position{line: 177, col: 19, offset: 4771},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 178, col: 1, offset: 4778},
			expr: &choiceExpr{
				pos: position{line: 178, col: 18, offset: 4795},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 178, col: 18, offset: 4795},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 178, col: 36, offset: 4813},
						val:        "[\\p{Nd}]",
						classes:    []*unicode.RangeTable{rangeTable("Nd")},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 180, col: 1, offset: 4823},
			expr: &actionExpr{
				pos: position{line: 180, col: 14, offset: 4836},
				run: (*parser).callonLitMatcher1,
				expr: &seqExpr{
					pos: position{line: 180, col: 14, offset: 4836},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 180, col: 14, offset: 4836},
							label: "lit",
							expr: &ruleRefExpr{
								pos:  position{line: 180, col: 18, offset: 4840},
								name: "StringLiteral",
							},
						},
						&labeledExpr{
							pos:   position{line: 180, col: 32, offset: 4854},
							label: "ignore",
							expr: &zeroOrOneExpr{
								pos: position{line: 180, col: 39, offset: 4861},
								expr: &litMatcher{
									pos:        position{line: 180, col: 39, offset: 4861},
									val:        "i",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 193, col: 1, offset: 5260},
			expr: &choiceExpr{
				pos: position{line: 193, col: 17, offset: 5276},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 193, col: 17, offset: 5276},
						run: (*parser).callonStringLiteral2,
						expr: &choiceExpr{
							pos: position{line: 193, col: 19, offset: 5278},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 193, col: 19, offset: 5278},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 19, offset: 5278},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 193, col: 23, offset: 5282},
											expr: &ruleRefExpr{
												pos:  position{line: 193, col: 23, offset: 5282},
												name: "DoubleStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 193, col: 41, offset: 5300},
											val:        "\"",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 193, col: 47, offset: 5306},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 47, offset: 5306},
											val:        "'",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 193, col: 51, offset: 5310},
											name: "SingleStringChar",
										},
										&litMatcher{
											pos:        position{line: 193, col: 68, offset: 5327},
											val:        "'",
											ignoreCase: false,
										},
									},
								},
								&seqExpr{
									pos: position{line: 193, col: 74, offset: 5333},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 193, col: 74, offset: 5333},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 193, col: 78, offset: 5337},
											expr: &ruleRefExpr{
												pos:  position{line: 193, col: 78, offset: 5337},
												name: "RawStringChar",
											},
										},
										&litMatcher{
											pos:        position{line: 193, col: 93, offset: 5352},
											val:        "`",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 195, col: 5, offset: 5425},
						run: (*parser).callonStringLiteral18,
						expr: &choiceExpr{
							pos: position{line: 195, col: 7, offset: 5427},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 195, col: 9, offset: 5429},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 9, offset: 5429},
											val:        "\"",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 195, col: 13, offset: 5433},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 13, offset: 5433},
												name: "DoubleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 195, col: 33, offset: 5453},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 195, col: 33, offset: 5453},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 195, col: 39, offset: 5459},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 195, col: 51, offset: 5471},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 51, offset: 5471},
											val:        "'",
											ignoreCase: false,
										},
										&zeroOrOneExpr{
											pos: position{line: 195, col: 55, offset: 5475},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 55, offset: 5475},
												name: "SingleStringChar",
											},
										},
										&choiceExpr{
											pos: position{line: 195, col: 75, offset: 5495},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 195, col: 75, offset: 5495},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 195, col: 81, offset: 5501},
													name: "EOF",
												},
											},
										},
									},
								},
								&seqExpr{
									pos: position{line: 195, col: 91, offset: 5511},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 195, col: 91, offset: 5511},
											val:        "`",
											ignoreCase: false,
										},
										&zeroOrMoreExpr{
											pos: position{line: 195, col: 95, offset: 5515},
											expr: &ruleRefExpr{
												pos:  position{line: 195, col: 95, offset: 5515},
												name: "RawStringChar",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 195, col: 110, offset: 5530},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleStringChar",
			pos:  position{line: 199, col: 1, offset: 5632},
			expr: &choiceExpr{
				pos: position{line: 199, col: 20, offset: 5651},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 199, col: 20, offset: 5651},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 199, col: 20, offset: 5651},
								expr: &choiceExpr{
									pos: position{line: 199, col: 23, offset: 5654},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 199, col: 23, offset: 5654},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 199, col: 29, offset: 5660},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 199, col: 36, offset: 5667},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 199, col: 42, offset: 5673},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 199, col: 55, offset: 5686},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 199, col: 55, offset: 5686},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 199, col: 60, offset: 5691},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 200, col: 1, offset: 5710},
			expr: &choiceExpr{
				pos: position{line: 200, col: 20, offset: 5729},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 200, col: 20, offset: 5729},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 200, col: 20, offset: 5729},
								expr: &choiceExpr{
									pos: position{line: 200, col: 23, offset: 5732},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 200, col: 23, offset: 5732},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 200, col: 29, offset: 5738},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 200, col: 36, offset: 5745},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 42, offset: 5751},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 200, col: 55, offset: 5764},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 200, col: 55, offset: 5764},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 200, col: 60, offset: 5769},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 201, col: 1, offset: 5788},
			expr: &seqExpr{
				pos: position{line: 201, col: 17, offset: 5804},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 201, col: 17, offset: 5804},
						expr: &litMatcher{
							pos:        position{line: 201, col: 18, offset: 5805},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 201, col: 22, offset: 5809},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 203, col: 1, offset: 5821},
			expr: &choiceExpr{
				pos: position{line: 203, col: 22, offset: 5842},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 203, col: 24, offset: 5844},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 203, col: 24, offset: 5844},
								val:        "\"",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 203, col: 30, offset: 5850},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 204, col: 7, offset: 5879},
						run: (*parser).callonDoubleStringEscape5,
						expr: &choiceExpr{
							pos: position{line: 204, col: 9, offset: 5881},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 204, col: 9, offset: 5881},
									name: "SourceChar",
								},
								&ruleRefExpr{
									pos:  position{line: 204, col: 22, offset: 5894},
									name: "EOL",
								},
								&ruleRefExpr{
									pos:  position{line: 204, col: 28, offset: 5900},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 207, col: 1, offset: 5965},
			expr: &choiceExpr{
				pos: position{line: 207, col: 22, offset: 5986},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 207, col: 24, offset: 5988},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 207, col: 24, offset: 5988},
								val:        "'",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 207, col: 30, offset: 5994},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 208, col: 7, offset: 6023},
						run: (*parser).callonSingleStringEscape5,
						expr: &choiceExpr{
							pos: position{line: 208, col: 9, offset: 6025},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 208, col: 9, offset: 6025},
									name: "SourceChar",
								},
								&ruleRefExpr{
									pos:  position{line: 208, col: 22, offset: 6038},
									name: "EOL",
								},
								&ruleRefExpr{
									pos:  position{line: 208, col: 28, offset: 6044},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 212, col: 1, offset: 6110},
			expr: &choiceExpr{
				pos: position{line: 212, col: 24, offset: 6133},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 212, col: 24, offset: 6133},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 43, offset: 6152},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 57, offset: 6166},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 69, offset: 6178},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 89, offset: 6198},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 213, col: 1, offset: 6217},
			expr: &choiceExpr{
				pos: position{line: 213, col: 20, offset: 6236},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 213, col: 20, offset: 6236},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 26, offset: 6242},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 32, offset: 6248},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 38, offset: 6254},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 44, offset: 6260},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 50, offset: 6266},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 56, offset: 6272},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 213, col: 62, offset: 6278},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 214, col: 1, offset: 6283},
			expr: &choiceExpr{
				pos: position{line: 214, col: 15, offset: 6297},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 214, col: 15, offset: 6297},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 214, col: 15, offset: 6297},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 26, offset: 6308},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 214, col: 37, offset: 6319},
								name: "OctalDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 215, col: 7, offset: 6336},
						run: (*parser).callonOctalEscape6,
						expr: &seqExpr{
							pos: position{line: 215, col: 7, offset: 6336},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 215, col: 7, offset: 6336},
									name: "OctalDigit",
								},
								&choiceExpr{
									pos: position{line: 215, col: 20, offset: 6349},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 215, col: 20, offset: 6349},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 215, col: 33, offset: 6362},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 215, col: 39, offset: 6368},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 218, col: 1, offset: 6429},
			expr: &choiceExpr{
				pos: position{line: 218, col: 13, offset: 6441},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 218, col: 13, offset: 6441},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 218, col: 13, offset: 6441},
								val:        "x",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 218, col: 17, offset: 6445},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 218, col: 26, offset: 6454},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 219, col: 7, offset: 6469},
						run: (*parser).callonHexEscape6,
						expr: &seqExpr{
							pos: position{line: 219, col: 7, offset: 6469},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 219, col: 7, offset: 6469},
									val:        "x",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 219, col: 13, offset: 6475},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 219, col: 13, offset: 6475},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 219, col: 26, offset: 6488},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 219, col: 32, offset: 6494},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 222, col: 1, offset: 6561},
			expr: &choiceExpr{
				pos: position{line: 223, col: 5, offset: 6586},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 223, col: 5, offset: 6586},
						run: (*parser).callonLongUnicodeEscape2,
						expr: &seqExpr{
							pos: position{line: 223, col: 5, offset: 6586},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 223, col: 5, offset: 6586},
									val:        "U",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 9, offset: 6590},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 18, offset: 6599},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 27, offset: 6608},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 36, offset: 6617},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 45, offset: 6626},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 54, offset: 6635},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 63, offset: 6644},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 223, col: 72, offset: 6653},
									name: "HexDigit",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 226, col: 7, offset: 6755},
						run: (*parser).callonLongUnicodeEscape13,
						expr: &seqExpr{
							pos: position{line: 226, col: 7, offset: 6755},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 226, col: 7, offset: 6755},
									val:        "U",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 226, col: 13, offset: 6761},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 226, col: 13, offset: 6761},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 226, col: 26, offset: 6774},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 226, col: 32, offset: 6780},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 229, col: 1, offset: 6843},
			expr: &choiceExpr{
				pos: position{line: 230, col: 5, offset: 6869},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 230, col: 5, offset: 6869},
						run: (*parser).callonShortUnicodeEscape2,
						expr: &seqExpr{
							pos: position{line: 230, col: 5, offset: 6869},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 230, col: 5, offset: 6869},
									val:        "u",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 9, offset: 6873},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 18, offset: 6882},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 27, offset: 6891},
									name: "HexDigit",
								},
								&ruleRefExpr{
									pos:  position{line: 230, col: 36, offset: 6900},
									name: "HexDigit",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 233, col: 7, offset: 7002},
						run: (*parser).callonShortUnicodeEscape9,
						expr: &seqExpr{
							pos: position{line: 233, col: 7, offset: 7002},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 233, col: 7, offset: 7002},
									val:        "u",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 233, col: 13, offset: 7008},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 233, col: 13, offset: 7008},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 26, offset: 7021},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 32, offset: 7027},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 237, col: 1, offset: 7091},
			expr: &charClassMatcher{
				pos:        position{line: 237, col: 14, offset: 7104},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 238, col: 1, offset: 7110},
			expr: &charClassMatcher{
				pos:        position{line: 238, col: 16, offset: 7125},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 239, col: 1, offset: 7131},
			expr: &charClassMatcher{
				pos:        position{line: 239, col: 12, offset: 7142},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 241, col: 1, offset: 7153},
			expr: &choiceExpr{
				pos: position{line: 241, col: 20, offset: 7172},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 241, col: 20, offset: 7172},
						run: (*parser).callonCharClassMatcher2,
						expr: &seqExpr{
							pos: position{line: 241, col: 20, offset: 7172},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 241, col: 20, offset: 7172},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 241, col: 24, offset: 7176},
									expr: &choiceExpr{
										pos: position{line: 241, col: 26, offset: 7178},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 241, col: 26, offset: 7178},
												name: "ClassCharRange",
											},
											&ruleRefExpr{
												pos:  position{line: 241, col: 43, offset: 7195},
												name: "ClassChar",
											},
											&seqExpr{
												pos: position{line: 241, col: 55, offset: 7207},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 241, col: 55, offset: 7207},
														val:        "\\",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 241, col: 60, offset: 7212},
														name: "UnicodeClassEscape",
													},
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 241, col: 82, offset: 7234},
									val:        "]",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 241, col: 86, offset: 7238},
									expr: &litMatcher{
										pos:        position{line: 241, col: 86, offset: 7238},
										val:        "i",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 245, col: 5, offset: 7345},
						run: (*parser).callonCharClassMatcher15,
						expr: &seqExpr{
							pos: position{line: 245, col: 5, offset: 7345},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 245, col: 5, offset: 7345},
									val:        "[",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 245, col: 9, offset: 7349},
									expr: &seqExpr{
										pos: position{line: 245, col: 11, offset: 7351},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 245, col: 11, offset: 7351},
												expr: &ruleRefExpr{
													pos:  position{line: 245, col: 14, offset: 7354},
													name: "EOL",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 245, col: 20, offset: 7360},
												name: "SourceChar",
											},
										},
									},
								},
								&choiceExpr{
									pos: position{line: 245, col: 36, offset: 7376},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 245, col: 36, offset: 7376},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 245, col: 42, offset: 7382},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 249, col: 1, offset: 7492},
			expr: &seqExpr{
				pos: position{line: 249, col: 18, offset: 7509},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 249, col: 18, offset: 7509},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 249, col: 28, offset: 7519},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 249, col: 32, offset: 7523},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 250, col: 1, offset: 7533},
			expr: &choiceExpr{
				pos: position{line: 250, col: 13, offset: 7545},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 250, col: 13, offset: 7545},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 250, col: 13, offset: 7545},
								expr: &choiceExpr{
									pos: position{line: 250, col: 16, offset: 7548},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 250, col: 16, offset: 7548},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 250, col: 22, offset: 7554},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 250, col: 29, offset: 7561},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 250, col: 35, offset: 7567},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 250, col: 48, offset: 7580},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 250, col: 48, offset: 7580},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 250, col: 53, offset: 7585},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 251, col: 1, offset: 7601},
			expr: &choiceExpr{
				pos: position{line: 251, col: 19, offset: 7619},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 251, col: 21, offset: 7621},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 251, col: 21, offset: 7621},
								val:        "]",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 251, col: 27, offset: 7627},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 252, col: 7, offset: 7656},
						run: (*parser).callonCharClassEscape5,
						expr: &seqExpr{
							pos: position{line: 252, col: 7, offset: 7656},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 252, col: 7, offset: 7656},
									expr: &litMatcher{
										pos:        position{line: 252, col: 8, offset: 7657},
										val:        "p",
										ignoreCase: false,
									},
								},
								&choiceExpr{
									pos: position{line: 252, col: 14, offset: 7663},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 252, col: 14, offset: 7663},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 252, col: 27, offset: 7676},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 252, col: 33, offset: 7682},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 256, col: 1, offset: 7748},
			expr: &seqExpr{
				pos: position{line: 256, col: 22, offset: 7769},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 256, col: 22, offset: 7769},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 257, col: 7, offset: 7782},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 257, col: 7, offset: 7782},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 258, col: 7, offset: 7811},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 258, col: 7, offset: 7811},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 258, col: 7, offset: 7811},
											expr: &litMatcher{
												pos:        position{line: 258, col: 8, offset: 7812},
												val:        "{",
												ignoreCase: false,
											},
										},
										&choiceExpr{
											pos: position{line: 258, col: 14, offset: 7818},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 258, col: 14, offset: 7818},
													name: "SourceChar",
												},
												&ruleRefExpr{
													pos:  position{line: 258, col: 27, offset: 7831},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 258, col: 33, offset: 7837},
													name: "EOF",
												},
											},
										},
									},
								},
							},
							&actionExpr{
								pos: position{line: 259, col: 7, offset: 7908},
								run: (*parser).callonUnicodeClassEscape13,
								expr: &seqExpr{
									pos: position{line: 259, col: 7, offset: 7908},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 259, col: 7, offset: 7908},
											val:        "{",
											ignoreCase: false,
										},
										&labeledExpr{
											pos:   position{line: 259, col: 11, offset: 7912},
											label: "ident",
											expr: &ruleRefExpr{
												pos:  position{line: 259, col: 17, offset: 7918},
												name: "IdentifierName",
											},
										},
										&litMatcher{
											pos:        position{line: 259, col: 32, offset: 7933},
											val:        "}",
											ignoreCase: false,
										},
									},
								},
							},
							&actionExpr{
								pos: position{line: 265, col: 7, offset: 8110},
								run: (*parser).callonUnicodeClassEscape19,
								expr: &seqExpr{
									pos: position{line: 265, col: 7, offset: 8110},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 265, col: 7, offset: 8110},
											val:        "{",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 265, col: 11, offset: 8114},
											name: "IdentifierName",
										},
										&choiceExpr{
											pos: position{line: 265, col: 28, offset: 8131},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 265, col: 28, offset: 8131},
													val:        "]",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 265, col: 34, offset: 8137},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 265, col: 40, offset: 8143},
													name: "EOF",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleCharUnicodeClass",
			pos:  position{line: 269, col: 1, offset: 8226},
			expr: &charClassMatcher{
				pos:        position{line: 269, col: 26, offset: 8251},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 271, col: 1, offset: 8262},
			expr: &actionExpr{
				pos: position{line: 271, col: 14, offset: 8275},
				run: (*parser).callonAnyMatcher1,
				expr: &litMatcher{
					pos:        position{line: 271, col: 14, offset: 8275},
					val:        ".",
					ignoreCase: false,
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 276, col: 1, offset: 8350},
			expr: &choiceExpr{
				pos: position{line: 276, col: 13, offset: 8362},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 276, col: 13, offset: 8362},
						run: (*parser).callonCodeBlock2,
						expr: &seqExpr{
							pos: position{line: 276, col: 13, offset: 8362},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 276, col: 13, offset: 8362},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 276, col: 17, offset: 8366},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 276, col: 22, offset: 8371},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 280, col: 5, offset: 8470},
						run: (*parser).callonCodeBlock7,
						expr: &seqExpr{
							pos: position{line: 280, col: 5, offset: 8470},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 280, col: 5, offset: 8470},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 280, col: 9, offset: 8474},
									name: "Code",
								},
								&ruleRefExpr{
									pos:  position{line: 280, col: 14, offset: 8479},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 284, col: 1, offset: 8544},
			expr: &zeroOrMoreExpr{
				pos: position{line: 284, col: 8, offset: 8551},
				expr: &choiceExpr{
					pos: position{line: 284, col: 10, offset: 8553},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 284, col: 10, offset: 8553},
							expr: &seqExpr{
								pos: position{line: 284, col: 12, offset: 8555},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 284, col: 12, offset: 8555},
										expr: &charClassMatcher{
											pos:        position{line: 284, col: 13, offset: 8556},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 284, col: 18, offset: 8561},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 284, col: 34, offset: 8577},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 284, col: 34, offset: 8577},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 284, col: 38, offset: 8581},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 284, col: 43, offset: 8586},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 286, col: 1, offset: 8594},
			expr: &zeroOrMoreExpr{
				pos: position{line: 286, col: 6, offset: 8599},
				expr: &choiceExpr{
					pos: position{line: 286, col: 8, offset: 8601},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 286, col: 8, offset: 8601},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 286, col: 21, offset: 8614},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 286, col: 27, offset: 8620},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 287, col: 1, offset: 8631},
			expr: &zeroOrMoreExpr{
				pos: position{line: 287, col: 5, offset: 8635},
				expr: &choiceExpr{
					pos: position{line: 287, col: 7, offset: 8637},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 287, col: 7, offset: 8637},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 20, offset: 8650},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 289, col: 1, offset: 8687},
			expr: &charClassMatcher{
				pos:        position{line: 289, col: 14, offset: 8700},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 290, col: 1, offset: 8708},
			expr: &litMatcher{
				pos:        position{line: 290, col: 7, offset: 8714},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 291, col: 1, offset: 8719},
			expr: &choiceExpr{
				pos: position{line: 291, col: 7, offset: 8725},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 291, col: 7, offset: 8725},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 7, offset: 8725},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 291, col: 10, offset: 8728},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 291, col: 16, offset: 8734},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 16, offset: 8734},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 291, col: 18, offset: 8736},
								expr: &ruleRefExpr{
									pos:  position{line: 291, col: 18, offset: 8736},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 37, offset: 8755},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 291, col: 43, offset: 8761},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 291, col: 43, offset: 8761},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 291, col: 46, offset: 8764},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 293, col: 1, offset: 8769},
			expr: &notExpr{
				pos: position{line: 293, col: 7, offset: 8775},
				expr: &anyMatcher{
					line: 293, col: 8, offset: 8776,
				},
			},
		},
	},
}

func (c *current) onGrammar1(initializer, rules interface{}) (interface{}, error) {
	pos := c.astPos()

	// create the grammar, assign its initializer
	g := ast.NewGrammar(pos)
	initSlice := toIfaceSlice(initializer)
	if len(initSlice) > 0 {
		g.Init = initSlice[0].(*ast.CodeBlock)
	}

	rulesSlice := toIfaceSlice(rules)
	g.Rules = make([]*ast.Rule, len(rulesSlice))
	for i, duo := range rulesSlice {
		g.Rules[i] = duo.([]interface{})[0].(*ast.Rule)
	}

	return g, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["initializer"], stack["rules"])
}

func (c *current) onInitializer1(code interface{}) (interface{}, error) {
	return code, nil
}

func (p *parser) callonInitializer1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInitializer1(stack["code"])
}

func (c *current) onRule1(name, display, expr interface{}) (interface{}, error) {
	pos := c.astPos()

	rule := ast.NewRule(pos, name.(*ast.Identifier))
	displaySlice := toIfaceSlice(display)
	if len(displaySlice) > 0 {
		rule.DisplayName = displaySlice[0].(*ast.StringLit)
	}
	rule.Expr = expr.(ast.Expression)

	return rule, nil
}

func (p *parser) callonRule1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule1(stack["name"], stack["display"], stack["expr"])
}

func (c *current) onChoiceExpr1(first, rest interface{}) (interface{}, error) {
	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}

	pos := c.astPos()
	choice := ast.NewChoiceExpr(pos)
	choice.Alternatives = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		choice.Alternatives = append(choice.Alternatives, sl.([]interface{})[3].(ast.Expression))
	}
	return choice, nil
}

func (p *parser) callonChoiceExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChoiceExpr1(stack["first"], stack["rest"])
}

func (c *current) onActionExpr1(expr, code interface{}) (interface{}, error) {
	if code == nil {
		return expr, nil
	}

	pos := c.astPos()
	act := ast.NewActionExpr(pos)
	act.Expr = expr.(ast.Expression)
	codeSlice := toIfaceSlice(code)
	act.Code = codeSlice[1].(*ast.CodeBlock)

	return act, nil
}

func (p *parser) callonActionExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onActionExpr1(stack["expr"], stack["code"])
}

func (c *current) onSeqExpr1(first, rest interface{}) (interface{}, error) {
	restSlice := toIfaceSlice(rest)
	if len(restSlice) == 0 {
		return first, nil
	}
	seq := ast.NewSeqExpr(c.astPos())
	seq.Exprs = []ast.Expression{first.(ast.Expression)}
	for _, sl := range restSlice {
		seq.Exprs = append(seq.Exprs, sl.([]interface{})[1].(ast.Expression))
	}
	return seq, nil
}

func (p *parser) callonSeqExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSeqExpr1(stack["first"], stack["rest"])
}

func (c *current) onLabeledExpr2(label, expr interface{}) (interface{}, error) {
	pos := c.astPos()
	lab := ast.NewLabeledExpr(pos)
	lab.Label = label.(*ast.Identifier)
	lab.Expr = expr.(ast.Expression)
	return lab, nil
}

func (p *parser) callonLabeledExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLabeledExpr2(stack["label"], stack["expr"])
}

func (c *current) onPrefixedExpr2(op, expr interface{}) (interface{}, error) {
	pos := c.astPos()
	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndExpr(pos)
		and.Expr = expr.(ast.Expression)
		return and, nil
	}
	not := ast.NewNotExpr(pos)
	not.Expr = expr.(ast.Expression)
	return not, nil
}

func (p *parser) callonPrefixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedExpr2(stack["op"], stack["expr"])
}

func (c *current) onPrefixedOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonPrefixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrefixedOp1()
}

func (c *current) onSuffixedExpr2(expr, op interface{}) (interface{}, error) {
	pos := c.astPos()
	opStr := op.(string)
	switch opStr {
	case "?":
		zero := ast.NewZeroOrOneExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "*":
		zero := ast.NewZeroOrMoreExpr(pos)
		zero.Expr = expr.(ast.Expression)
		return zero, nil
	case "+":
		one := ast.NewOneOrMoreExpr(pos)
		one.Expr = expr.(ast.Expression)
		return one, nil
	default:
		return nil, errors.New("unknown operator: " + opStr)
	}
}

func (p *parser) callonSuffixedExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedExpr2(stack["expr"], stack["op"])
}

func (c *current) onSuffixedOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSuffixedOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSuffixedOp1()
}

func (c *current) onPrimaryExpr7(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonPrimaryExpr7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimaryExpr7(stack["expr"])
}

func (c *current) onRuleRefExpr1(name interface{}) (interface{}, error) {
	ref := ast.NewRuleRefExpr(c.astPos())
	ref.Name = name.(*ast.Identifier)
	return ref, nil
}

func (p *parser) callonRuleRefExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRuleRefExpr1(stack["name"])
}

func (c *current) onSemanticPredExpr1(op, code interface{}) (interface{}, error) {
	opStr := op.(string)
	if opStr == "&" {
		and := ast.NewAndCodeExpr(c.astPos())
		and.Code = code.(*ast.CodeBlock)
		return and, nil
	}
	not := ast.NewNotCodeExpr(c.astPos())
	not.Code = code.(*ast.CodeBlock)
	return not, nil
}

func (p *parser) callonSemanticPredExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredExpr1(stack["op"], stack["code"])
}

func (c *current) onSemanticPredOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonSemanticPredOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSemanticPredOp1()
}

func (c *current) onIdentifier1(ident interface{}) (interface{}, error) {
	astIdent := ast.NewIdentifier(c.astPos(), string(c.text))
	if reservedWords[astIdent.Val] {
		return astIdent, errors.New("identifier is a reserved word")
	}
	return astIdent, nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1(stack["ident"])
}

func (c *current) onIdentifierName1() (interface{}, error) {
	return ast.NewIdentifier(c.astPos(), string(c.text)), nil
}

func (p *parser) callonIdentifierName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierName1()
}

func (c *current) onLitMatcher1(lit, ignore interface{}) (interface{}, error) {
	rawStr := lit.(*ast.StringLit).Val
	s, err := strconv.Unquote(rawStr)
	if err != nil {
		// an invalid string literal raises an error in the escape rules,
		// so simply replace the literal with an empty string here to
		// avoid a cascade of errors.
		s = ""
	}
	m := ast.NewLitMatcher(c.astPos(), s)
	m.IgnoreCase = ignore != nil
	return m, nil
}

func (p *parser) callonLitMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLitMatcher1(stack["lit"], stack["ignore"])
}

func (c *current) onStringLiteral2() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), string(c.text)), nil
}

func (p *parser) callonStringLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral2()
}

func (c *current) onStringLiteral18() (interface{}, error) {
	return ast.NewStringLit(c.astPos(), "``"), errors.New("string literal not terminated")
}

func (p *parser) callonStringLiteral18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral18()
}

func (c *current) onDoubleStringEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonDoubleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleStringEscape5()
}

func (c *current) onSingleStringEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonSingleStringEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSingleStringEscape5()
}

func (c *current) onOctalEscape6() (interface{}, error) {
	return nil, errors.New("invalid octal escape")
}

func (p *parser) callonOctalEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOctalEscape6()
}

func (c *current) onHexEscape6() (interface{}, error) {
	return nil, errors.New("invalid hexadecimal escape")
}

func (p *parser) callonHexEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexEscape6()
}

func (c *current) onLongUnicodeEscape2() (interface{}, error) {
	return validateUnicodeEscape(string(c.text), "invalid Unicode escape")

}

func (p *parser) callonLongUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape2()
}

func (c *current) onLongUnicodeEscape13() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonLongUnicodeEscape13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape13()
}

func (c *current) onShortUnicodeEscape2() (interface{}, error) {
	return validateUnicodeEscape(string(c.text), "invalid Unicode escape")

}

func (p *parser) callonShortUnicodeEscape2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape2()
}

func (c *current) onShortUnicodeEscape9() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonShortUnicodeEscape9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape9()
}

func (c *current) onCharClassMatcher2() (interface{}, error) {
	pos := c.astPos()
	cc := ast.NewCharClassMatcher(pos, string(c.text))
	return cc, nil
}

func (p *parser) callonCharClassMatcher2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher2()
}

func (c *current) onCharClassMatcher15() (interface{}, error) {
	return ast.NewCharClassMatcher(c.astPos(), "[]"), errors.New("character class not terminated")
}

func (p *parser) callonCharClassMatcher15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassMatcher15()
}

func (c *current) onCharClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonCharClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassEscape5()
}

func (c *current) onUnicodeClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape5()
}

func (c *current) onUnicodeClassEscape13(ident interface{}) (interface{}, error) {
	if !unicodeClasses[ident.(*ast.Identifier).Val] {
		return nil, errors.New("invalid Unicode class escape")
	}
	return nil, nil

}

func (p *parser) callonUnicodeClassEscape13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape13(stack["ident"])
}

func (c *current) onUnicodeClassEscape19() (interface{}, error) {
	return nil, errors.New("Unicode class not terminated")

}

func (p *parser) callonUnicodeClassEscape19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape19()
}

func (c *current) onAnyMatcher1() (interface{}, error) {
	any := ast.NewAnyMatcher(c.astPos(), ".")
	return any, nil
}

func (p *parser) callonAnyMatcher1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnyMatcher1()
}

func (c *current) onCodeBlock2() (interface{}, error) {
	pos := c.astPos()
	cb := ast.NewCodeBlock(pos, string(c.text))
	return cb, nil
}

func (p *parser) callonCodeBlock2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock2()
}

func (c *current) onCodeBlock7() (interface{}, error) {
	return nil, errors.New("code block not terminated")
}

func (p *parser) callonCodeBlock7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock7()
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
