// chris 071415

package main

import (
	"log"
	"os"
	"fmt"

	"go/ast"
	"go/parser"
	"go/token"
)

// Program name.  Set by init.
var prog string

func validator(name string, s *ast.StructType) {
	fmt.Println(name)
	for _, fld := range(s.Fields.List) {
		nam := fld.Names[0].Name
		typ := fld.Type.(*ast.Ident)
		fmt.Printf("%s %s\n", nam, typ)
	}
}

func main() {
	log.SetFlags(0)

	f, err := parser.ParseFile(token.NewFileSet(), "-", os.Stdin, 0)
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

func usage() {
	log.Printf("usage: %s file.v", path.Base(os.Args[0]))
	os.Exit(2)
}

func init() {
	log.SetFlags(0)
	prog = path.Base(os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	filename := os.Args[1]
	if filename == "" {
		usage()
	}

	err := parse(filename)
	if err != nil {
		log.Fatal(err)
	}
}
