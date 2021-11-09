package tokens

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType = uint32
type Operator = int

type Token struct {
	Type       TokenType
	Content    string
	ValueInt   int
	ValueFloat float64
	Line       int
}

const (
	Identifier TokenType = iota
	String
	Integer
	Float
	ScopeOpen
	ScopeClosed
	ParenOpen
	ParenClosed
	Comma
	Assignment
	Boolean
	Scoper
	Dot
	OperatorType
)

const (
	OperatorAdd Operator = iota
	OperatorSub
	OperatorMul
	OperatorDiv
	OperatorMod
	OperatorOr
	OperatorAnd
	OperatorXor
	OperatorLor
	OperatorLand
	OperatorNot
	OperatorEq
	OperatorLt
	OperatorLe
	OperatorGt
	OperatorGe
	OperatorLsh
	OperatorRsh
)

var operators = map[string]Operator{}

func init() {
	operators["+"] = OperatorAdd
	operators["-"] = OperatorSub
	operators["*"] = OperatorMul
	operators["/"] = OperatorDiv
	operators["%"] = OperatorMod
	operators["|"] = OperatorOr
	operators["&"] = OperatorAnd
	operators["^"] = OperatorXor
	operators["||"] = OperatorLor
	operators["&&"] = OperatorLand
	operators["~"] = OperatorNot
	operators["!"] = OperatorNot
	operators["=="] = OperatorEq
	operators["<"] = OperatorLt
	operators["<="] = OperatorLe
	operators[">"] = OperatorGt
	operators[">="] = OperatorGe
	operators["<<"] = OperatorLsh
	operators[">>"] = OperatorRsh
}

const windowsLineSpererator = "\r\n"
const commentIntroduction = "//"

type CodeLexer struct {
	code        []rune
	words       []Token
	currentWord int
}

func (C *CodeLexer) append(word Token) {
	C.words[C.currentWord] = word
	C.currentWord++
	if C.currentWord >= len(C.words) {
		old := C.words
		C.words = make([]Token, len(C.words)*2)
		copy(C.words, old)
	}
}

func Peek(runes []rune, index int) (rune, bool) {
	if index < len(runes) {
		return runes[index], true
	}
	return 0, false
}

func Lexerp(code string) []Token {
	a, err := Lexer(code)
	if err != nil {
		panic(err)
	}
	return a
}

func Lexer(code string) ([]Token, error) {
	parser := CodeLexer{
		[]rune(code),
		make([]Token, 64),
		0,
	}
	words, err := parser.Lexer()
	if err != nil {
		return words, err
	}
	return words[0:parser.currentWord], err
}

