// chris 071615 Routines dealing with import management.

package main

import (
	"fmt"

	"go/ast"
	"go/token"
)

// prependImport prepends,in place, an import of the given name to the
// list of imports in the *ast.File.
func prependImport(astfile *ast.File, name string) {
	nopos := token.Pos(0)
	comment := fmt.Sprintf("// *** %s IMPORT ADDED BY %s ***", name, progUpper)
	litvalue := fmt.Sprintf("\"%s\"", name)

	decl := &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				&ast.Comment{
					Slash: nopos,
					Text:  comment,
				},
			},
		},
		TokPos: nopos,
		Tok:    token.IMPORT,
		Lparen: nopos,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Doc:  nil,
				Name: nil,
				Path: &ast.BasicLit{
					ValuePos: nopos,
					Kind:     token.STRING,
					Value:    litvalue,
				},
				Comment: nil,
				EndPos:  nopos,
			},
		},
		Rparen: nopos,
	}

	astfile.Decls = append([]ast.Decl{decl}, astfile.Decls...)
}
