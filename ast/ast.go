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

type Scope struct {
	Body Node
}

type Assignment struct {
	Identifier Node
	Value      Node
}

type Expression struct {
	Callee  Node
	ArgList []Node
}

type MemberSelector struct {
	Identifier string
	Property   Node
}

type Program = Block

type Parser struct {
	tokens []tokens.Token
	index  int
}

func Parsep(lexed []tokens.Token) Node {
	val, err := Parse(lexed)
	if err != nil {
		panic(err)
	}
	return val
}

func Parse(lexed []tokens.Token) (Node, error) {
	parser := Parser{
		tokens: lexed,
	}
	return parser.parse()
}

func (P *Parser) parse() (Node, error) {
	l := 64
	body := make([]Node, l)
	bindex := 0
	peek, peeked := P.peek()
	returnOnScopeClose := peeked && peek.Type == tokens.ScopeOpen

	if returnOnScopeClose {
		P.next()
	}

	for P.index < len(P.tokens) {
		if returnOnScopeClose {
			peek, peeked = P.peek()
			if peeked && peek.Type == tokens.ScopeClosed {
				P.next()
				break
			}
		}
		node, err := P.pullValue()
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
		} else if requiresComma {
			if peek.Type == tokens.ParenClosed {
				return args, nil
			}
			return nil, fmt.Errorf("expected comma line %d", peek.Line)
		} else if peek.Type == tokens.Comma {
			return nil, fmt.Errorf("Unexpected comma line %d", peek.Line)
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
	case tokens.Scoper:
		if peek.Type == tokens.ScopeOpen {
			b, err := P.parse()
			if err != nil {
				return nil, err
			}
			return Scope{b}, nil
		}
	case tokens.Identifier:
		identifier, err := P.parseIdentifier(next)
		if err != nil {
			return nil, err
		}
		peek, peeked := P.peek()
		if !peeked {
			return identifier, nil
		}
		switch peek.Type {
		case tokens.ParenOpen:
			P.next()
			args, err := P.argList()
			if err != nil {
				return nil, err
			}
			return Expression{identifier, args}, nil
		case tokens.Assignment:
			P.next()
			node, err := P.pullValue()
			if err != nil {
				return nil, err
			}
			return Assignment{identifier, node}, nil
		}
		return Variable{next.Content}, nil
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

func (P *Parser) parseIdentifier(last tokens.Token) (Node, error) {
	peek, peeked := P.peek()
	if !peeked {
		return Variable{last.Content}, nil
	}
	switch peek.Type {
	case tokens.Dot:
		P.next()
		current, has := P.next()
		if !has {
			return nil, fmt.Errorf("Identifier expected")
		}
		node, err := P.parseIdentifier(current)
		if err != nil {
			return nil, err
		}
		return MemberSelector{last.Content, node}, nil
	}
	return Variable{last.Content}, nil
}
