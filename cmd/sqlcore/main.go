package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	SOURCE      = "./extern/db/psql"
	DESTINATION = "./core/models"
)

type TypeInfo struct {
	Name   string
	Fields []FieldInfo
}

type FieldInfo struct {
	Name    string
	Type    string
	JSONTag string
}

type InterfaceInfo struct {
	Name    string
	Methods []MethodInfo
}

type MethodInfo struct {
	Name    string
	Params  string
	Returns string
}

type GeneratorConfig struct {
	SourcePath      string
	DestinationPath string
	Template        string
	PackageName     string
}

func generateStructs(config GeneratorConfig) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, config.SourcePath, nil, 0)
	if err != nil {
		panic(err)
	}

	var types []TypeInfo
	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				types = append(types, extractTypeInfo(typeSpec.Name.Name, structType))
			}
		}
		return true
	})

	if len(types) > 0 {
		tmpl := template.Must(template.New("model").Parse(config.Template))
		file, err := os.Create(config.DestinationPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = tmpl.Execute(file, struct {
			Package string
			Types   []TypeInfo
		}{
			Package: config.PackageName,
			Types:   types,
		})
		if err != nil {
			panic(err)
		}
	}
}

func generateInterfaces(config GeneratorConfig) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, config.SourcePath, nil, 0)
	if err != nil {
		panic(err)
	}

	var interfaces []InterfaceInfo
	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				interfaces = append(interfaces, extractInterfaceInfo(typeSpec.Name.Name, interfaceType))
			}
		}
		return true
	})

	if len(interfaces) > 0 {
		tmpl := template.Must(template.New("interface").Parse(config.Template))
		file, err := os.Create(config.DestinationPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = tmpl.Execute(file, struct {
			Package    string
			Interfaces []InterfaceInfo
		}{
			Package:    config.PackageName,
			Interfaces: interfaces,
		})
		if err != nil {
			panic(err)
		}
	}
}

func extractTypeInfo(name string, structType *ast.StructType) TypeInfo {
	var fields []FieldInfo
	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 {
			fieldName := field.Names[0].Name
			fieldType := getFieldType(field.Type)
			jsonTag := extractJSONTag(field.Tag)

			fields = append(fields, FieldInfo{
				Name:    fieldName,
				Type:    convertType(fieldType),
				JSONTag: jsonTag,
			})
		}
	}
	return TypeInfo{Name: name, Fields: fields}
}
func extractInterfaceInfo(name string, interfaceType *ast.InterfaceType) InterfaceInfo {
	var methods []MethodInfo
	for _, method := range interfaceType.Methods.List {
		if len(method.Names) > 0 {
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			params := extractFuncParams(funcType)
			returns := extractFuncReturns(funcType)

			methods = append(methods, MethodInfo{
				Name:    method.Names[0].Name,
				Params:  params,
				Returns: returns,
			})
		}
	}
	return InterfaceInfo{Name: name, Methods: methods}
}

func extractFuncParams(funcType *ast.FuncType) string {
	var params []string
	if funcType.Params != nil {
		for _, param := range funcType.Params.List {
			paramType := getFieldType(param.Type)
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
				}
			} else {
				params = append(params, paramType)
			}
		}
	}
	return strings.Join(params, ", ")
}

func extractFuncReturns(funcType *ast.FuncType) string {
	if funcType.Results == nil {
		return ""
	}

	var returns []string
	for _, result := range funcType.Results.List {
		returnType := getFieldType(result.Type)
		returns = append(returns, returnType)
	}

	if len(returns) == 0 {
		return ""
	} else if len(returns) == 1 {
		return returns[0]
	}
	return fmt.Sprintf("(%s)", strings.Join(returns, ", "))
}

func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func convertType(original string) string {
	switch original {
	case "time.Time":
		return "time.Time"
	default:
		return "string"
	}
}

func extractJSONTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	tagValue := strings.Trim(tag.Value, "`")
	return tagValue
}

const baseModelTemplate = `package models

import (
	"time"
)
{{range .Types}}
type {{.Name}} struct {
  {{- range .Fields}}
  {{.Name}} {{.Type}} ` + "`{{.JSONTag}}`" + `
  {{- end}}
}
{{end}}`

const queryModelTemplate = `package models
{{range .Types}}
type {{.Name}} struct {
  {{- range .Fields}}
  {{.Name}} {{.Type}} ` + "`{{.JSONTag}}`" + `
  {{- end}}
}
{{end}}`

const interfaceTemplate = `package models

import (
	"context"
)
{{range .Interfaces}}
type {{.Name}} interface {
	{{- range .Methods}}
	{{.Name}}({{.Params}}) {{.Returns}}
	{{- end}}
}
{{end}}`

func main() {
	if err := os.MkdirAll(DESTINATION, fs.ModePerm); err != nil {
		panic(err)
	}

	generateStructs(GeneratorConfig{
		SourcePath:      filepath.Join(SOURCE, "models.go"),
		DestinationPath: filepath.Join(DESTINATION, "models.go"),
		Template:        baseModelTemplate,
		PackageName:     "models",
	})

	files, err := filepath.Glob(SOURCE + "/*.sql.go")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		baseName := filepath.Base(file)
		packageName := strings.TrimSuffix(baseName, ".sql.go")

		generateStructs(GeneratorConfig{
			SourcePath:      file,
			DestinationPath: filepath.Join(DESTINATION, strings.TrimSuffix(baseName, ".go")+".go"),
			Template:        queryModelTemplate,
			PackageName:     packageName,
		})
	}

	generateInterfaces(GeneratorConfig{
		SourcePath:      filepath.Join(SOURCE, "querier.go"),
		DestinationPath: filepath.Join(DESTINATION, "interfaces.go"),
		Template:        interfaceTemplate,
		PackageName:     "models",
	})
}
