package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"zimlit/graphene/token"

	"github.com/fatih/color"
)

type tmpLexErr struct {
	col   int
	line  int
	msg   string
	fname string
}

func (l *tmpLexErr) Error() string { return "" }

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
		fmt.Fprintln(&str, err.Error())
	}

	return str.String()
}

type Lexer struct {
	col      int
	pos      int
	line     int
	lineStr  string
	source   []rune
	fname    string
	keywords map[string]token.TokenKind
}

func NewLexer(source string, fname string) Lexer {
	l := Lexer{
		col:     1,
		pos:     0,
		line:    1,
		lineStr: "",
		source:  []rune(source),
		fname:   fname,
	}
	l.keywords = make(map[string]token.TokenKind)
	l.keywords["int"] = token.INTK
	l.keywords["float"] = token.FLOATK
	l.keywords["let"] = token.LET
	l.keywords["nil"] = token.NIL
	l.keywords["if"] = token.IF
	l.keywords["else"] = token.ELSE

	return l
}

func (l *Lexer) newTmpErr(msg string) tmpLexErr {
	return tmpLexErr{
		col:   l.col,
		line:  l.line,
		msg:   msg,
		fname: l.fname,
	}
}

func (l *Lexer) newTmpErrAt(msg string, col int) tmpLexErr {
	return tmpLexErr{
		col:   col,
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

func (l *Lexer) newTokenAt(literal string, kind token.TokenKind, col int) token.Token {
	return token.Token{
		Kind:    kind,
		Literal: literal,
		Line:    l.line,
		Col:     col,
	}
}

func (l *Lexer) peek() rune {
	return l.source[l.pos]
}

func (l *Lexer) peekNext() rune {
	return l.source[l.pos+1]
}

func (l *Lexer) advance() {
	l.lineStr += string(l.source[l.pos])
	l.pos++
	l.col++
}

func (l *Lexer) num() (*token.Token, *tmpLexErr) {
	val := ""
	dot_count := 0
	col := l.col
	for ; l.pos < len(l.source); l.advance() {
		val += string(l.peek())
		if l.peek() == '.' {
			dot_count++
		}

		if !(unicode.IsDigit(l.peekNext()) || l.peekNext() == '.') {
			break
		}
	}

	var tok token.Token

	if dot_count > 1 {
		e := l.newTmpErrAt("to many dots in number literal", col)
		return nil, &e
	} else if dot_count == 1 {
		tok = l.newTokenAt(val, token.FLOAT, col)
	} else {
		tok = l.newTokenAt(val, token.INT, col)
	}

	return &tok, nil
}

func (l *Lexer) ident() token.Token {
	val := ""
	col := l.col

	for ; l.pos < len(l.source); l.advance() {
		val += string(l.peek())

		if !(unicode.IsLetter(l.peekNext()) || unicode.IsDigit(l.peekNext()) || l.peekNext() == '_') {
			break
		}
	}
	t := l.keywords[val]
	if t == 0 {
		return l.newTokenAt(val, token.IDENT, col)
	}
	return l.newTokenAt(val, t, col)
}

func (l *Lexer) Lex() ([]token.Token, []string, LexErrs) {
	toks := []token.Token{}
	var errs LexErrs = nil
	tmps := []tmpLexErr{}
	lines := []string{}
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
		case '(':
			toks = append(toks, l.newToken("(", token.LPAREN))
		case ')':
			toks = append(toks, l.newToken(")", token.RPAREN))
		case '=':
			toks = append(toks, l.newToken("=", token.EQ))
		case ':':
			toks = append(toks, l.newToken(":", token.COLON))
		case ' ':
		case '\t':
		case '\n':
			for _, err := range tmps {
				errs = append(errs, l.newLexErr(err))
			}
			lines = append(lines, l.lineStr)
			l.line++
			l.lineStr = ""
		default:
			if unicode.IsDigit(l.peek()) {
				t, err := l.num()

				if err != nil {
					tmps = append(tmps, *err)
				} else {
					toks = append(toks, *t)
				}

			} else if unicode.IsLetter(l.peek()) || l.peek() == '_' {
				t := l.ident()
				toks = append(toks, t)
			} else {
				tmps = append(tmps, l.newTmpErr(fmt.Sprintf("Unexpected character '%s'", string(l.peek()))))

			}

		}
	}

	if errs != nil {
		return nil, nil, errs
	}
	return toks, lines, nil
}
