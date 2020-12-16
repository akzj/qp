package qp

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type Lexer struct {
	reader *bufio.Reader
	line   int
	token  Token
	err    error
}

func (l *Lexer) Finish() bool {
	return l.err != nil && l.token.typ == EOFType
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func (l *Lexer) ahead() (byte, error) {
	c, err := l.reader.ReadByte()
	if err != nil {
		l.err = err
		return 0, err
	}
	_ = l.reader.UnreadByte()
	return c, nil
}

func (l *Lexer) Get() (byte, error) {
	c, err := l.reader.ReadByte()
	if err != nil {
		l.err = err
		return 0, err
	}
	return c, nil
}

func (l *Lexer) Peek() Token {
	if l.token.typ != EOFType {
		return l.token
	}
	for {
		c, err := l.Get()
		if err != nil {
			return EmptyToken
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
				_, _ = l.Get()
				token = IncOperatorToken
			} else {
				token = AddOperatorToken
			}
		case c == '(':
			token = LeftParenthesisToken
		case c == ')':
			token = RightParenthesisToken
		case c == '{':
			token = LeftBraceToken
		case c == '}':
			token = RightBraceToken
		case c == '[':
			token = LeftBracketToken
		case c == ']':
			token = RightBracketToken
		case c == '<':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.Get()
				token = LessEqualToken
			} else {
				token = LessToken
			}
		case c == '>':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.Get()
				token = GreaterEqualToken
			} else {
				token = GreaterToken
			}
		case c == '*':
			token = MulOperatorToken
		case '0' <= c && c <= '9':
			token = l.parseNumToken(c)
		case c == '!':
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.Get()
				token = NoEqualToken
			} else {
				log.Panicf("unknown token`%s`", string(c))
			}
		case c == '-':
			token = SubOperatorToken
		case c == '=':
			token = AssignToken
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.Get()
				token = EqualToken
			}
		case c == ',':
			token = CommaToken
		case c == ';':
			token = SemicolonToken
		case c == ':':
			token = ColonToken
		case c == '.':
			token = PeriodToken
		case c == '"':
			token = l.parseString(false)
		case c == '`':
			token = l.parseString(true)
		case c == '/':
			if ahead, _ := l.ahead(); ahead == '/' { //
				_, _ = l.Get()
				token = Token{
					typ:  CommentType,
					val:  l.readline(),
					line: l.line,
				}
			}
		case c == '|':
			if c, _ := l.ahead(); c == '|' {
				l.Get()
				token = OrToken
			}
		case c == '&':
			if c, _ := l.ahead(); c == '&' {
				l.Get()
				token = AndToken
			}
		default:
			log.Panicln(string(c), l.line)
		}
		token.line = l.line
		l.token = token
		return token
	}
}

func (l *Lexer) parseString(multiline bool) Token {
	var buffer bytes.Buffer
	for {
		c, err := l.Get()
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
		typ:  StringType,
		val:  buffer.String(),
		line: 0,
	}
}

func (l *Lexer) parseNumToken(c byte) Token {
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
		typ: IntType,
		val: buf.String(),
	}
}

func (l *Lexer) next() {
	l.token = EmptyToken
}

func (l *Lexer) parseLabel(c byte) Token {
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
				typ: KeywordType[keyword],
			}
		}
	}

	return Token{
		typ: IDType,
		val: buf.String(),
	}
}

func (l *Lexer) Line() int {
	return l.line
}

func (l *Lexer) readline() string {
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
		_, _ = l.Get()
		line = append(line, c)
	}
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func newLexer(reader io.Reader) *Lexer {
	return &Lexer{
		token:  EmptyToken,
		reader: bufio.NewReader(reader),
	}
}
