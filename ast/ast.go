package ast

import (
	"fmt"

	"github.com/worldOneo/rutist/tokens"
)

type Type uint64

const (
	TypeCall Type = iota
	TypeVariable
)

type Node interface{}

type Variable struct {
	Name string
}

type String struct {
	Value string
}

type Float struct {
	Value float64
}

type Int struct {
	Value int
}

type Bool struct {
	Value bool
}

func (B Bool) Bool() Bool {
	return Bool{B.Value}
}

func (B Int) Bool() Bool {
	return Bool{B.Value != 0}
}

func (B Float) Bool() Bool {
	return Bool{B.Value != 0}
}

func (B String) Bool() Bool {
	return Bool{B.Value != ""}
}

type Block struct {
	Body []Node
}

type Assignment struct {
	Identifier string
	Value      Node
}

type Expression struct {
	Identifier string
	ArgList    []Node
}

type Program = Block

type Parser struct {
	tokens []tokens.Token
	index  int
}

func Parse(lexed []tokens.Token) (Node, error) {
	parser := Parser{
		tokens: lexed,
	}
	return parser.parse()
}

func (p *Parser) parse() (Node, error) {
	l := 64
	body := make([]Node, l)
	bindex := 0
	for p.index < len(p.tokens) {
		node, err := p.pullValue()
		if err != nil {
			return nil, err
		}
		body[bindex] = node
		bindex++
		if bindex >= l {
			old := body
			l *= 2
			body = make([]Node, l)
			copy(body, old)
		}
	}
	return Block{body[0:bindex]}, nil
}

func (P *Parser) peek() (tokens.Token, bool) {
	if P.index < len(P.tokens) {
		return P.tokens[P.index], true
	}
	return tokens.Token{}, false
}

func (P *Parser) next() (tokens.Token, bool) {
	if P.index < len(P.tokens) {
		P.index++
		return P.tokens[P.index-1], true
	}
	return tokens.Token{}, false
}

func (P *Parser) argList() ([]Node, error) {
	args := make([]Node, 0)
	requiresComma := false

	for peek, peeked := P.peek(); peeked && peek.Type != tokens.ParenClosed; peek, peeked = P.peek() {
		if requiresComma && peek.Type == tokens.Comma {
			requiresComma = false
			P.next()
			continue
		} else if requiresComma || peek.Type == tokens.Comma {
			return nil, fmt.Errorf("unexpected comma")
		}
		arg, err := P.pullValue()
		requiresComma = true
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	P.next()
	return args, nil
}

func (P *Parser) pullValue() (Node, error) {
	next, has := P.next()
	if !has {
		return nil, fmt.Errorf("Expected value")
	}

	peek, _ := P.peek()
	switch next.Type {
	case tokens.Identifier:
		if peek.Type == tokens.ParenOpen {
			P.next()
			args, err := P.argList()
			if err != nil {
				return nil, err
			}
			return Expression{next.Content, args}, nil
		} else if peek.Type == tokens.ScopeOpen {
			P.next()
			return P.parse()
		} else if peek.Type == tokens.Assignment {
			P.next()
			node, err := P.pullValue()
			if err != nil {
				return nil, err
			}
			return Assignment{next.Content, node}, nil
		}
		return String{next.Content}, nil
	case tokens.Float:
		return Float{next.ValueFloat}, nil
	case tokens.Integer:
		return Int{next.ValueInt}, nil
	case tokens.String:
		return String{next.Content}, nil
	case tokens.Boolean:
		return Bool{next.ValueInt == 1}, nil
	}
	return nil, fmt.Errorf("Identifier Expected")
}
