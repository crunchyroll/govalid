// chris 081015 Context for code generation.

package main

import (
	"bytes"
	"fmt"
)

type variableType struct {
	name     string
	typeExpr string
}

type generationContext struct {
	*bytes.Buffer

	// Set of import names.
	imports map[string]struct{}

	// Set of variables.
	variables map[string]string
}

func newContext() *generationContext {
	return &generationContext{
		Buffer:    new(bytes.Buffer),
		imports:   make(map[string]struct{}),
		variables: make(map[string]string),
	}
}

// write uses fmt.Sprintf on its arguments and writes the resultant
// string into the buffer.
func (ctx *generationContext) write(format string, a ...interface{}) {
	ctx.Buffer.WriteString(fmt.Sprintf(format, a...))
}

func (ctx *generationContext) addImport(name string) {
	ctx.imports[name] = struct{}{}
}

func (ctx *generationContext) addVariable(name, typeExpr string) {
	ctx.variables[name] = typeExpr
}

func (ctx *generationContext) getImports() []string {
	r := make([]string, 0, 10)
	for name, _ := range ctx.imports {
		r = append(r, name)
	}
	return r
}

func (ctx *generationContext) getVariables() []variableType {
	r := make([]variableType, 0, 10)
	for name, typeExpr := range ctx.variables {
		r = append(r, variableType{name, typeExpr})
	}
	return r
}
