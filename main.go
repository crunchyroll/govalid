// chris 071415

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// Program name variables.  Set by init.
var prog string
var progUpper string

// process parses the input, finds struct definitions, and invokes the
// validator code to produce validators for the struct definitions.  It
// prints the original code plus the generated code to standard out.
func process(filename string, input io.Reader) error {
	dst := os.Stdout // Destination.

	// Parse first before outputting anything.
	fset := token.NewFileSet()
	mode := parser.DeclarationErrors | parser.AllErrors
	astfile, err := parser.ParseFile(fset, filename, input, mode)
	if err != nil {
		return err
	}

	// Buffer validator function code before outputting anything.
	// We do this because we need to know whether we need to augment
	// the import list before outputting any declarations (imports
	// must precede declarations).
	b := newBuf()

	// Isolate the struct types--the things for which we want to
	// generate validator functions.
	for _, obj := range astfile.Scope.Objects {
		if obj.Kind != ast.Typ {
			continue
		}
		ts, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}
		if s.Fields == nil {
			return fmt.Errorf("type %s struct has empty field list %v", ts.Name, ts)
		}

		// Ok, we isolated the struct type, now output a
		// validator for it.
		validator(b, ts.Name.Name, s)
	}

	// Add strconv import if needed.  Also, make more generic if
	// need be.  (E.g., adding other imports besides strconv, doing
	// non-linear search through existing imports, etc.)
	if b.needsStrconv && !hasImport(astfile, "strconv") {
		prependImport(astfile, "strconv")
	}

	// Output header comment.
	_, err = fmt.Printf("// *** GENERATED BY %s; DO NOT EDIT ***\n\n", progUpper)
	if err != nil {
		return err
	}

	// Next, output original code.
	err = printer.Fprint(dst, fset, astfile)
	if err != nil {
		return err
	}

	// Output generated code (from the buf).  Includes a leading
	// newline to separate from the original code.
	io.Copy(dst, b)

	return nil
}

func usage() {
	log.Printf("usage: %s [file.v]", path.Base(os.Args[0]))
	os.Exit(2)
}

func init() {
	log.SetFlags(0)
	prog = path.Base(os.Args[0])
	progUpper = strings.ToUpper(prog)
}

func main() {
	var filename string
	var input    io.Reader
	if len(os.Args) == 1 {
		filename = "-"
		input = os.Stdin
	} else if len(os.Args) == 2 {
		filename = os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		input = file
	} else {
		usage()
	}
	err := process(filename, input)
	if err != nil {
		log.Fatal(err)
	}
}
