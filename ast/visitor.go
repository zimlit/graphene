package visitor

import "zimlit/graphene/ast"

type Visitor[R any] interface {
	visitBinary(b ast.Binary) R
	visitUnary(u ast.Unary) R
	visitLiteral(l ast.Literal) R
	visitGrouping(g ast.Grouping) R
}
