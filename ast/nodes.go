package ast

import (
	"github.com/worldOneo/rutist/tokens"
)

type Node interface {
	Token() tokens.Token
	File() string
	SetToken(tokens.Token)
}

type Meta struct {
	At tokens.Token
	F  string
}

func (M Meta) Token() tokens.Token {
	return M.At
}

func (M Meta) File() string {
	return M.F
}

func (M *Meta) SetToken(t tokens.Token) {
	M.At = t
}

func NewMeta(t tokens.Token, file string) *Meta {
	return &Meta{t, file}
}

type Identifier struct {
	Name string
	*Meta
}

type String struct {
	Value string
	*Meta
}

type Float struct {
	Value float64
	*Meta
}

type Int struct {
	Value int
	*Meta
}

type Bool struct {
	Value bool
	*Meta
}
type Block struct {
	Body []Node
	*Meta
}

type Scope struct {
	Body Node
	*Meta
}

type Assignment struct {
	Identifier Node
	Value      Node
	*Meta
}

type Expression struct {
	Callee  Node
	ArgList []Node
	*Meta
}

type FunctionDefinition struct {
	Scope   Node
	ArgList []Identifier
	*Meta
}

type MemberSelector struct {
	Object   Node
	Property Node
	*Meta
}

type UnaryExpression struct {
	Operation tokens.Operator
	Value     Node
	*Meta
}

type BinaryExpression struct {
	Operation tokens.Operator
	Left      Node
	Right     Node
	*Meta
}

func walkTree(tree Node, f func(node Node)) {
	f(tree)
	switch n := tree.(type) {
	case Identifier, Float, String, Int, Bool:
		return
	case Block:
		for _, n := range n.Body {
			walkTree(n, f)
		}
	case Scope:
		walkTree(n.Body, f)
	case Assignment:
		walkTree(n.Identifier, f)
		walkTree(n.Value, f)
	case Expression:
		walkTree(n.Callee, f)
		for _, n := range n.ArgList {
			walkTree(n, f)
		}
	case FunctionDefinition:
		walkTree(n.Scope, f)
		for _, n := range n.ArgList {
			walkTree(n, f)
		}
	case MemberSelector:
		walkTree(n.Object, f)
		walkTree(n.Property, f)
	case BinaryExpression:
		walkTree(n.Left, f)
		walkTree(n.Right, f)
	case UnaryExpression:
		walkTree(n.Value, f)
	}
}