func (C *CodeLexer) Lexer() ([]Token, error) {
	lineComment := false
	line := 0
	buff := strings.Builder{}
	for i := 0; i < len(C.code); i++ {
		c := C.code[i]
		safeInc := func() (rune, bool) {
			i++
			if i < len(C.code) {
				return C.code[i], true
			}
			return 0, false
		}
		n, peeked := Peek(C.code, i+1)
		if isNewLine(c) {
			if isNewLine(n) {
				i++
			}
			line++
			lineComment = false
			continue
		}

		if isSpace(c) || lineComment {
			continue
		}

		if isLineComment(c, n) {
			lineComment = true
		}

		if isSpecialChar(c) {
			switch c {
			case '{':
				C.append(scopeOpenToken(line))
			case '}':
				C.append(scopeClosedToken(line))
			case '(':
				C.append(Token{ParenOpen, "(", 0, 0, line})
			case ')':
				C.append(Token{ParenClosed, ")", 0, 0, line})
			case ',':
				C.append(Token{Comma, ",", 0, 0, line})
			case '+', '-', '/', '*', '%', '=', '>', '<', '~', '!', '|', '&':
				sign := string(c)
				if isSpecialChar(n) {
					sign += string(n)
				}
				if sign == "=" {
					C.append(Token{Assignment, "=", 0, 0, line})
					continue
				}

				var ok bool
				var operator Operator
				if operator, ok = operators[sign]; ok {
					C.append(Token{OperatorType, sign, operator, 0, line})
				} else {
					return nil, fmt.Errorf("Invalid operator at line %d", line)
				}
				if ok {
					i += len(sign) - 1
				}
			case '@':
				C.append(Token{Scoper, "@", 0, 0, line})
			case '.':
				C.append(Token{Dot, ".", 0, 0, line})
			}
			continue
		}

		if isAlpha(c) {
			buff.Reset()
			var inc bool
			for isAlpha(c) {
				buff.WriteRune(c)
				c, inc = safeInc()
				if !inc {
					break
				}
			}
			i--
			val := buff.String()
			switch val {
			case "true":
				C.append(Token{Boolean, "true", 1, 0, line})
			case "false":
				C.append(Token{Boolean, "false", 0, 0, line})
			default:
				C.append(Token{Identifier, val, 0, 0, line})
			}
			continue
		}

		if isStringBegin(c) {
			buff.Reset()
			escaped := false
			i++
			for {
				c := C.code[i]
				if isStringBegin(c) {
					break
				}
				if escaped {
					buff.WriteRune(getEscapedCharacter(c))
					escaped = false
					i++
					continue
				}
				if isEscapeChar(c) {
					escaped = true
					i++
					continue
				}
				buff.WriteRune(c)
				n, peeked = Peek(C.code, i+1)
				if !peeked {
					return []Token{}, fmt.Errorf("Incomplete string")
				}
				i++
			}
			C.append(Token{String, buff.String(), 0, 0, line})
			continue
		}

		if isDigit(c) {
			buff.Reset()
			float := false
			for isDigit(c) || isNumericalSkipChar(c) || c == '.' {
				if isNumericalSkipChar(c) {
					continue
				}
				if c == '.' {
					float = true
				}
				buff.WriteRune(c)
				i++
				c = C.code[i]
			}
			i--
			str := buff.String()
			if !float {
				intVal, err := strconv.Atoi(str)
				if err != nil {
					return []Token{}, fmt.Errorf("Unparseble int literal")
				}
				C.append(intToken(str, intVal, line))
				continue
			}
			floatVal, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return []Token{}, fmt.Errorf("Unparseble float literal")
			}
			C.append(floatToken(str, floatVal, line))
			continue
		}
	}
	return C.words, nil
}

func isNumericalSkipChar(b rune) bool {
	return b == '_'
}

func isAlpha(b rune) bool {
	return !isDigit(b) && !isSpecialChar(b) && !isSpace(b) && !isStringBegin(b)
}

func isDigit(b rune) bool {
	return b >= '0' && b <= '9'
}

func isStringBegin(b rune) bool {
	return b == '"'
}

func isSpecialChar(b rune) bool {
	return b == '{' || b == '}' ||
		b == '(' || b == ')' ||
		b == ',' || b == '.' ||
		b == '<' || b == '>' ||
		b == '~' || b == '!' ||
		b == '|' || b == '&' ||
		b == '+' || b == '-' ||
		b == '*' || b == '/' ||
		b == '%' ||
		isEqual(b) || isScoper(b)
}

func isSpace(b rune) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isNewLine(b rune) bool {
	return b == '\n' || b == '\r'
}

func isLineComment(b rune, c rune) bool {
	return b == c && b == '/'
}

func isEscapeChar(b rune) bool {
	return b == '\\'
}

func isEqual(b rune) bool {
	return b == '='
}

func isScoper(b rune) bool {
	return b == '@'
}

func getEscapedCharacter(b rune) rune {
	switch b {
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case '"':
		return '"'
	default:
		return b
	}
}

func scopeOpenToken(line int) Token {
	return Token{ScopeOpen, "{", 0, 0, line}
}

func scopeClosedToken(line int) Token {
	return Token{ScopeClosed, "}", 0, 0, line}
}

func intToken(str string, val int, line int) Token {
	return Token{Integer, str, val, 0, line}
}

func floatToken(str string, val float64, line int) Token {
	return Token{Float, str, 0, val, line}
}

func stringToken(content string, line int) Token {
	return Token{String, content, 0, 0, line}
}

func identifierToken(id string, line int) Token {
	return Token{Identifier, id, 0, 0, line}
}
