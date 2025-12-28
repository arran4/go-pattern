package dsl

import (
	"unicode"
)

type TokenType int

const (
	EOF TokenType = iota
	IDENT
	STRING
	NUMBER
	PIPE   // |
	LPAREN // (
	RPAREN // )
	EQUALS // =
	CARET  // ^
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += 1
}

// peekChar removed as unused

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token

	switch l.ch {
	case '|':
		tok = newToken(PIPE, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '=':
		tok = newToken(EQUALS, l.ch)
	case '^':
		tok = newToken(CARET, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = IDENT
			return tok
		} else if isDigit(l.ch) || l.ch == '-' || l.ch == '#' || l.ch == '.' {
			// Numbers, Colors (#...), Decimals, Negatives
			tok.Literal = l.readNumberOrColor()
			tok.Type = NUMBER // Treating loosely as a value
			return tok
		} else {
			tok = newToken(EOF, l.ch) // Unknown
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readNumberOrColor() string {
	position := l.pos
	// Allow characters for numbers, hex colors, and simple strings that start with these
	for isDigit(l.ch) || isLetter(l.ch) || l.ch == '.' || l.ch == '-' || l.ch == '#' {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readString() string {
	position := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
