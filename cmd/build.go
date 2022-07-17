/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
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
