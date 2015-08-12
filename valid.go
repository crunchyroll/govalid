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

// It would be nice if we didn't have as much duplication of generated
// code between the validators.

// validateString writes validator code for a string.
func validateString(ctx *generationContext, fieldname string, meta *fieldMetadata) {
	ctx.addVariable(fmt.Sprintf("field_%s", fieldname), "string")
	ctx.addVariable("ok", "bool")

	ctx.write("\tfield_%s, ok = data[\"%s\"]\n", fieldname, fieldname)

	if meta.max != "" || meta.min != "" {
		ctx.write("\tif ok {\n")
		if meta.max != "" {
			ctx.addImport("errors")
			ctx.write("\t\tif len(field_%s) > %s {\n", fieldname, meta.max)
			ctx.write("\t\t\treturn nil, errors.New(\"%s can have a length of at most %s\")\n", fieldname, meta.max)
			ctx.write("\t\t}\n")
		}
		if meta.min != "" {
			ctx.addImport("errors")
			ctx.write("\t\tif len(field_%s) < %s {\n", fieldname, meta.min)
			ctx.write("\t\t\treturn nil, errors.New(\"%s must have a length of at least %s\")\n", fieldname, meta.min)
			ctx.write("\t\t}\n")
		}
		ctx.write("\t\tret.%s = field_%s\n", fieldname, fieldname)
		ctx.write("\t} else {\n")
	} else {
		ctx.write("\tif !ok {\n")
	}
	if meta.def != nil {
		ctx.write("\t\t// %s is optional.\n", fieldname)
		if *meta.def == "" {
			ctx.write("\t\t// Zero value already set.\n")
		} else {
			ctx.write("\t\tret.%s = %s\n", fieldname, *meta.def)
		}
	} else {
		ctx.addImport("errors")
		ctx.write("\t\treturn nil, errors.New(\"%s is required\")\n", fieldname)
	}
	ctx.write("\t}\n")
}

// validateBool writes validator code for a bool.
func validateBool(ctx *generationContext, fieldname string, meta *fieldMetadata) {
	ctx.addVariable(fmt.Sprintf("field_%s_s", fieldname), "string")
	ctx.addVariable("ok", "bool")
	ctx.addVariable("err", "error")

	ctx.write("\tfield_%s_s, ok = data[\"%s\"]\n", fieldname, fieldname)

	ctx.write("\tif ok {\n")
	ctx.write("\t\tret.%s, err = strconv.ParseBool(field_%s_s)\n", fieldname, fieldname)
	ctx.write("\t\tif err != nil {\n")
	ctx.write("\t\t\treturn nil, err\n")
	ctx.write("\t\t}\n")
	ctx.write("\t} else {\n")
	if meta.def != nil {
		ctx.write("\t\t// %s is optional.\n", fieldname)
		if *meta.def == "" {
			ctx.write("\t\t// Zero value already set.\n")
		} else {
			ctx.write("\t\tret.%s = %s\n", fieldname, *meta.def)
		}
	} else {
		ctx.addImport("errors")
		ctx.write("\t\treturn nil, errors.New(\"%s is required\")\n", fieldname)
	}
	ctx.write("\t}\n")
}

