package lexer

import (
	"bufio"
	"bytes"
	"gitlab.com/akzj/qp"
	"io"
	"log"
)

type Lexer struct {
	reader *bufio.Reader
	line   int
	token  qp.Token
	err    error
}

func (l *Lexer) Finish() bool {
	return l.err != nil && l.token.Typ == qp.EOFType
}

func IsSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func IsLetter(c byte) bool {
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

func (l *Lexer) Peek() qp.Token {
	if l.token.Typ != qp.EOFType {
		return l.token
	}
	for {
		c, err := l.Get()
		if err != nil {
			return qp.EmptyToken
		}
		var token qp.Token
		switch {
		case IsSpace(c):
			if c == '\n' {
				l.line++
			}
			continue
		case IsLetter(c):
			token = l.parseLabel(c)
		case c == '+':
			if a, _ := l.ahead(); a == '+' {
				_, _ = l.Get()
				token = qp.IncOperatorToken
			} else {
				token = qp.AddOperatorToken
			}
		case c == '(':
			token = qp.LeftParenthesisToken
		case c == ')':
			token = qp.RightParenthesisToken
		case c == '{':
			token = qp.LeftBraceToken
		case c == '}':
			token = qp.RightBraceToken
		case c == '[':
			token = qp.LeftBracketToken
		case c == ']':
			token = qp.RightBracketToken
		case c == '<':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.Get()
				token = qp.LessEqualToken
			} else {
				token = qp.LessToken
			}
		case c == '>':
			if ahead, _ := l.ahead(); ahead == '=' {
				_, _ = l.Get()
				token = qp.GreaterEqualToken
			} else {
				token = qp.GreaterToken
			}
		case c == '*':
			token = qp.MulOperatorToken
		case '0' <= c && c <= '9':
			token = l.parseNumToken(c)
		case c == '!':
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.Get()
				token = qp.NoEqualToken
			} else {
				log.Panicf("unknown token`%s`", string(c))
			}
		case c == '-':
			token = qp.SubOperatorToken
		case c == '=':
			token = qp.AssignToken
			if c, _ = l.ahead(); c == '=' {
				_, _ = l.Get()
				token = qp.EqualToken
			}
		case c == ',':
			token = qp.CommaToken
		case c == ';':
			token = qp.SemicolonToken
		case c == ':':
			token = qp.ColonToken
		case c == '.':
			token = qp.PeriodToken
		case c == '"':
			token = l.parseString(false)
		case c == '`':
			token = l.parseString(true)
		case c == '/':
			if ahead, _ := l.ahead(); ahead == '/' { //
				_, _ = l.Get()
				token = qp.Token{
					Typ:  qp.CommentType,
					Val:  l.readline(),
					Line: l.line,
				}
			}
		case c == '|':
			if c, _ := l.ahead(); c == '|' {
				l.Get()
				token = qp.OrToken
			}
		case c == '&':
			if c, _ := l.ahead(); c == '&' {
				l.Get()
				token = qp.AndToken
			}
		default:
			log.Panicln(string(c), l.line)
		}
		token.Line = l.line
		l.token = token
		return token
	}
}

func (l *Lexer) parseString(multiline bool) qp.Token {
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
	return qp.Token{
		Typ:  qp.StringType,
		Val:  buffer.String(),
		Line: 0,
	}
}

func (l *Lexer) parseNumToken(c byte) qp.Token {
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.reader.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if IsDigit(c) {
			buf.WriteByte(c)
		} else {
			if err := l.reader.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	return qp.Token{
		Typ: qp.IntType,
		Val: buf.String(),
	}
}

func (l *Lexer) Next() {
	l.token = qp.EmptyToken
}

func (l *Lexer) parseLabel(c byte) qp.Token {
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.reader.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if IsLetter(c) || IsDigit(c) || c == '_' {
			buf.WriteByte(c)
		} else {
			if err := l.reader.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	for _, keyword := range qp.Keywords {
		if keyword == buf.String() {
			return qp.Token{
				Typ: qp.KeywordType[keyword],
			}
		}
	}

	return qp.Token{
		Typ: qp.IDType,
		Val: buf.String(),
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

func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func New(reader io.Reader) *Lexer {
	return &Lexer{
		token:  qp.EmptyToken,
		reader: bufio.NewReader(reader),
	}
}
