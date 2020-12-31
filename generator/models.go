package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	funcs = template.FuncMap{
		"lowerCamel":          strcase.ToLowerCamel,
		"toCamel":             strcase.ToCamel,
		"toScreamingSnake":    strcase.ToScreamingSnake,
		"toSnake":             strcase.ToSnake,
		"fieldWithDecoration": getFieldDecoration,
	}
	// guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
)

// ModelData holds the data for the model definition
type ModelData struct {
	Name string
	// Date       string
	// Short      string
	// Link       string
	// TimeToRead string
	ClassDeclaration, ImportStatements, StaticFields, Fields, Transformer string
	Dependencies                                                          []string
	Attributes                                                            []*Attribute
}

// ModelGenerator Object
type ModelGenerator struct {
	Config *ModelConfig
}

// ModelConfig holds the configuration for the model definition
type ModelConfig struct {
	Entity                 *Entity
	Template               *template.Template
	Destination, ModelName string
	// IsIndex                bool
	// Writer *ModelWriter
}

// ModelWriter writes model files
type ModelWriter struct {
	ModelName string
	// BlogDescription string
	// BlogAuthor      string
	// BlogURL         string
}

func getFieldDecoration(attribute Attribute) string {
	var fieldDeclaration string
	fieldName := strcase.ToLowerCamel(attribute.Name)
	switch attribute.Meta.AttributeType {
	case "string":
		fieldDeclaration = fmt.Sprintf("String %s", fieldName)
	}
	if attribute.Meta.Required == true {

	}
	return fieldDeclaration
}

func getRelativePathDiff(source, target string) string {
	// targpath := "/lib/features/home/data/models"
	// basepath := "/lib/core/types"

	relpath, _ := filepath.Rel(source, target)
	// relpath = filepath.ToSlash(relpath)
	return relpath
}
func getDependencies(config *ModelConfig) ([]string, error) {
	result := []string{}
	entity := config.Entity
	for _, attribute := range entity.Attributes {
		if attribute.Meta.AttributeType == "multi" || attribute.Meta.AttributeType == "single" {
			// fmt.Println("path is ", getRelativePathDiff(entity.Meta.Path, attribute.Meta.Association.Entity.Meta.Path))
			relativePath := getRelativePathDiff(entity.Meta.Path, attribute.Meta.Association.Entity.Meta.Path)
			fmt.Println("Relative Path:", relativePath)
			fileName := fmt.Sprintf("%s%s", strcase.ToSnake(attribute.Meta.Association.Entity.Name), ".dart")
			result = append(result, filepath.ToSlash(filepath.Join(relativePath, fileName)))
		}
	}
	return result, nil
}