// validateUint writes validator code for a uint of the given bitSize.
func validateUint(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) {
	ctx.addVariable(fmt.Sprintf("field_%s_s", fieldname), "string")
	ctx.addVariable("ok", "bool")
	ctx.addVariable(fmt.Sprintf("field_%s", fieldname), "uint64")
	ctx.addVariable("err", "error")

	ctx.write("\tfield_%s_s, ok = data[\"%s\"]\n", fieldname, fieldname)
	ctx.write("\tif ok {\n")
	ctx.write("\t\tfield_%s, err = strconv.ParseUint(field_%s_s, 0, %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\t\tif err != nil {\n")
	ctx.write("\t\t\treturn nil, err\n")
	ctx.write("\t\t}\n")

	if meta.max != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s > %s{\n", fieldname, meta.max)
		ctx.write("\t\t\treturn nil, errors.New(\"%s can be at most %s\")\n", fieldname, meta.max)
		ctx.write("\t\t}\n")
	}
	if meta.min != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s < %s {\n", fieldname, meta.min)
		ctx.write("\t\t\treturn nil, errors.New(\"%s must be at least %s\")\n", fieldname, meta.min)
		ctx.write("\t\t}\n")
	}

	// Have to cast since ParseUint returns a uint64.
	if bitSize == 0 {
		ctx.write("\t\tret.%s = uint(field_%s)\n", fieldname, fieldname)
	} else if bitSize != 64 {
		ctx.write("\t\tret.%s = uint%d(field_%s)\n", fieldname, bitSize, fieldname)
	} else {
		ctx.write("\t\tret.%s = field_%s\n", fieldname, fieldname)
	}

	ctx.write("\t} else {\n")

	if meta.def != nil {
		ctx.write("\t\t// %s is optional.\n", fieldname)
		if *meta.def == "" {
			ctx.write("\t\t// Zero value already set.\n")
		} else {
			ctx.write("\t\tret.%s = %s\n", fieldname, *meta.def)
		}
	} else {
		ctx.addImport("errors")
		ctx.write("\t\treturn nil, errors.New(\"%s is required\")\n", fieldname)
	}

	ctx.write("\t}\n")
}

// validateInt writes validator code for an int of the given bitSize.
func validateInt(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) {
	ctx.addVariable(fmt.Sprintf("field_%s_s", fieldname), "string")
	ctx.addVariable("ok", "bool")
	ctx.addVariable(fmt.Sprintf("field_%s", fieldname), "int64")
	ctx.addVariable("err", "error")

	ctx.write("\tfield_%s_s, ok = data[\"%s\"]\n", fieldname, fieldname)
	ctx.write("\tif ok {\n")
	ctx.write("\t\tfield_%s, err = strconv.ParseInt(field_%s_s, 0, %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\t\tif err != nil {\n")
	ctx.write("\t\t\treturn nil, err\n")
	ctx.write("\t\t}\n")

	if meta.max != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s > %s{\n", fieldname, meta.max)
		ctx.write("\t\t\treturn nil, errors.New(\"%s can be at most %s\")\n", fieldname, meta.max)
		ctx.write("\t\t}\n")
	}
	if meta.min != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s < %s {\n", fieldname, meta.min)
		ctx.write("\t\t\treturn nil, errors.New(\"%s must be at least %s\")\n", fieldname, meta.min)
		ctx.write("\t\t}\n")
	}

	// Have to cast since ParseInt returns an int64.
	if bitSize == 0 {
		ctx.write("\t\tret.%s = int(field_%s)\n", fieldname, fieldname)
	} else if bitSize != 64 {
		ctx.write("\t\tret.%s = int%d(field_%s)\n", fieldname, bitSize, fieldname)
	} else {
		ctx.write("\t\tret.%s = field_%s\n", fieldname, fieldname)
	}

	ctx.write("\t} else {\n")

	if meta.def != nil {
		ctx.write("\t\t// %s is optional.\n", fieldname)
		if *meta.def == "" {
			ctx.write("\t\t// Zero value already set.\n")
		} else {
			ctx.write("\t\tret.%s = %s\n", fieldname, *meta.def)
		}
	} else {
		ctx.addImport("errors")
		ctx.write("\t\treturn nil, errors.New(\"%s is required\")\n", fieldname)
	}

	ctx.write("\t}\n")
}

