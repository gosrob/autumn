package astutil

import "go/ast"

type Visitor struct {
	visitor func(ast.Node)
}

func NewVisitor(visitor func(ast.Node)) Visitor {
	return Visitor{
		visitor: visitor,
	}
}

// Visit implements ast.Visitor.
func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	v.visitor(node)
	return v
}

var _ ast.Visitor = (*Visitor)(nil)
