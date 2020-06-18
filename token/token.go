package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal string) *Token {
	return &Token{tokenType, literal}
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	COMMENT TokenType = "COMMENT"

	// Identifiers and literals.
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	// Operators.
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	LT       TokenType = "<"
	GT       TokenType = ">"
	MOD      TokenType = "%"

	EQ      TokenType = "=="
	NOT_EQ  TokenType = "!="
	GT_EQ   TokenType = ">="
	LT_EQ   TokenType = "<="
	AND     TokenType = "&&"
	OR      TokenType = "||"
	DECLARE TokenType = ":="

	// Delimiters.
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// Keywords.
	FUNCTION TokenType = "fn"
	TRUE     TokenType = "true"
	FALSE    TokenType = "false"
	IF       TokenType = "if"
	ELSE     TokenType = "else"
	RETURN   TokenType = "return"
	NIL      TokenType = "nil"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"nil":    NIL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
