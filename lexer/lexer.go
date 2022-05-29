package lexer

import (
	"fmt"
	"strings"
	"zimlit/graphene/token"

	"github.com/fatih/color"
)

type tmpLexErr struct {
	col   int
	line  int
	msg   string
	fname string
}

type LexErr struct {
	col     int
	line    int
	lineStr string
	msg     string
	fname   string
}

func (l *LexErr) Error() string {
	r := color.New(color.FgHiRed, color.Bold).FprintfFunc()
	w := color.New(color.FgHiWhite, color.Bold).FprintfFunc()
	b := color.New(color.FgHiBlue, color.Bold).FprintfFunc()
	var str strings.Builder
	r(&str, "error")
	fmt.Fprint(&str, ": ")
	w(&str, "%s\n", l.msg)
	b(&str, " --> ")
	fmt.Fprintf(&str, "%s:%d:%d\n", l.fname, l.line, l.col)
	b(&str, "  |\n")
	b(&str, "%d | ", l.line)
	fmt.Fprintln(&str, l.lineStr)
	b(&str, "  |")
	for i := 0; i < l.col; i++ {
		fmt.Fprint(&str, " ")
	}
	r(&str, "^ %s\n", l.msg)

	return str.String()
}

type LexErrs []LexErr

func (l *LexErrs) Error() string {
	var str strings.Builder
	for _, err := range *l {
		fmt.Fprint(&str, err.Error())
		fmt.Fprintln(&str, "")
	}

	return str.String()
}

type Lexer struct {
	col     int
	pos     int
	line    int
	lineStr string
	source  []rune
	fname   string
}

func NewLexer(source string, fname string) Lexer {
	return Lexer{
		col:     1,
		pos:     0,
		line:    1,
		lineStr: "",
		source:  []rune(source),
		fname:   fname,
	}
}

func (l *Lexer) newTmpErr(msg string) tmpLexErr {
	return tmpLexErr{
		col:   l.col,
		line:  l.line,
		msg:   msg,
		fname: l.fname,
	}
}

func (l *Lexer) newLexErr(tmp tmpLexErr) LexErr {
	return LexErr{
		col:     tmp.col,
		line:    tmp.line,
		msg:     tmp.msg,
		lineStr: l.lineStr,
		fname:   tmp.fname,
	}
}

func (l *Lexer) newToken(literal string, kind token.TokenKind) token.Token {
	return token.Token{
		Kind:    kind,
		Literal: literal,
		Line:    l.line,
		Col:     l.col,
	}
}

func (l *Lexer) peek() rune {
	return l.source[l.pos]
}

func (l *Lexer) advance() {
	l.lineStr += string(l.source[l.pos])
	l.pos++
	l.col++
}

func (l *Lexer) Lex() ([]token.Token, LexErrs) {
	toks := []token.Token{}
	var errs LexErrs = nil
	tmps := []tmpLexErr{}
	for ; l.pos < len(l.source); l.advance() {
		switch l.peek() {
		case '+':
			toks = append(toks, l.newToken("+", token.PLUS))
		case '-':
			toks = append(toks, l.newToken("-", token.MINUS))
		case '*':
			toks = append(toks, l.newToken("*", token.STAR))
		case '/':
			toks = append(toks, l.newToken("/", token.SLASH))
		case ' ':
		case '\t':
		case '\n':
			for _, err := range tmps {
				errs = append(errs, l.newLexErr(err))
			}
			l.line++
		default:
			tmps = append(tmps, l.newTmpErr(fmt.Sprintf("Unexpected character '%s'", string(l.peek()))))
		}
	}

	if errs != nil {
		return nil, errs
	}
	return toks, nil
}
