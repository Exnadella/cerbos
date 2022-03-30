// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func NewParser() (*participle.Parser, error) {
	lex, err := lexer.New(lexer.Rules{
		"Root": {
			{Name: "Ident", Pattern: `[a-zA-Z]\w*(\.\w+)*`},
			{Name: "Assign", Pattern: `=`, Action: lexer.Push("Assign")},
			{Name: "Whitespace", Pattern: `\s+`},
		},
		"Assign": {
			{Name: "Any", Pattern: `.*$`, Action: lexer.Pop()},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create lexer: %w", err)
	}

	parser, err := participle.Build(&REPLDirective{},
		participle.Lexer(lex),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser: %w", err)
	}

	return parser, nil
}

//nolint:govet
type REPLDirective struct {
	Exit  bool          `parser:"@('q'|'quit'|'exit')"`
	Reset bool          `parser:"| @'reset'"`
	Vars  bool          `parser:"| @'vars'"`
	Help  bool          `parser:"| @('h' | 'help')"`
	Let   *LetDirective `parser:"| @@"`
}

type LetDirective struct {
	Name string `parser:"'let' @Ident"`
	Expr string `parser:"'=' @Any"`
}
