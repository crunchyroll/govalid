// chris 071615 Generation of validator code.

package main

import (
	"fmt"
	"strings"
	"unicode"

	"go/ast"
	"go/token"
	"unicode/utf8"
)

// validateString writes validator code for a string.
func validateString(ctx *generationContext, fieldname string, meta *fieldMetadata) error {
	ctx.write("\tret.%s = data[\"%s\"]\n", fieldname, fieldname)

	return nil
}

// validateBool writes validator code for a bool.
func validateBool(ctx *generationContext, fieldname string, meta *fieldMetadata) error {
	ctx.addVariable("err", "error")

	ctx.write("\tret.%s, err = strconv.ParseBool(data[\"%s\"])\n", fieldname, fieldname)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")

	return nil
}

// It would be nice if we didn't have as much duplication of generated
// code between the numeric validators.

// validateUint writes validator code for a uint of the given bitSize.
func validateUint(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) error {
	ctx.addVariable("err", "error")

	ctx.write("\tvar %stmp uint64\n", fieldname)
	ctx.write("\t%stmp, err = strconv.ParseUint(data[\"%s\"], 0, %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")
	// Have to cast since ParseUint returns a uint64.
	if bitSize == 0 {
		ctx.write("\tret.%s = uint(%stmp)\n", fieldname, fieldname)
	} else if bitSize != 64 {
		ctx.write("\tret.%s = uint%d(%stmp)\n", fieldname, bitSize, fieldname)
	} else {
		ctx.write("\tret.%s = %stmp\n", fieldname, fieldname)
	}

	return nil
}

// validateInt writes validator code for an int of the given bitSize.
func validateInt(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) error {
	ctx.addVariable("err", "error")

	ctx.write("\tvar %stmp int64\n", fieldname)
	ctx.write("\t%stmp, err = strconv.ParseInt(data[\"%s\"], 0, %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")
	// Have to cast since ParseInt returns an int64.
	if bitSize == 0 {
		ctx.write("\tret.%s = int(%stmp)\n", fieldname, fieldname)
	} else if bitSize != 64 {
		ctx.write("\tret.%s = int%d(%stmp)\n", fieldname, bitSize, fieldname)
	} else {
		ctx.write("\tret.%s = %stmp\n", fieldname, fieldname)
	}

	return nil
}

// validateFloat writes validator code for a float of the given bitSize.
func validateFloat(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) error {
	ctx.addVariable("err", "error")

	ctx.write("\tvar %stmp float64\n", fieldname)
	ctx.write("\t%stmp, err = strconv.ParseFloat(data[\"%s\"], %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")
	// Have to cast since ParseFloat returns a float64.  Superfluous
	// if bitSize is 64, but whatever.
	ctx.write("\tret.%s = float%d(%stmp)\n", fieldname, bitSize, fieldname)

	return nil
}

// validateSimpleType delegates validator code generation given the name
// of the type.
func validateSimpleType(ctx *generationContext, fieldname string, typename string, meta *fieldMetadata) error {
	switch typename {
	case "string":
		if err := checkFieldMetadata(meta); err != nil {
			return err
		}
		return validateString(ctx, fieldname, meta)

	case "bool":
		ctx.addImport("strconv")
		return validateBool(ctx, fieldname, meta)

	case "uint":
		ctx.addImport("strconv")
		return validateUint(ctx, fieldname, meta, 0)
	case "uint8":
		ctx.addImport("strconv")
		return validateUint(ctx, fieldname, meta, 8)
	case "uint16":
		ctx.addImport("strconv")
		return validateUint(ctx, fieldname, meta, 16)
	case "uint32":
		ctx.addImport("strconv")
		return validateUint(ctx, fieldname, meta, 32)
	case "uint64":
		ctx.addImport("strconv")
		return validateUint(ctx, fieldname, meta, 64)

	case "int":
		ctx.addImport("strconv")
		return validateInt(ctx, fieldname, meta, 0)
	case "int8":
		ctx.addImport("strconv")
		return validateInt(ctx, fieldname, meta, 8)
	case "int16":
		ctx.addImport("strconv")
		return validateInt(ctx, fieldname, meta, 16)
	case "int32":
		ctx.addImport("strconv")
		return validateInt(ctx, fieldname, meta, 32)
	case "int64":
		ctx.addImport("strconv")
		return validateInt(ctx, fieldname, meta, 64)

	case "float32":
		ctx.addImport("strconv")
		return validateFloat(ctx, fieldname, meta, 32)
	case "float64":
		ctx.addImport("strconv")
		return validateFloat(ctx, fieldname, meta, 64)
	}

	return nil
}

