package main

import (
	"log"
	"os"

	"go/ast"
	"go/parser"
	"go/token"
)

func validator(name string, s *ast.StructType) {
	log.Print(name)
	for _, fld := range(s.Fields.List) {
		nam := fld.Names[0].Name
		typ := fld.Type.(*ast.Ident)
		log.Printf("%s %s", nam, typ)
	}
}

func main() {
	log.SetFlags(0)

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "-", os.Stdin, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range f.Scope.Objects {
		if obj.Kind != ast.Typ { continue }
		ts, ok := obj.Decl.(*ast.TypeSpec)
		if !ok { continue }
		s, ok := ts.Type.(*ast.StructType)
		if !ok { continue }
		if s.Fields == nil {
			log.Fatalf("type %s struct has empty field list %v", ts.Name, ts)
		}
		validator(ts.Name.Name, s)
	}
}