// WriteModel writes a Model.dart file
func (mw *ModelWriter) WriteModel(path, classDeclaration, importStatements, staticFields, fields, transformer string, w io.Writer, config *ModelConfig) error {
	filePath := filepath.Join(path, fmt.Sprintf("%s%s", mw.ModelName, "_model.dart"))
	// config.Template.Funcs(funcs)
	// f, err := os.Create(filePath)
	// if err != nil {
	// 	return fmt.Errorf("error creating file %s: %v", filePath, err)
	// }
	// defer f.Close()
	bw := bufio.NewWriter(w)
	dependencies, err := getDependencies(config)
	fmt.Println("dependencies ", dependencies)
	if err != nil {
		return fmt.Errorf("error retrieving dependencies: %v", err)
	}
	md := ModelData{
		ClassDeclaration: classDeclaration,
		ImportStatements: importStatements,
		StaticFields:     staticFields,
		Fields:           fields,
		Transformer:      transformer,
		Name:             mw.ModelName,
		Dependencies:     dependencies,
		Attributes:       config.Entity.Attributes,
	}

	if err := config.Template.Execute(bw, md); err != nil {
		return fmt.Errorf("error executing template %s: %v", filePath, err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

// Generate generates a model
func (g *ModelGenerator) Generate() error {
	const (
		// master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		classDeclaration = `class {{toCamel .Name}}Model {`
		fields           = `{{range .}} {{"\n"}}{{"\t"}}String {{"_"}}{{lowerCamel .Name}}{{";"}}{{end}} `
		staticFields     = `{{range .}} {{"\n"}}{{"\t"}}static const {{toScreamingSnake .Name}}{{" = \""}}{{lowerCamel .Name}}{{"\";"}}{{end}} `
		staticImport     = `import {{.}}`
	)
	var (
		funcs = template.FuncMap{
			"lowerCamel":       strcase.ToLowerCamel,
			"toCamel":          strcase.ToCamel,
			"toScreamingSnake": strcase.ToScreamingSnake,
		}
		// guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	entity := g.Config.Entity
	destination := g.Config.Destination
	// t := g.Config.Template
	// t.Funcs(funcs)
	fmt.Printf("\tGenerating Model: %v%s...\n", entity.Name, "_model")
	// staticPath := filepath.Join(destination, entity.Name)
	staticPath := filepath.Join(destination, "lib", "features", "home", "models")
	if err := os.MkdirAll(staticPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory at %s: %v", staticPath, err)
	}
	importBlock := bytes.Buffer{}
	WriteModelImport(&importBlock, g.Config)
	classDeclTemplate, err := template.New("classDeclaration").Funcs(funcs).Parse(classDeclaration)
	if err != nil {
		return fmt.Errorf("error creating template %s: %v", "classDeclaration", err)
	}
	fieldsTemplate, err := template.New("fields").Funcs(funcs).Parse(fields)
	if err != nil {
		return fmt.Errorf("error creating template %s: %v", "fields", err)
	}
	staticFieldsTemplate, err := template.New("static_fields").Funcs(funcs).Parse(staticFields)
	if err != nil {
		return fmt.Errorf("error creating template %s: %v", "static_fields", err)
	}
	classDeclBlock := bytes.Buffer{}
	err = classDeclTemplate.Execute(&classDeclBlock, entity)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", "classDeclaration", err)
	}
	fieldsBlock := bytes.Buffer{}
	err = fieldsTemplate.Execute(&fieldsBlock, entity.Attributes)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", "fields", err)
	}
	staticFieldsBlock := bytes.Buffer{}
	err = staticFieldsTemplate.Execute(&staticFieldsBlock, entity.Attributes)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", "static_fields", err)
	}
	// md := ModelData{
	// 	ClassDeclaration: classDeclBlock.String(),
	// 	ImportStatements: importBlock.String(),
	// 	StaticFields:     staticFieldsBlock.String(),
	// 	Fields:           fieldsBlock.String(),
	// 	Transformer:      "transformer",
	// 	Name:             "User Model",
	// }
	filePath := filepath.Join(staticPath, fmt.Sprintf("%s%s", entity.Name, "_model.dart"))
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	modelWriter := &ModelWriter{
		ModelName: entity.Name,
	}
	if err := modelWriter.WriteModel(staticPath, classDeclBlock.String(), importBlock.String(), staticFieldsBlock.String(), fieldsBlock.String(), "transformers", f, g.Config); err != nil {
		return err
	}
	fmt.Printf("\tFinished generating Model: %s...\n", entity.Name)
	return nil
}

//Filter accepts a function and returns a filtered slice of Attributes
func Filter(arr []*Attribute, cond func(*Attribute) bool) []*Attribute {
	result := []*Attribute{}
	// fmt.Println("slice is ", arr)
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

//GetTypeImports will iterate through attributes and return import statements for each type
func GetTypeImports(w io.Writer, arr []*Attribute) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintln(bw, "import '../../../../core/types/gender.dart';")
	fmt.Fprint(bw, "import '../../../../core/types/vital_status.dart';")
	return bw.Flush()
	// return nil
}

//WriteModelImport writes model import statements
func WriteModelImport(w io.Writer, config *ModelConfig) error {
	entity := config.Entity
	// filePath := ""//filepath.Join(path, fmt.Sprintf("%s%s", i.ModelName, "_model.dart"))
	// f, err := os.Create(filePath)
	// f.Write()
	baseImports := `import 'package:freezed_annotation/freezed_annotation.dart';

%s
import '../../../../core/util/mapper.dart';
import '../../domain/entities/character.dart';

part '%s_model.freezed.dart';
part '%s_model.g.dart';
`
	entityName := strcase.ToSnake(entity.Name)
	typeAttributes := Filter(entity.Attributes, func(val *Attribute) bool {
		// fmt.Println("attribute type is ", val.Meta.AttributeType)
		return val.Meta.AttributeType == "enum"
	})
	importBlock := bytes.Buffer{}
	// importWriter := bufio.NewWriter(&importBlock)
	GetTypeImports(&importBlock, typeAttributes)
	// fmt.Println("enumAttributes are ", typeAttributes)
	fmt.Println("Entity name is ", entity.Name)
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, baseImports, importBlock.String(), entityName, entityName)
	// fmt.Fprint(bw, "world!")
	return bw.Flush() // Don't forget to flush!
	// return nil
}
