// chris 071615 Routines dealing with import management.

package main

import (
	"fmt"

	"go/ast"
	"go/token"
)

// findImports returns the imports mentioned in the given *ast.File as
// simple strings.  Basically, it pulls out the bit in quotes.
func findImports(astfile *ast.File) []string {
	r := make([]string, 0, len(astfile.Imports))
	for _, is := range astfile.Imports {
		if is.Path.Value == "" {
			continue
		}
		// Cut off leading and trailing double quote characters.
		v := is.Path.Value[1 : len(is.Path.Value)-1]
		r = append(r, v)
	}
	return r
}

// hasImport determines if the given *ast.File already has an import of
// the given name.
func hasImport(astfile *ast.File, name string) bool {
	for _, existingImport := range findImports(astfile) {
		if existingImport == name {
			return true
		}
	}
	return false
}

// prependImport prepends,in place, an import of the given name to the
// list of imports in the *ast.File.
func prependImport(astfile *ast.File, name string) {
	litvalue := fmt.Sprintf("\"%s\"", name)
	commentText := fmt.Sprintf("// *** %s IMPORT ADDED BY %s ***", name, args.progUpper)
	comment := &ast.CommentGroup{
		List: []*ast.Comment{
			&ast.Comment{
				Slash: token.NoPos,
				Text:  commentText,
			},
		},
	}

	decl := &ast.GenDecl{
		Doc:    comment,
		TokPos: token.NoPos,
		Tok:    token.IMPORT,
		Lparen: token.NoPos,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Doc:  nil,
				Name: nil,
				Path: &ast.BasicLit{
					ValuePos: token.NoPos,
					Kind:     token.STRING,
					Value:    litvalue,
				},
				Comment: nil,
				EndPos:  token.NoPos,
			},
		},
		Rparen: token.NoPos,
	}

	astfile.Decls = append([]ast.Decl{decl}, astfile.Decls...)
}
