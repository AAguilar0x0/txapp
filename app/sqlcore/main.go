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
	"sync"
	"text/template"

	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"golang.org/x/mod/modfile"
)

const (
	SOURCE      = "./extern/db/psql/dal"
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
	Template        *template.Template
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
		file, err := os.Create(config.DestinationPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = config.Template.Execute(file, types)
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
		file, err := os.Create(config.DestinationPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = config.Template.Execute(file, interfaces)
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
					if name.Name == "newID" && paramType == "string" {
						continue
					}
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
		if returnType == "error" {
			returnType = "*apierrors.APIError"
		}
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
	default:
		return original
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
{{range .}}
type {{.Name}} struct {
  {{- range .Fields}}
  {{.Name}} {{.Type}} ` + "`{{.JSONTag}}`" + `
  {{- end}}
}
{{end}}`

const queryModelTemplate = `package models
{{range .}}
type {{.Name}} struct {
  {{- range .Fields}}
  {{.Name}} {{.Type}} ` + "`{{.JSONTag}}`" + `
  {{- end}}
}
{{end}}`

const interfaceTemplate = `package models

import (
	"context"
	"%s/core/pkg/apierrors"
)
{{range .}}
type {{.Name}} interface {
	{{- range .Methods}}
	{{.Name}}({{.Params}}) {{.Returns}}
	{{- end}}
}
{{end}}`

func getModuleName() string {
	goModBytes, err := os.ReadFile("go.mod")
	assert.NoError(err, "Failed to read go.mod file.")

	modFile, err := modfile.Parse("go.mod", goModBytes, nil)
	assert.NoError(err, "Failed to parse go.mod file.")

	if modFile.Module == nil || modFile.Module.Mod.Path == "" {
		assert.NoError(err, "No module declaration found in go.mod.")
	}

	return modFile.Module.Mod.Path
}

func main() {
	if err := os.MkdirAll(DESTINATION, fs.ModePerm); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	model := template.Must(template.New("model").Parse(baseModelTemplate))
	query := template.Must(template.New("query").Parse(queryModelTemplate))
	intrfc := template.Must(template.New("interface").Parse(fmt.Sprintf(interfaceTemplate, getModuleName())))

	wg.Add(1)
	go func() {
		defer wg.Done()
		generateStructs(GeneratorConfig{
			SourcePath:      filepath.Join(SOURCE, "models.go"),
			DestinationPath: filepath.Join(DESTINATION, "models.go"),
			Template:        model,
		})
	}()

	files, err := filepath.Glob(SOURCE + "/*.sql.go")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			baseName := filepath.Base(file)
			generateStructs(GeneratorConfig{
				SourcePath:      file,
				DestinationPath: filepath.Join(DESTINATION, baseName),
				Template:        query,
			})
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		generateInterfaces(GeneratorConfig{
			SourcePath:      filepath.Join(SOURCE, "querier.go"),
			DestinationPath: filepath.Join(DESTINATION, "querier.go"),
			Template:        intrfc,
		})
	}()

	wg.Wait()
}