// validateFloat writes validator code for a float of the given bitSize.
func validateFloat(ctx *generationContext, fieldname string, meta *fieldMetadata, bitSize int) {
	ctx.addVariable(fmt.Sprintf("field_%s_s", fieldname), "string")
	ctx.addVariable("ok", "bool")
	ctx.addVariable(fmt.Sprintf("field_%s", fieldname), "float64")
	ctx.addVariable("err", "error")

	ctx.write("\tfield_%s_s, ok = data[\"%s\"]\n", fieldname, fieldname)
	ctx.write("\tif ok {\n")
	ctx.write("\t\tfield_%s, err = strconv.ParseFloat(field_%s_s, 0, %d)\n", fieldname, fieldname, bitSize)
	ctx.write("\t\tif err != nil {\n")
	ctx.write("\t\t\treturn nil, err\n")
	ctx.write("\t\t}\n")

	if meta.max != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s > %s{\n", fieldname, meta.max)
		ctx.write("\t\t\treturn nil, errors.New(\"%s can be at most %s\")\n", fieldname, meta.max)
		ctx.write("\t\t}\n")
	}
	if meta.min != "" {
		ctx.addImport("errors")
		ctx.write("\t\tif field_%s < %s {\n", fieldname, meta.min)
		ctx.write("\t\t\treturn nil, errors.New(\"%s must be at least %s\")\n", fieldname, meta.min)
		ctx.write("\t\t}\n")
	}

	// Have to cast since ParseFloat returns a float64.
	if bitSize == 32 {
		ctx.write("\t\tret.%s = float32(field_%s)\n", fieldname, fieldname)
	} else { // 64
		ctx.write("\t\tret.%s = field_%s\n", fieldname, fieldname)
	}

	ctx.write("\t} else {\n")

	if meta.def != nil {
		ctx.write("\t\t// %s is optional.\n", fieldname)
		if *meta.def == "" {
			ctx.write("\t\t// Zero value already set.\n")
		} else {
			ctx.write("\t\tret.%s = %s\n", fieldname, *meta.def)
		}
	} else {
		ctx.addImport("errors")
		ctx.write("\t\treturn nil, errors.New(\"%s is required\")\n", fieldname)
	}

	ctx.write("\t}\n")
}

// validateSimpleType delegates validator code generation given the name
// of the type.
func validateSimpleType(ctx *generationContext, fieldname string, typename string, meta *fieldMetadata) {
	switch typename {
	case "string":
		validateString(ctx, fieldname, meta)

	case "bool":
		ctx.addImport("strconv")
		validateBool(ctx, fieldname, meta)

	case "uint":
		ctx.addImport("strconv")
		validateUint(ctx, fieldname, meta, 0)
	case "uint8":
		ctx.addImport("strconv")
		validateUint(ctx, fieldname, meta, 8)
	case "uint16":
		ctx.addImport("strconv")
		validateUint(ctx, fieldname, meta, 16)
	case "uint32":
		ctx.addImport("strconv")
		validateUint(ctx, fieldname, meta, 32)
	case "uint64":
		ctx.addImport("strconv")
		validateUint(ctx, fieldname, meta, 64)

	case "int":
		ctx.addImport("strconv")
		validateInt(ctx, fieldname, meta, 0)
	case "int8":
		ctx.addImport("strconv")
		validateInt(ctx, fieldname, meta, 8)
	case "int16":
		ctx.addImport("strconv")
		validateInt(ctx, fieldname, meta, 16)
	case "int32":
		ctx.addImport("strconv")
		validateInt(ctx, fieldname, meta, 32)
	case "int64":
		ctx.addImport("strconv")
		validateInt(ctx, fieldname, meta, 64)

	case "float32":
		ctx.addImport("strconv")
		validateFloat(ctx, fieldname, meta, 32)
	case "float64":
		ctx.addImport("strconv")
		validateFloat(ctx, fieldname, meta, 64)
	}
}

// validateUrl writes validator code for a *mail.Address.
func validateMailAddress(ctx *generationContext, fieldname string, meta *fieldMetadata) {
	ctx.addVariable("err", "error")

	if meta.max != "" {
		ctx.addImport("errors")
		ctx.write("\tif len(data[\"%s\"]) > %s {\n", fieldname, meta.max)
		ctx.write("\t\treturn nil, errors.New(\"%s can have a length of at most %s\")\n", fieldname, meta.max)
		ctx.write("\t}\n")
	}
	if meta.min != "" {
		ctx.addImport("errors")
		ctx.write("\tif len(data[\"%s\"]) < %s {\n", fieldname, meta.min)
		ctx.write("\t\treturn nil, errors.New(\"%s must have a length of at least %s\")\n", fieldname, meta.min)
		ctx.write("\t}\n")
	}

	ctx.write("\tret.%s, err = mail.ParseAddress(data[\"%s\"])\n", fieldname, fieldname)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")
}

