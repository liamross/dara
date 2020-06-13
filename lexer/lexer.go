package lexer

import (
	"dara/token"
	"fmt"
	"strings"
)

// Lexer allows you to extract tokens from a dara input string.
type Lexer struct {
	input    string
	line     int
	position int
	ch       byte
}

// New returns a primed `Lexer`.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, position: -1}
	l.advance()
	return l
}

func (l *Lexer) Position() int {
	return l.line
}

// Scan scans the entire input and returns all of the tokens.
func (l *Lexer) Scan() []*token.Token {
	var tokens []*token.Token

	for !l.isAtEnd() {
		tokens = append(tokens, l.NextToken())
	}
	tokens = append(tokens, l.NextToken())

	return tokens
}

// NextToken scans through and returns individual tokens.
func (l *Lexer) NextToken() *token.Token {
	var tok *token.Token

	l.skipWhitespace()

	switch l.ch {
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)

	case '<':
		tok = l.lookAhead('=', token.LT_EQ, token.LT)
	case '>':
		tok = l.lookAhead('=', token.GT_EQ, token.GT)
	case '=':
		tok = l.lookAhead('=', token.EQ, token.ASSIGN)
	case '!':
		tok = l.lookAhead('=', token.NOT_EQ, token.BANG)
	case '&':
		if l.peek() == '&' {
			l.advance()
			tok = token.New(token.AND, "&&")
		}
	case '|':
		if l.peek() == '|' {
			l.advance()
			tok = token.New(token.OR, "||")
		}

	case '/':
		switch l.peek() {
		case '/':
			return token.New(token.COMMENT, l.lineComment())
		case '*':
			return token.New(token.COMMENT, l.blockComment())
		default:
			tok = newToken(token.SLASH, l.ch)
		}

	case '"':
		return token.New(token.STRING, l.string('"'))
	case '\'':
		return token.New(token.STRING, l.string('\''))

	case 0:
		tok = token.New(token.EOF, "")
	default:
		switch {
		case isAlpha(l.ch):
			return l.readIdentifier()
		case isDigit(l.ch):
			return l.readNumber()
		default:
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.advance()
	}
}

func (l *Lexer) advance() {
	if l.position+1 >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position+1]
	}
	l.position++
}

func (l *Lexer) peek() byte {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return l.input[l.position+1]
}

func (l *Lexer) lookAhead(check byte, a, b token.TokenType) *token.Token {
	if l.peek() == check {
		ch := l.ch
		l.advance()
		return token.New(a, string(ch)+string(l.ch))
	}
	return newToken(b, l.ch)
}

func (l *Lexer) readIdentifier() *token.Token {
	start := l.position
	for isAlphaNumeric(l.ch) {
		l.advance()
	}
	literal := l.input[start:l.position]
	return token.New(token.LookupIdent(literal), literal)
}

func (l *Lexer) readNumber() *token.Token {
	start := l.position
	for isDigit(l.ch) {
		l.advance()
	}
	if l.ch == '.' && isDigit(l.peek()) {
		l.advance()
		for isDigit(l.ch) {
			l.advance()
		}
	}
	return token.New(token.NUMBER, l.input[start:l.position])
}

func (l *Lexer) string(end byte) string {
	l.advance()
	start := l.position

	for l.ch != end && !l.isAtEnd() {
		if l.ch == '\n' {
			l.line++
		}
		l.advance()
	}

	literal := l.input[start:l.position]

	if l.isAtEnd() {
		// TODO: Unterminated string error.
		fmt.Println("Unterminated string error")
		return literal
	}

	l.advance()
	return literal
}

func (l *Lexer) lineComment() string {
	l.advance()
	l.advance()
	start := l.position

	for l.ch != '\n' && !l.isAtEnd() {
		l.advance()
	}

	return strings.TrimSpace(l.input[start:l.position])
}

func (l *Lexer) blockComment() string {
	l.advance()
	l.advance()
	start := l.position

	for !(l.ch == '*' && l.peek() == '/') && !l.isAtEnd() {
		if l.ch == '\n' {
			l.line++
		}
		l.advance()
	}

	literal := strings.TrimSpace(l.input[start:l.position])

	if l.isAtEnd() {
		// TODO: Unterminated block comment error.
		fmt.Println("Unterminated block comment error")
		return literal
	}

	l.advance()
	l.advance()
	return literal
}

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func newToken(tokenType token.TokenType, ch byte) *token.Token {
	return token.New(tokenType, string(ch))
}