// validateUrl writes validator code for a *mail.Address.
func validateMailAddress(ctx *generationContext, fieldname string, meta *fieldMetadata) error {
	ctx.addVariable("err", "error")

	ctx.write("\tret.%s, err = mail.ParseAddress(data[\"%s\"])\n", fieldname, fieldname)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")

	return nil
}

// validateUrl writes validator code for a *url.URL.
func validateUrl(ctx *generationContext, fieldname string, meta *fieldMetadata) error {
	ctx.addVariable("err", "error")

	ctx.write("\tret.%s, err = url.Parse(data[\"%s\"])\n", fieldname, fieldname)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")

	return nil
}

func makeFunctionName(structname string) string {
	first, _ := utf8.DecodeRune([]byte(structname))
	isPublic := unicode.IsUpper(first)
	if isPublic {
		return fmt.Sprintf("Validate%s", structname)
	}
	return fmt.Sprintf("validate%s", strings.Title(structname))
}

func validatorImpl(ctx *generationContext, structtype *ast.StructType) error {
	for _, field := range structtype.Fields.List {
		fieldname := field.Names[0].Name

		var tagstring string
		if field.Tag != nil && field.Tag.Kind == token.STRING {
			tagstring = field.Tag.Value
		} else {
			tagstring = ""
		}
		meta := parseFieldMetadata(tagstring)

		switch field.Type.(type) {

		// We'll look for a simple type.
		case *ast.Ident:
			ident := field.Type.(*ast.Ident)
			typename := ident.Name
			ctx.write("\t// %s %s\n", fieldname, typename)
			if err := validateSimpleType(ctx, fieldname, typename, meta); err != nil {
				return err
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
			ctx.write("\t// %s *%s.%s\n", fieldname, pkgname, typename)
			if pkgname == "mail" && typename == "Address" {
				if err := checkFieldMetadata(meta); err != nil {
					return err
				}
				if err := validateMailAddress(ctx, fieldname, meta); err != nil {
					return err
				}
			} else if pkgname == "url" && typename == "URL" {
				if err := checkFieldMetadata(meta); err != nil {
					return err
				}
				if err := validateUrl(ctx, fieldname, meta); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validator writes validator code for the given struct.  It iterates
// through the struct fields, and for those for which it can generate
// validator code, it does so.  It returns whether or not the strconv
// package is needed by the generated code.
func validator(ctx *generationContext, structname string, structtype *ast.StructType) error {
	// First, buffer inner contents of the function into a secondary
	// context.  This is so we can know what variables we'll need to
	// declare at the top of the function.
	ctx2 := newContext()

	// Generate the inner implementation of the validator function.
	if err := validatorImpl(ctx2, structtype); err != nil {
		return err
	}

	// Now that that's succeeded, we can actually output all of the
	// code.
	funcname := makeFunctionName(structname)

	ctx.write("\n") // Newline to separate from prior content.

	// Add descriptive comment, GoDoc/golint compatible.
	ctx.write("// %s reads data from the given map of strings to\n", funcname)
	ctx.write("// strings and validates the data into a new *%s.\n", structname)
	ctx.write("// Fields named in a %s will be recognized as keys.\n", structname)
	ctx.write("// Keys in the input data that are not fields in the\n")
	ctx.write("// %s will be ignored.  If there is an error\n", structname)
	ctx.write("// validating any fields, an appropriate error will\n")
	ctx.write("// be returned.\n")

	ctx.write("func %s(data map[string]string) (*%s, error) {\n", funcname, structname)
	ctx.write("\tret := new(%s)\n", structname)

	// Delcare variables needed by the implementation.
	for _, x := range ctx2.getVariables() {
		ctx.write("\tvar %s %s\n", x.name, x.typeExpr)
	}

	// Copy over the inner implementation body itself.
	if _, err := ctx.Buffer.ReadFrom(ctx2.Buffer); err != nil {
		return err
	}

	ctx.write("\n")
	ctx.write("\treturn ret, nil\n")
	ctx.write("}\n")

	// Migrate needed imports to parent context so caller can get
	// ahold of them.
	for _, importName := range ctx2.getImports() {
		ctx.addImport(importName)
	}

	return nil
}
