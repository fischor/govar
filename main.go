package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	if err := search(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func search(out, errOut io.Writer) error {
	if len(os.Args[1:]) != 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\tgovar PATH_TO_GO_FILE VAR_OR_CONST\n")
		os.Exit(1)
	}
	filename := os.Args[1]
	varname := os.Args[2]

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fs := token.NewFileSet()
	fs.AddFile(filename, fs.Base(), len(b))

	f, err := parser.ParseFile(fs, filename, b, parser.AllErrors)
	if err != nil {
		return err
	}

	v := visitor{out: out, errOut: errOut, varname: varname}
	ast.Walk(v, f)

	// ast walked and not exited already. That means not found.
	return fmt.Errorf("no such var or const: %s", varname)
}

type visitor struct {
	out     io.Writer
	errOut  io.Writer
	varname string
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.File:
		// Keep searching the ast.File node.
		return v
	case *ast.GenDecl:
		if d.Tok != token.CONST && d.Tok != token.VAR {
			return nil
		}
		for _, spec := range d.Specs {
			if vSpec, ok := spec.(*ast.ValueSpec); ok {
				for i, ident := range vSpec.Names {
					if ident.Name != v.varname {
						continue
					}
					// Found the identifier.
					lit, ok := vSpec.Values[i].(*ast.BasicLit)
					if !ok {
						fmt.Fprintf(v.errOut, "%s is declared but is not assigned to a literal\n", v.varname)
						os.Exit(1)
					}
					var val string
					switch lit.Kind {
					case token.CHAR:
						r, _, _, err := strconv.UnquoteChar(lit.Value, '\'')
						if err != nil {
							fmt.Fprintf(v.errOut, "error unquoting char: %v\n", err)
							os.Exit(1)
						}
						val = string(r)
					case token.STRING:
						var err error
						val, err = strconv.Unquote(lit.Value)
						if err != nil {
							fmt.Fprintf(v.errOut, "error unquoting string: %v\n", err)
							os.Exit(1)
						}
					default:
						val = lit.Value
					}
					fmt.Fprintf(v.out, "%s\n", val)
					os.Exit(0)
				}
			}
		}
		// Do not walk down the AssignStmt.
		return nil
	default:
		// Do not walk down anything other than the ast.File. We are only
		// looking for top-level vars and const.
		return nil
	}
}
