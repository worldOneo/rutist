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

type Identifier struct {
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

type FunctionDefinition struct {
	Scope   Node
	ArgList []Identifier
}

type MemberSelector struct {
	Object   Node
	Property Node
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

func (P *Parser) checkAppendage(prev Node) (Node, error) {
	peek, peeked := P.peek()
	if !peeked {
		return prev, nil
	}
	switch peek.Type {
	case tokens.Dot:
		P.next()
		val, err := P.pullValue()
		if err != nil {
			return nil, err
		}
		return P.checkAppendage(MemberSelector{prev, val})
	}
	return prev, nil
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

func (P *Parser) argList(identifierOnly bool) ([]Node, error) {
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
	v, err := P._pullValue()
	if err != nil {
		return nil, err
	}
	return P.checkAppendage(v)
}

func (P *Parser) _pullValue() (Node, error) {
	next, has := P.peek()
	if !has {
		return nil, fmt.Errorf("Expected value")
	}

	switch next.Type {
	case tokens.ScopeOpen:
		b, err := P.parse()
		if err != nil {
			return nil, err
		}
		return Scope{b}, nil
	case tokens.ParenOpen:
		P.next()
		args, err := P.argList(true)
		if err != nil {
			return nil, err
		}
		arglist := make([]Identifier, len(args))
		for i := 0; i < len(args); i++ {
			arg, ok := args[i].(Identifier)
			if !ok {
				return nil, fmt.Errorf("Identifier expected line: %d", next.Line)
			}
			arglist[i] = arg
		}
		b, err := P.parse()
		if err != nil {
			return nil, err
		}
		return FunctionDefinition{b, arglist}, nil
	case tokens.Identifier:
		P.next()
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
			args, err := P.argList(false)
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
		return Identifier{next.Content}, nil
	case tokens.Float:
		P.next()
		return Float{next.ValueFloat}, nil
	case tokens.Integer:
		P.next()
		return Int{next.ValueInt}, nil
	case tokens.String:
		P.next()
		return String{next.Content}, nil
	case tokens.Boolean:
		P.next()
		return Bool{next.ValueInt == 1}, nil
	}
	return nil, fmt.Errorf("Identifier Expected")
}

func (P *Parser) parseIdentifier(last tokens.Token) (Node, error) {
	peek, peeked := P.peek()
	if !peeked {
		return Identifier{last.Content}, nil
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
		return MemberSelector{Identifier{last.Content}, node}, nil
	}
	return Identifier{last.Content}, nil
}
