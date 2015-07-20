// chris 071615 Validator code.

package main

import (
	"fmt"
	"strings"
	"unicode"

	"go/ast"
	"unicode/utf8"
)

// validateString writes validator code for a string to the given *buf.
func validateString(b *buf, fldname string) {
	b.writef("\tret.%s = data[\"%s\"]\n", fldname, fldname)
}

// validateBool writes validator code for a bool to the given *buf.
func validateBool(b *buf, fldname string) {
	b.writef("\tret.%s, err = strconv.ParseBool(data[\"%s\"])\n", fldname, fldname)
	b.writef("\tif err != nil {\n")
	b.writef("\t\treturn nil, err\n")
	b.writef("\t}\n")
}

// validateUint writes validator code for a uint of the given bitSize to
// the given *buf.
func validateUint(b *buf, fldname string, bitSize int) {
	b.writef("\t%sTmp, err = strconv.ParseUint(data[\"%s\"], 0, %d)\n", fldname, fldname, bitSize)
	b.writef("\tif err != nil {\n")
	b.writef("\t\treturn nil, err\n")
	b.writef("\t}\n")
	// Have to cast since ParseUint returns a uint64.  Superfluous
	// if bitSize is 64, but whatever.
	b.writef("\tret.%s = uint%d(%sTmp)\n", fldname, bitSize, fldname)
}

// validateInt writes validator code for an int of the given bitSize to
// the given *buf.
func validateInt(b *buf, fldname string, bitSize int) {
	b.writef("\t%sTmp, err = strconv.ParseInt(data[\"%s\"], 0, %d)\n", fldname, fldname, bitSize)
	b.writef("\tif err != nil {\n")
	b.writef("\t\treturn nil, err\n")
	b.writef("\t}\n")
	// Have to cast since ParseUint returns an int64.  Superfluous
	// if bitSize is 64, but whatever.
	b.writef("\tret.%s = int%d(%sTmp)\n", fldname, bitSize, fldname)
}

// validator writes validator code for the given struct to the given
// *buf.  It iterates through the struct fields, and for those for which
// it can generate validator code, it does so.
func validator(b *buf, name string, s *ast.StructType) {
	first, _ := utf8.DecodeRune([]byte(name))
	isPublic := unicode.IsUpper(first)
	var fname string
	if isPublic {
		fname = fmt.Sprintf("Validate%s", name)
	} else {
		fname = fmt.Sprintf("validate%s", strings.Title(name))
	}

	b.writef("\n") // Newline to separate from above content.
	b.writef("func %s(data map[string]string) (*%s, error) {\n", fname, name)
	b.writef("\tret := new(%s)\n", name)

	for _, fld := range s.Fields.List {
		nam := fld.Names[0].Name
		typ, ok := fld.Type.(*ast.Ident)
		if !ok {
			continue
		}
		b.writef("\t// %s %s\n", nam, typ)
		switch typ.Name {
		case "string":
			validateString(b, nam)

		case "bool":
			validateBool(b, nam)
			b.needsStrconv = true

		case "uint8":
			validateUint(b, nam, 8)
		case "uint16":
			validateUint(b, nam, 16)
		case "uint32":
			validateUint(b, nam, 32)
		case "uint64":
			validateUint(b, nam, 64)

		case "int8":
			validateInt(b, nam, 8)
		case "int16":
			validateInt(b, nam, 16)
		case "int32":
			validateInt(b, nam, 32)
		case "int64":
			validateInt(b, nam, 64)
		}
	}

	b.writef("\t\n")
	b.writef("\treturn ret, nil\n")
	b.writef("}\n")
}
