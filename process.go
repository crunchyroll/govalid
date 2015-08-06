// chris 080615

package main

import (
	"bytes"
	"fmt"
	"io"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// process parses the source, finds struct definitions, and invokes the
// validator code to produce validators for the struct definitions.  It
// prints the original code plus the generated code to the given
// destination writer.
func process(dst io.Writer, srcname string, src io.Reader) error {
	// Parse first before outputting anything.
	fset := token.NewFileSet()
	mode := parser.DeclarationErrors | parser.AllErrors
	astfile, err := parser.ParseFile(fset, srcname, src, mode)
	if err != nil {
		return err
	}

	// Buffer validator function code before outputting anything.
	// We do this because we need to know whether we need to augment
	// the import list before outputting any declarations (imports
	// must precede declarations).
	buf := new(bytes.Buffer)

	needsStrconv := false

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
		if validator(buf, ts.Name.Name, s) {
			needsStrconv = true
		}
	}

	// Add strconv import if needed.  Also, make more generic if
	// need be.  (E.g., adding other imports besides strconv, doing
	// non-linear search through existing imports, etc.)
	if needsStrconv && !hasImport(astfile, "strconv") {
		prependImport(astfile, "strconv")
	}

	// Output header comment.
	_, err = io.WriteString(dst, fmt.Sprintf("// *** GENERATED BY %s; DO NOT EDIT ***\n\n", args.progUpper))
	if err != nil {
		return err
	}

	// Next, output original code.
	err = printer.Fprint(dst, fset, astfile)
	if err != nil {
		return err
	}

	// Output generated code (from the buffer).
	io.Copy(dst, buf)

	return nil
}
