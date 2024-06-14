package app

import "github.com/gosrob/autumn/internal/util/astutil"

type ProjectScanner struct {
	astutil.ProjectScanner
	scaned []scanedItem

	structs []StructDefinition
	funcs   []FuncDefinition
}

type scanedItem struct{}

func (s *ProjectScanner) ScanDirectory(dir string) error {
	s.ProjectScanner.ScanDirectory(dir)

	return nil
}

func (s *ProjectScanner) Parse(d astutil.Declaration) error {
	if d.Type == astutil.FuncDecl {
		bean := FuncDefinition{}
		n := d.ASTNode
		bean.Name = extractFuncName(n)
		bean.Params = extractParam(n)
		bean.Results = extractResult(n)
		s.funcs = append(s.funcs, bean)
	}
	if d.Type == astutil.StructDecl {
		bean := StructDefinition{}
		bean.Name = d.FullIdentity
		fds, err := extractFields(d.ASTNode, nil, "")
		if err != nil {
			return err
		}
		bean.Fields = fds
		s.structs = append(s.structs, bean)
	}
	return nil
}
