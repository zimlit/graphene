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
			p := parser.NewParser(toks, lines)
			tree, errs := p.Parse()
			if errs != nil {
				fmt.Println(errs)
			} else {
				fmt.Println(tree)
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