// validateUrl writes validator code for a *url.URL.
func validateUrl(ctx *generationContext, fieldname string, meta *fieldMetadata) {
	ctx.addVariable("err", "error")

	if meta.max != "" {
		ctx.addImport("errors")
		ctx.write("\tif len(data[\"%s\"]) > %s {\n", fieldname, meta.max)
		ctx.write("\t\treturn nil, errors.New(\"%s can have a length of at most %s\")\n", fieldname, meta.max)
		ctx.write("\t}\n")
	}
	if meta.min != "" {
		ctx.addImport("errors")
		ctx.write("\tif len(data[\"%s\"]) < %s {\n", fieldname, meta.min)
		ctx.write("\t\treturn nil, errors.New(\"%s must have a length of at least %s\")\n", fieldname, meta.min)
		ctx.write("\t}\n")
	}

	ctx.write("\tret.%s, err = url.Parse(data[\"%s\"])\n", fieldname, fieldname)
	ctx.write("\tif err != nil {\n")
	ctx.write("\t\treturn nil, err\n")
	ctx.write("\t}\n")
}

func makeFunctionName(structname string) string {
	first, _ := utf8.DecodeRune([]byte(structname))
	isPublic := unicode.IsUpper(first)
	if isPublic {
		return fmt.Sprintf("Validate%s", structname)
	}
	return fmt.Sprintf("validate%s", strings.Title(structname))
}

func validatorImpl(ctx *generationContext, structtype *ast.StructType) {
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
			ctx.write("\n\t// %s %s\n", fieldname, typename)
			validateSimpleType(ctx, fieldname, typename, meta)

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
			ctx.write("\n\t// %s *%s.%s\n", fieldname, pkgname, typename)
			if pkgname == "mail" && typename == "Address" {
				validateMailAddress(ctx, fieldname, meta)
			} else if pkgname == "url" && typename == "URL" {
				validateUrl(ctx, fieldname, meta)
			}
		}
	}
}

func declareVariables(ctx *generationContext, vars []variableType) {
	if len(vars) == 0 {
		return
	}

	// Compute maximum length of variable names so we can do
	// gofmt-compatible alignment.
	max := 0
	for _, x := range vars {
		if len(x.name) > max {
			max = len(x.name)
		}
	}

	ctx.write("\tvar (\n")
	for _, x := range vars {
		nspaces := max - len(x.name) + 1
		spaces := strings.Repeat(" ", nspaces)
		ctx.write("\t\t%s%s%s\n", x.name, spaces, x.typeExpr)
	}
	ctx.write("\t)\n")
}

// validator writes validator code for the given struct.  It iterates
// through the struct fields, and for those for which it can generate
// validator code, it does so.  It returns whether or not the strconv
// package is needed by the generated code.
func validator(ctx *generationContext, structname string, structtype *ast.StructType) {
	// First, buffer inner contents of the function into a secondary
	// context.  This is so we can know what variables we'll need to
	// declare at the top of the function.
	ctx2 := newContext()

	// Generate the inner implementation of the validator function.
	validatorImpl(ctx2, structtype)

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
	declareVariables(ctx, ctx2.getVariables())

	// Copy over the inner implementation body itself.  Because
	// we're reading from a buffer, there's no actual error to
	// handle here.
	ctx.Buffer.ReadFrom(ctx2.Buffer)

	ctx.write("\n")
	ctx.write("\treturn ret, nil\n")
	ctx.write("}\n")

	// Migrate needed imports to parent context so caller can get
	// ahold of them.
	for _, importName := range ctx2.getImports() {
		ctx.addImport(importName)
	}
}
