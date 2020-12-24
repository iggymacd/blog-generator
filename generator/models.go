package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
)

// ModelData holds the data for the model definition
type ModelData struct {
	Name string
	// Date       string
	// Short      string
	// Link       string
	// TimeToRead string
	ImportStatements, StaticFields, Fields, Transformer string
}

// // Field holds the data for each field
// type Field struct {
// 	Name string
// }

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
	Writer *ModelWriter
}

// ModelWriter writes model files
type ModelWriter struct {
	ModelName string
	// BlogDescription string
	// BlogAuthor      string
	// BlogURL         string
}

// WriteModel writes a Model.dart file
func (i *ModelWriter) WriteModel(path, importStatements, staticFields, fields, transformer string, t *template.Template) error {
	filePath := filepath.Join(path, fmt.Sprintf("%s%s", i.ModelName, "Model.dart"))
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	// metaDesc := metaDescription
	// if metaDescription == "" {
	// 	metaDesc = i.BlogDescription
	// }
	// hlbuf := bytes.Buffer{}
	// hlw := bufio.NewWriter(&hlbuf)
	// formatter := html.New(html.WithClasses(true))
	// formatter.WriteCSS(hlw, styles.MonokaiLight)
	// hlw.Flush()
	w := bufio.NewWriter(f)
	td := ModelData{
		ImportStatements: importStatements,
		StaticFields:     staticFields,
		Fields:           fields,
		Transformer:      transformer,
		Name:             i.ModelName,
	}
	// td := IndexData{
	// 	Name:            i.BlogAuthor,
	// 	Year:            time.Now().Year(),
	// 	HTMLTitle:       getHTMLTitle(pageTitle, i.BlogTitle),
	// 	PageTitle:       pageTitle,
	// 	Content:         content,
	// 	CanonicalLink:   buildCanonicalLink(path, i.BlogURL),
	// 	MetaDescription: metaDesc,
	// 	HighlightCSS:    template.CSS(hlbuf.String()),
	// }
	if err := t.Execute(w, td); err != nil {
		return fmt.Errorf("error executing template %s: %v", filePath, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

// Generate generates a model
func (g *ModelGenerator) Generate() error {
	entity := g.Config.Entity
	destination := g.Config.Destination
	t := g.Config.Template
	fmt.Printf("\tGenerating Model: %v%s...\n", entity.Name, "Model")
	staticPath := filepath.Join(destination, entity.Name)
	if err := os.Mkdir(staticPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory at %s: %v", staticPath, err)
	}
	const (
		// master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		fields       = `{{range .}} {{"\n"}}{{"\t"}}String {{"_"}}{{lowerCamel .Name}}{{";"}}{{end}} `
		staticFields = `{{range .}} {{"\n"}}{{"\t"}}static const {{toScreamingSnake .Name}}{{" = \""}}{{lowerCamel .Name}}{{"\";"}}{{end}} `
	)
	var (
		funcs = template.FuncMap{
			"lowerCamel":       strcase.ToLowerCamel,
			"toScreamingSnake": strcase.ToScreamingSnake,
		}
		// guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)

	fieldsTemplate, err := template.New("fields").Funcs(funcs).Parse(fields)
	if err != nil {
		return fmt.Errorf("error creating template %s: %v", "fields", err)
	}
	staticFieldsTemplate, err := template.New("static_fields").Funcs(funcs).Parse(staticFields)
	if err != nil {
		return fmt.Errorf("error creating template %s: %v", "fields", err)
	}
	fieldsBlock := bytes.Buffer{}
	err = fieldsTemplate.Execute(&fieldsBlock, entity.Attributes)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", "fields", err)
	}
	staticFieldsBlock := bytes.Buffer{}
	err = staticFieldsTemplate.Execute(&staticFieldsBlock, entity.Attributes)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", "fields", err)
	}
	// if post.ImagesDir != "" {
	// 	if err := copyImagesDir(post.ImagesDir, staticPath); err != nil {
	// 		return err
	// 	}
	// }
	if err := g.Config.Writer.WriteModel(staticPath, "imports", staticFieldsBlock.String(), fieldsBlock.String(), "transformers", t); err != nil {
		return err
	}
	// if err := g.Config.Writer.WriteIndexHTML(staticPath, post.Meta.Title, post.Meta.Short, template.HTML(string(post.HTML)), t); err != nil {
	// 	return err
	// }
	fmt.Printf("\tFinished generating Model: %s...\n", entity.Name)
	return nil
}
