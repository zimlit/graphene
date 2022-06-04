package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"zimlit/graphene/lexer"
	"zimlit/graphene/parser"

	"github.com/alecthomas/kong"
)

type Context struct{}

type BuildCmd struct {
	Output string `short:"o" help:"Specify ooutput file."`

	Path string `arg:"" name:"path" help:"Path to input file." type:"path"`
}

func (b *BuildCmd) Run(ctx *Context) error {
	dat, err := os.ReadFile(b.Path)
	if err != nil {
		return err
	}

	var out *os.File = os.Stdin
	if b.Output != "" {
		f, err := os.Create(b.Output)
		if err != nil {
			return err
		}
		out = f
	}

	l := lexer.NewLexer(string(dat), b.Path)
	toks, lines, errs := l.Lex()
	if errs != nil {
		fmt.Println(errs.Error())
		return nil
	}
	p := parser.NewParser(toks, lines, b.Path)
	c, len := p.Parse()

	for i := 0; i < len; i++ {
		res := <-c
		if res.Err != nil {
			fmt.Print(res.Err.Error())
		} else {
			for _, expr := range res.Exprs {
				out.WriteString(fmt.Sprintln(expr))
			}
		}
	}

	return nil
}

type ReplCmd struct {
}

func (r *ReplCmd) Run(xtx *Context) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		l := lexer.NewLexer(line, "stdin")
		toks, lines, errs := l.Lex()
		if errs != nil {
			fmt.Print(errs.Error())
		} else {
			fmt.Println(toks)
			p := parser.NewParser(toks, lines, "stdin")
			c, len := p.Parse()
			for i := 0; i < len; i++ {
				res := <-c
				if res.Err != nil {
					fmt.Print(res.Err.Error())
				} else {
					fmt.Println(res.Exprs)
				}
			}
		}
	}
}

var CLI struct {
	Build BuildCmd `cmd:"" help:"Build a file."`
	Repl  ReplCmd  `cmd:"" default:"1" help:"Open Repl."`
}

func main() {
	ctx := kong.Parse(&CLI)

	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
