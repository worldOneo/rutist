package ast

import (
	"fmt"

	"github.com/worldOneo/rutist/tokens"
)

type Type uint64

type Program = Block

type Parser struct {
	tokens []tokens.Token
	index  int
	file string
}

func (P *Parser) Meta(t tokens.Token) *Meta {
	return NewMeta(t, P.file)
}

func Parsep(lexed []tokens.Token) Node {
	val, err := Parse(lexed, "constant.go")
	if err != nil {
		panic(err)
	}
	return val
}

func Parse(lexed []tokens.Token, file string) (Node, error) {
	parser := Parser{
		tokens: lexed,
		file: file,
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
	return Block{body[0:bindex], P.Meta(peek)}, nil
}

func (P *Parser) checkAppendage(prev Node) (Node, error) {
	peek, peeked := P.peek()
	if !peeked {
		return prev, nil
	}
	switch peek.Type {
	case tokens.Dot:
		P.next()
		val, err := P._pullValue()
		if err != nil {
			return nil, err
		}
		return P.checkAppendage(MemberSelector{prev, val, P.Meta(peek)})
	case tokens.ParenOpen:
		P.next()
		args, err := P.argList(false)
		if err != nil {
			return nil, err
		}
		return P.checkAppendage(Expression{prev, args, P.Meta(peek)})
	case tokens.Assignment:
		P.next()
		node, err := P.pullValue()
		if err != nil {
			return nil, err
		}
		return Assignment{prev, node, P.Meta(peek)}, nil
	case tokens.OperatorType:
		P.next()
		node, err := P.pullValue()
		if err != nil {
			return nil, err
		}
		if bin, ok := node.(BinaryExpression); ok {
			if bin.Operation > peek.ValueInt {
				expr := BinaryExpression{peek.ValueInt, prev, bin.Left, P.Meta(peek)}
				bin.Left = expr
				return bin, nil
			}
		}
		return BinaryExpression{peek.ValueInt, prev, node, P.Meta(peek)}, nil
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
		return Scope{b, P.Meta(next)}, nil
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
		return FunctionDefinition{b, arglist, P.Meta(next)}, nil
	case tokens.Identifier:
		P.next()
		identifier, err := P.parseIdentifier(next)
		if err != nil {
			return nil, err
		}
		return identifier, nil
	case tokens.Float:
		P.next()
		return Float{next.ValueFloat, P.Meta(next)}, nil
	case tokens.Integer:
		P.next()
		return Int{next.ValueInt, P.Meta(next)}, nil
	case tokens.String:
		P.next()
		return String{next.Content, P.Meta(next)}, nil
	case tokens.Boolean:
		P.next()
		return Bool{next.ValueInt == 1, P.Meta(next)}, nil
	case tokens.OperatorType:
		P.next()
		val, err := P.pullValue()
		if err != nil {
			return nil, err
		}
		return UnaryExpression{next.ValueInt, val, P.Meta(next)}, nil
	}
	return nil, fmt.Errorf("Identifier Expected")
}

func (P *Parser) parseIdentifier(last tokens.Token) (Node, error) {
	peek, peeked := P.peek()
	if !peeked {
		return Identifier{last.Content, P.Meta(last)}, nil
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
		return MemberSelector{Identifier{last.Content, P.Meta(last)}, node, P.Meta(peek)}, nil
	}
	return Identifier{last.Content, P.Meta(last)}, nil
}
