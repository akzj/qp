package qp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type lexer struct {
	reader *bufio.Reader
	line   int
	token  Token
	err    error
}

func (l *lexer) finish() bool {
	return l.err != nil && l.token.typ == EOFTokenType
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func (l *lexer) ahead() (byte, error) {
	c, err := l.reader.ReadByte()
	if err != nil {
		l.err = err
		return 0, err
	}
	_ = l.reader.UnreadByte()
	return c, nil
}

func (l *lexer) get() (byte, error) {
	c, err := l.reader.ReadByte()
	if err != nil {
		l.err = err
		return 0, err
	}
	return c, nil
}

func (l *lexer) peek() Token {
	if l.token.typ != EOFTokenType {
		return l.token
	}
	for {
		c, err := l.get()
		if err != nil {
			return emptyToken
		}
		var token Token
		switch {
		case isSpace(c):
			if c == '\n' {
				l.line++
			}
			continue
		case isLetter(c):
			token = l.parseLabel(c)
		case c == '+':
			if a, _ := l.ahead(); a == '+' {
				_, _ = l.get()
				token = incOperatorToken
			} else {
				token = addOperatorToken
			}
		case c == '(':
			token = leftParenthesisToken
		case c == ')':
			token = rightParenthesisToken
		case c == '{':
			token = leftBraceToken
		case c == '}':
			fmt.Println(c)
			token = rightBraceToken
		case c == '<':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.get()
				token = lessEqualToken
			} else {
				token = lessToken
			}
		case c == '>':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.get()
				token = greaterEqualToken
			} else {
				token = greaterToken
			}
		case c == '*':
			token = mulOperatorToken
		case '0' <= c && c <= '9':
			token = l.parseNumToken(c)
		case c == '=':
			token = assignToken
		case c == ',':
			token = commaToken
		default:
			token = unknownToken
		}
		l.token = token
		return token
	}
}

func (l *lexer) parseNumToken(c byte) Token {
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.reader.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if isDigit(c) {
			buf.WriteByte(c)
		} else {
			if err := l.reader.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	return Token{
		typ: intTokenType,
		val: buf.String(),
	}
}

func (l *lexer) next() {
	l.token = emptyToken
}

func (l *lexer) parseLabel(c byte) Token {
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.reader.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if isLetter(c) {
			buf.WriteByte(c)
		} else {
			if err := l.reader.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	for _, keyword := range Keywords {
		if keyword == buf.String() {
			return Token{
				typ: keywordTokenType[keyword],
			}
		}
	}

	return Token{
		typ: labelType,
		val: buf.String(),
	}
}

func (l *lexer) Line() int {
	return l.line
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{
		token:  emptyToken,
		reader: bufio.NewReader(reader),
	}
}
