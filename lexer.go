package qp

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type lexer struct {
	reader *bufio.Reader
	line   int
	token  Token
	err    error
}

func (l *lexer) finish() bool {
	return l.err != nil && l.token.typ == EOFType
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
	if l.token.typ != EOFType {
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
			token = rightBraceToken
		case c == '[':
			token = leftBracketToken
		case c == ']':
			token = rightBracketToken
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
		case c == '!':
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.get()
				token = NoEqualToken
			} else {
				log.Panicf("unknown token`%s`", string(c))
			}
		case c == '-':
			token = subOperatorToken
		case c == '=':
			token = assignToken
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.get()
				token = equalToken
			}
		case c == ',':
			token = commaToken
		case c == ';':
			token = semicolonToken
		case c == ':':
			token = colonToken
		case c == '.':
			token = periodToken
		case c == '"':
			token = l.parseString(false)
		case c == '`':
			token = l.parseString(true)
		case c == '/':
			if ahead, _ := l.ahead(); ahead == '/' { //
				_, _ = l.get()
				token = Token{
					typ:  commentTokenType,
					val:  l.readline(),
					line: l.line,
				}
			}
		case c == '|':
			if c, _ := l.ahead(); c == '|' {
				l.get()
				token = orToken
			}
		case c == '&':
			if c, _ := l.ahead(); c == '&' {
				l.get()
				token = andToken
			}
		default:
			log.Panicln(string(c), l.line)
		}
		token.line = l.line
		l.token = token
		return token
	}
}

func (l *lexer) parseString(multiline bool) Token {
	var buffer bytes.Buffer
	for {
		c, err := l.get()
		if err != nil {
			l.err = err
			break
		}
		if c == '\n' {
			if multiline == false {
				log.Panic("parse string failed", string(c))
			} else {
				l.line++
			}
		}
		if c == '\\' {
			c, err = l.ahead()
			if err != nil {
				l.err = err
				break
			}
			if c == '"' {
				buffer.WriteByte('"')
				continue
			}
		}
		if c == '"' && multiline == false {
			break
		}
		if c == '`' && multiline {
			break
		}
		buffer.WriteByte(c)
	}
	return Token{
		typ:  stringType,
		val:  buffer.String(),
		line: 0,
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
		typ: intType,
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
		if isLetter(c) || isDigit(c) || c == '_' {
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
		typ: IDType,
		val: buf.String(),
	}
}

func (l *lexer) Line() int {
	return l.line
}

func (l *lexer) readline() string {
	var line []byte
	for {
		c, err := l.ahead()
		if err != nil {
			l.err = err
			return string(line)
		}
		if c == '\n' {
			return string(line)
		}
		_, _ = l.get()
		line = append(line, c)
	}
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
