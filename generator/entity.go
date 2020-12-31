package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Entity holds data for a entity
type Entity struct {
	Name       string
	Meta       *EntityMeta
	Attributes []*Attribute
	// HTML      []byte
	// Meta      *Meta
	// ImagesDir string
	// Images    []string
}

// EntityMeta is a data container for Entity Metadata
type EntityMeta struct {
	Description string
	DisplayName string
	// Date       string
	Permissions []string
	Attributes  []*Attribute
	Path        string
	// ParsedDate time.Time
}

// Attribute holds data for an Attribute
type Attribute struct {
	Name string
	Meta *AttributeMeta
	// HTML      []byte
	// Meta      *Meta
	// ImagesDir string
	// Images    []string
}

//Association is a data container for an entity relationship
type Association struct {
	Type        string
	Cardinality string
	Entity      *Entity
}

// AttributeMeta is a data container for Attribute Metadata
type AttributeMeta struct {
	AttributeType string
	DisplayName   string
	Association   *Association
	Path          string
	Required      bool
	// Date       string
	// Permissions []string
	// ParsedDate time.Time
}

// ByDateDesc is the sorting object for entitys
// type ByDateDesc []*Entity

// EntityGenerator object
type EntityGenerator struct {
	Config *EntityConfig
}

// EntityConfig holds the entity's configuration
type EntityConfig struct {
	Entity      *Entity
	Destination string
	Template    *template.Template
	Writer      *IndexWriter
}

// Generate generates a entity
// func (g *EntityGenerator) Generate() error {
// 	entity := g.Config.Entity
// 	destination := g.Config.Destination
// 	t := g.Config.Template
// 	fmt.Printf("\tGenerating Entity: %s...\n", entity.Meta.Title)
// 	staticPath := filepath.Join(destination, entity.Name)
// 	if err := os.Mkdir(staticPath, os.ModePerm); err != nil {
// 		return fmt.Errorf("error creating directory at %s: %v", staticPath, err)
// 	}
// 	if entity.ImagesDir != "" {
// 		if err := copyImagesDir(entity.ImagesDir, staticPath); err != nil {
// 			return err
// 		}
// 	}

// 	if err := g.Config.Writer.WriteIndexHTML(staticPath, entity.Meta.Title, entity.Meta.Short, template.HTML(string(entity.HTML)), t); err != nil {
// 		return err
// 	}
// 	fmt.Printf("\tFinished generating Entity: %s...\n", entity.Meta.Title)
// 	return nil
// }

func newEntity(path string) (*Entity, error) {
	meta, err := getEntityMeta(path)
	if err != nil {
		return nil, err
	}
	// html, err := getHTML(path)
	// if err != nil {
	// 	return nil, err
	// }
	attributes, err := getAttributes(path)
	fmt.Println("attributes are ", attributes)
	// imagesDir, images, err := getImages(path)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)

	return &Entity{Name: name, Meta: meta, Attributes: attributes}, // Meta: meta, HTML: html, ImagesDir: imagesDir, Images: images
		nil
}

// func copyImagesDir(source, destination string) (err error) {
// 	path := filepath.Join(destination, "images")
// 	if err := os.Mkdir(path, os.ModePerm); err != nil {
// 		return fmt.Errorf("error creating images directory at %s: %v", path, err)
// 	}
// 	files, err := ioutil.ReadDir(source)
// 	if err != nil {
// 		return fmt.Errorf("error reading directory %s: %v", path, err)
// 	}
// 	for _, file := range files {
// 		src := filepath.Join(source, file.Name())
// 		dst := filepath.Join(path, file.Name())
// 		if err := copyFile(src, dst); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func getEntityMeta(path string) (*EntityMeta, error) {
	filePath := filepath.Join(path, "meta.yml")
	metaraw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := EntityMeta{}
	err = yaml.Unmarshal(metaraw, &meta)
	if err != nil {
		return nil, fmt.Errorf("error reading yml in %s: %v", filePath, err)
	}
	// parsedDate, err := time.Parse(dateFormat, meta.Date)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing date in %s: %v", filePath, err)
	// }
	// meta.ParsedDate = parsedDate
	return &meta, nil
}

// func getHTML(path string) ([]byte, error) {
// 	filePath := filepath.Join(path, "entity.md")
// 	input, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
// 	}
// 	html := blackfriday.MarkdownCommon(input)
// 	replaced, err := replaceCodeParts(html)
// 	if err != nil {
// 		return nil, fmt.Errorf("error during syntax highlighting of %s: %v", filePath, err)
// 	}
// 	return []byte(replaced), nil

// }

func getAttributes(path string) ([]*Attribute, error) {
	// dirPath := filepath.Join(path, "images")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while reading folder %s: %v", path, err)
	}
	attributes := []*Attribute{}
	for _, file := range files {
		fmt.Println("file is ", file.Name())
		fmt.Println("IsDir is ", file.IsDir())
		if file.IsDir() {
			attrPath := filepath.Join(path, file.Name())
			attribute, _ := getAttribute(file.Name(), attrPath)
			attributes = append(attributes, attribute)
		}
		// images = append(images, file.Name())
	}
	return attributes, nil
}

func getAttribute(attrName, path string) (*Attribute, error) {
	filePath := filepath.Join(path, "meta.yml")
	metaraw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := AttributeMeta{}
	err = yaml.Unmarshal(metaraw, &meta)
	if err != nil {
		return nil, fmt.Errorf("error reading yml in %s: %v", filePath, err)
	}
	// parsedDate, err := time.Parse(dateFormat, meta.Date)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing date in %s: %v", filePath, err)
	// }
	// meta.ParsedDate = parsedDate
	return &Attribute{Name: attrName, Meta: &meta}, nil

	// return nil, nil
}

// func replaceCodeParts(htmlFile []byte) (string, error) {
// 	byteReader := bytes.NewReader(htmlFile)
// 	doc, err := goquery.NewDocumentFromReader(byteReader)
// 	if err != nil {
// 		return "", fmt.Errorf("error while parsing html: %v", err)
// 	}
// 	// find code-parts via css selector and replace them with highlighted versions
// 	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
// 		class, _ := s.Attr("class")
// 		lang := strings.TrimPrefix(class, "language-")
// 		oldCode := s.Text()
// 		lexer := lexers.Get(lang)
// 		formatter := html.New(html.WithClasses(true))
// 		iterator, err := lexer.Tokenise(nil, string(oldCode))
// 		if err != nil {
// 			fmt.Printf("ERROR during syntax highlighting, %v", err)
// 		}
// 		b := bytes.Buffer{}
// 		buf := bufio.NewWriter(&b)
// 		err = formatter.Format(buf, styles.GitHub, iterator)
// 		if err != nil {
// 			fmt.Printf("ERROR during syntax highlighting, %v", err)
// 		}
// 		buf.Flush()
// 		s.SetHtml(b.String())
// 	})
// 	new, err := doc.Html()
// 	if err != nil {
// 		return "", fmt.Errorf("error while generating html: %v", err)
// 	}
// 	// replace unnecessarily added html tags
// 	new = strings.Replace(new, "<html><head></head><body>", "", 1)
// 	new = strings.Replace(new, "</body></html>", "", 1)
// 	return new, nil
// }

// func (p ByDateDesc) Len() int {
// 	return len(p)
// }

// func (p ByDateDesc) Swap(i, j int) {
// 	p[i], p[j] = p[j], p[i]
// }

// func (p ByDateDesc) Less(i, j int) bool {
// 	return p[i].Meta.ParsedDate.After(p[j].Meta.ParsedDate)
// }
