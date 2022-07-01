/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"zimlit/graphene/ast"
	"zimlit/graphene/lexer"
	"zimlit/graphene/parser"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "graphene",
	Short: "Graphene is a compiler for the graphene lanuage",
	Long:  `TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		rl, err := readline.New("> ")
		if err != nil {
			log.Fatal(err)
		}
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}
			l := lexer.NewLexer(line, "stdin")
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
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.graphene.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags()
}
