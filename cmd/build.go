/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"zimlit/graphene/ast"
	"zimlit/graphene/lexer"
	"zimlit/graphene/parser"

	"github.com/spf13/cobra"
)

var output string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds the directory passed in [path] as a graphene project",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("can only take one path as argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		buf, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		source := string(buf)
		l := lexer.NewLexer(source, "stdin")
		toks, lines, errs := l.Lex()
		if errs != nil {
			fmt.Print(errs.Error())
		} else {
			p := parser.NewParser(toks, lines, "stdin")
			c := make(chan parser.ParseResult)
			go p.Parse(c)
			exprs := [][]ast.Expr{}
			parse_res := <-c
			if parse_res.Err != nil {
				fmt.Println(parse_res.Err.Error())
			} else {
				exprs = append(exprs, parse_res.Exprs)
			}
			for _, x := range exprs {
				for _, expr := range x {
					fmt.Println(expr.String())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
  buildCmd.Flags().StringVarP(&output, "output", "o", "", "file to write output to")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
