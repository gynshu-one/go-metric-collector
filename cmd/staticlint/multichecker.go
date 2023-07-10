package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "noOsExit",
	Doc:  "reports the usage of os.Exit in main function",
	Run:  runNoOsExit,
}

// runNoOsExit reports the usage of os.Exit in main function
func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	// Loop through all files in the package
	for _, f := range pass.Files {
		for _, d := range f.Decls {
			fn, o := d.(*ast.FuncDecl)
			if !o {
				continue
			}

			// Check if the function is main
			if fn.Name.Name != "main" {
				continue
			}

			// Loop through all statements in the function
			for _, stmt := range fn.Body.List {

				// Check if the statement is an expression statement
				exprStmt, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}

				// Check if the expression is a call expression
				callExpr, ok := exprStmt.X.(*ast.CallExpr)
				if !ok {
					continue
				}

				// And finally if the function being called is os.Exit
				ident, ok := callExpr.Fun.(*ast.Ident)
				if !ok {
					continue
				}

				// If it is, report it
				if ident.Name != "os.Exit" {
					continue
				}
				pass.Reportf(callExpr.Pos(), "os.Exit should not be called in main function")
			}
		}
	}
	return nil, nil
}

func main() {
	analyzers := []*analysis.Analyzer{
		Analyzer,
		// Errcheck is a fork of the original errcheck tool with support for Go modules.
		errcheck.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		// All analyzers starting in SA "class"
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}

		// And one from other "class"
		if v.Analyzer.Name == "S1023" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// Run all analyzers
	multichecker.Main(analyzers...)
}
