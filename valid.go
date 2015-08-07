// chris 071615 Validator code.

package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"go/ast"
	"unicode/utf8"
)

// write uses fmt.Sprintf on its arguments and writes the resultant
// string into the given buffer.
func write(buf *bytes.Buffer, format string, a ...interface{}) {
	buf.WriteString(fmt.Sprintf(format, a...))
}

// validateString writes validator code for a string to the given
// buffer.
func validateString(buf *bytes.Buffer, fldname string) {
	write(buf, "\tret.%s = data[\"%s\"]\n", fldname, fldname)
}

// validateBool writes validator code for a bool to the given buffer.
func validateBool(buf *bytes.Buffer, fldname string) {
	write(buf, "\tret.%s, err = strconv.ParseBool(data[\"%s\"])\n", fldname, fldname)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
}

// It would be nice if we didn't have as much duplication of generated
// code between the numeric validators.

// validateUint writes validator code for a uint of the given bitSize to
// the given buffer.
func validateUint(buf *bytes.Buffer, fldname string, bitSize int) {
	write(buf, "\tvar %stmp uint64\n", fldname)
	write(buf, "\t%stmp, err = strconv.ParseUint(data[\"%s\"], 0, %d)\n", fldname, fldname, bitSize)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
	// Have to cast since ParseUint returns a uint64.
	if bitSize == 0 {
		write(buf, "\tret.%s = uint(%stmp)\n", fldname, fldname)
	} else if bitSize != 64 {
		write(buf, "\tret.%s = uint%d(%stmp)\n", fldname, bitSize, fldname)
	} else {
		write(buf, "\tret.%s = %stmp\n", fldname, fldname)
	}
}

// validateInt writes validator code for an int of the given bitSize to
// the given buffer.
func validateInt(buf *bytes.Buffer, fldname string, bitSize int) {
	write(buf, "\tvar %stmp int64\n", fldname)
	write(buf, "\t%stmp, err = strconv.ParseInt(data[\"%s\"], 0, %d)\n", fldname, fldname, bitSize)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
	// Have to cast since ParseInt returns an int64.
	if bitSize == 0 {
		write(buf, "\tret.%s = int(%stmp)\n", fldname, fldname)
	} else if bitSize != 64 {
		write(buf, "\tret.%s = int%d(%stmp)\n", fldname, bitSize, fldname)
	} else {
		write(buf, "\tret.%s = %stmp\n", fldname, fldname)
	}
}

// validateFloat writes validator code for a float of the given bitSize to
// the given buffer.
func validateFloat(buf *bytes.Buffer, fldname string, bitSize int) {
	write(buf, "\tvar %stmp float64\n", fldname)
	write(buf, "\t%stmp, err = strconv.ParseFloat(data[\"%s\"], %d)\n", fldname, fldname, bitSize)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
	// Have to cast since ParseFloat returns a float64.  Superfluous
	// if bitSize is 64, but whatever.
	write(buf, "\tret.%s = float%d(%stmp)\n", fldname, bitSize, fldname)
}

// validateSimpleType delegates validator code generation given the name
// of the type.
func validateSimpleType(buf *bytes.Buffer, fieldname string, typename string) (needsStrconv bool) {
	switch typename {
	case "string":
		validateString(buf, fieldname)
		return false

	case "bool":
		validateBool(buf, fieldname)
		return true

	case "uint":
		validateUint(buf, fieldname, 0)
		return true
	case "uint8":
		validateUint(buf, fieldname, 8)
		return true
	case "uint16":
		validateUint(buf, fieldname, 16)
		return true
	case "uint32":
		validateUint(buf, fieldname, 32)
		return true
	case "uint64":
		validateUint(buf, fieldname, 64)
		return true

	case "int":
		validateInt(buf, fieldname, 0)
		return true
	case "int8":
		validateInt(buf, fieldname, 8)
		return true
	case "int16":
		validateInt(buf, fieldname, 16)
		return true
	case "int32":
		validateInt(buf, fieldname, 32)
		return true
	case "int64":
		validateInt(buf, fieldname, 64)
		return true

	case "float32":
		validateFloat(buf, fieldname, 32)
		return true
	case "float64":
		validateFloat(buf, fieldname, 64)
		return true
	}

	return false
}

// validateUrl writes validator code for a *mail.Address to the given
// buffer.
func validateMailAddress(buf *bytes.Buffer, fieldname string) {
	write(buf, "\tret.%s, err = mail.ParseAddress(data[\"%s\"])\n", fieldname, fieldname)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
}

// validateUrl writes validator code for a *url.URL to the given
// buffer.
func validateUrl(buf *bytes.Buffer, fieldname string) {
	write(buf, "\tret.%s, err = url.Parse(data[\"%s\"])\n", fieldname, fieldname)
	write(buf, "\tif err != nil {\n")
	write(buf, "\t\treturn nil, err\n")
	write(buf, "\t}\n")
}

// validator writes validator code for the given struct to the given
// buffer.  It iterates through the struct fields, and for those for
// which it can generate validator code, it does so.  It returns whether
// or not the strconv package is needed by the generated code.
func validator(buf *bytes.Buffer, structname string, s *ast.StructType) (needsStrconv bool) {
	first, _ := utf8.DecodeRune([]byte(structname))
	isPublic := unicode.IsUpper(first)
	var funcname string
	if isPublic {
		funcname = fmt.Sprintf("Validate%s", structname)
	} else {
		funcname = fmt.Sprintf("validate%s", strings.Title(structname))
	}

	write(buf, "\n") // Newline to separate from above content.

	write(buf, "// %s reads data from the given map of strings to\n", funcname)
	write(buf, "// strings and validates the data into a new *%s.\n", structname)
	write(buf, "// Fields named in a %s will be recognized as keys.\n", structname)
	write(buf, "// Keys in the input data that are not fields in the\n")
	write(buf, "// %s will be ignored.  If there is an error\n", structname)
	write(buf, "// validating any fields, an appropriate error will\n")
	write(buf, "// be returned.\n")

	write(buf, "func %s(data map[string]string) (*%s, error) {\n", funcname, structname)
	write(buf, "\tret := new(%s)\n", structname)

	// This declaration will cause a compile error if there are no
	// fields in the struct for which we can generate validators.
	write(buf, "\tvar err error\n")

	for _, field := range s.Fields.List {
		fieldname := field.Names[0].Name
		switch field.Type.(type) {

		// We'll look for a simple type.
		case *ast.Ident:
			ident := field.Type.(*ast.Ident)
			typename := ident.Name
			write(buf, "\t// %s %s\n", fieldname, typename)
			if validateSimpleType(buf, fieldname, typename) {
				needsStrconv = true
			}

		// We'll look for a pointer type.
		case *ast.StarExpr:
			star := field.Type.(*ast.StarExpr)
			sel, ok := star.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			pkg, ok2 := sel.X.(*ast.Ident)
			if !ok2 {
				continue
			}
			pkgname := pkg.Name
			typename := sel.Sel.Name
			write(buf, "\t// %s *%s.%s\n", fieldname, pkgname, typename)
			if pkgname == "mail" && typename == "Address" {
				validateMailAddress(buf, fieldname)
			} else if pkgname == "url" && typename == "URL" {
				validateUrl(buf, fieldname)
			}
		}
	}

	write(buf, "\n")
	write(buf, "\treturn ret, nil\n")
	write(buf, "}\n")

	return needsStrconv
}
