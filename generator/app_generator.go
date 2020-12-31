package generator

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/iggymacd/blog-generator/config"
)

// // Meta is a data container for Metadata
// type Meta struct {
// 	Title      string
// 	Short      string
// 	Date       string
// 	Tags       []string
// 	ParsedDate time.Time
// }

// // IndexData is a data container for the landing page
// type IndexData struct {
// 	HTMLTitle       string
// 	PageTitle       string
// 	Content         template.HTML
// 	Year            int
// 	Name            string
// 	CanonicalLink   string
// 	MetaDescription string
// 	HighlightCSS    template.CSS
// }

// // Generator interface
// type AppGenerator interface {
// 	Generate() error
// }

// AppGenerator object
type AppGenerator struct {
	Config *AppConfig
}

// AppConfig holds the sources and destination folder
type AppConfig struct {
	Sources     []string
	Destination string
	Config      *config.Config
}

// NewApp creates a new SiteGenerator
func NewApp(config *AppConfig) *AppGenerator {
	return &AppGenerator{Config: config}
}

// Generate starts the app generation
func (g *AppGenerator) Generate() error {
	fmt.Println("sdk is ", g.Config.Config.App.UI.Sdk)
	templatePath := filepath.Join("static", g.Config.Config.App.UI.Sdk, "models", "template.dart")
	fmt.Println("Generating App...")
	sources := g.Config.Sources
	destination := g.Config.Destination
	if err := clearAndCreateDestination(destination); err != nil {
		return err
	}
	if err := clearAndCreateDestination(filepath.Join(destination, "lib")); err != nil {
		return err
	}
	t, err := getModelTemplate(templatePath)
	if err != nil {
		return err
	}
	var entities []*Entity
	for _, path := range sources {
		fmt.Println("path is ", path)
		entity, err := newEntity(path)
		if err != nil {
			return err
		}
		entities = append(entities, entity)
	}
	// sort.Sort(ByDateDesc(entities))
	if err := runAppTasks(entities, t, destination, g.Config.Config); err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

// CreateFlutterApp uses flutter sdk to create new project
func (g *AppGenerator) CreateFlutterApp(cfg *config.Config) error {
	// from := cfg.Generator.Repo
	to := cfg.Generator.Tmp
	appname := cfg.App.Appname
	// branch := cfg.Generator.Branch
	// if branch == "" {
	// 	branch = "master"
	// }

	// fmt.Printf("Fetching data from %s into %s...\n", from, to)
	// if err := createFolderIfNotExist(to); err != nil {
	// 	return nil, err
	// }
	// if err := clearFolder(to); err != nil {
	// 	return nil, err
	// }
	if err := flutterCreate(to, appname); err != nil {
		return err
	}
	return nil
	// dirs, err := getContentFolders(to)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Print("Fetching complete.\n", dirs)
	// return dirs, nil
}

func flutterCreate(path, projectName string) error {
	cmdName := "flutter"
	initArgs := []string{"create", projectName}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating project at %s: %v", path, err)
	}
	// remoteArgs := []string{"remote", "add", "origin", repositoryURL}
	// cmd = exec.Command(cmdName, remoteArgs...)
	// cmd.Dir = path
	// if err := cmd.Run(); err != nil {
	// 	return fmt.Errorf("error setting remote %s: %v", repositoryURL, err)
	// }
	// pullArgs := []string{"pull", "origin", branch}
	// cmd = exec.Command(cmdName, pullArgs...)
	// cmd.Dir = path
	// if err := cmd.Run(); err != nil {
	// 	return fmt.Errorf("error pulling %s at %s: %v", branch, path, err)
	// }
	return nil
}

func getModelTemplate(path string) (*template.Template, error) {
	t, err := template.New("entity").Funcs(funcs).ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %v", path, err)
	}
	return t, nil
}

func runAppTasks(entities []*Entity, t *template.Template, destination string, cfg *config.Config) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}

	// indexWriter := &IndexWriter{
	// 	BlogURL:         cfg.Blog.URL,
	// 	BlogTitle:       cfg.Blog.Title,
	// 	BlogDescription: cfg.Blog.Description,
	// 	BlogAuthor:      cfg.Blog.Author,
	// }

	// modelWriter := &ModelWriter{
	// 	ModelName: cfg.App.,
	// }

	//entities
	for _, entity := range entities {
		pg := ModelGenerator{&ModelConfig{
			Entity:      entity,
			Destination: destination,
			Template:    t,
			// Writer: &ModelWriter{
			// 	ModelName: entity.Name,
			// },
		}}
		generators = append(generators, &pg)
	}
	// tagPostsMap := createTagPostsMap(entities)
	// // frontpage
	// fg := ListingGenerator{&ListingConfig{
	// 	Posts:       entities[:getNumOfPagesOnFrontpage(entities, cfg.Blog.Frontpageposts)],
	// 	Template:    t,
	// 	Destination: destination,
	// 	PageTitle:   "",
	// 	IsIndex:     true,
	// 	Writer:      indexWriter,
	// }}
	// // archive
	// ag := ListingGenerator{&ListingConfig{
	// 	Posts:       entities,
	// 	Template:    t,
	// 	Destination: filepath.Join(destination, "archive"),
	// 	PageTitle:   "Archive",
	// 	IsIndex:     false,
	// 	Writer:      indexWriter,
	// }}
	// // tags
	// tg := TagsGenerator{&TagsConfig{
	// 	TagPostsMap: tagPostsMap,
	// 	Template:    t,
	// 	Destination: destination,
	// 	Writer:      indexWriter,
	// }}

	// staticURLs := []string{}
	// for _, staticURL := range cfg.Blog.Statics.Templates {
	// 	staticURLs = append(staticURLs, staticURL.Dest)
	// }
	// // sitemap
	// sg := SitemapGenerator{&SitemapConfig{
	// 	Posts:       entities,
	// 	TagPostsMap: tagPostsMap,
	// 	Destination: destination,
	// 	BlogURL:     cfg.Blog.URL,
	// 	Statics:     staticURLs,
	// }}
	// // rss
	// rg := RSSGenerator{&RSSConfig{
	// 	Posts:           entities,
	// 	Destination:     destination,
	// 	DateFormat:      cfg.Blog.Dateformat,
	// 	Language:        cfg.Blog.Language,
	// 	BlogURL:         cfg.Blog.URL,
	// 	BlogDescription: cfg.Blog.Description,
	// 	BlogTitle:       cfg.Blog.Title,
	// }}
	// // statics
	// fileToDestination := map[string]string{}
	// for _, static := range cfg.Blog.Statics.Files {
	// 	fileToDestination[static.Src] = filepath.Join(destination, static.Dest)
	// }
	// templateToFile := map[string]string{}
	// for _, static := range cfg.Blog.Statics.Templates {
	// 	templateToFile[static.Src] = filepath.Join(destination, static.Dest, "index.html")
	// }
	// statg := StaticsGenerator{&StaticsConfig{
	// 	FileToDestination: fileToDestination,
	// 	TemplateToFile:    templateToFile,
	// 	Template:          t,
	// 	Writer:            indexWriter,
	// }}
	// generators = append(generators, &fg, &ag, &tg, &sg, &statg)
	// if cfg.Generator.UseRSS {
	// 	generators = append(generators, &rg)
	// }

	for _, generator := range generators {
		wg.Add(1)
		go func(g Generator) {
			defer wg.Done()
			pool <- struct{}{}
			defer func() { <-pool }()
			if err := g.Generate(); err != nil {
				errors <- err
			}
		}(generator)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return nil
	case err := <-errors:
		if err != nil {
			return err
		}
	}
	return nil
}

// func clearAndCreateDestination(path string) error {
// 	if err := os.RemoveAll(path); err != nil {
// 		if !os.IsNotExist(err) {
// 			return fmt.Errorf("error removing folder at destination %s: %v ", path, err)
// 		}
// 	}
// 	return os.Mkdir(path, os.ModePerm)
// }

// IndexWriter writer index.html files
// type IndexWriter struct {
// 	BlogTitle       string
// 	BlogDescription string
// 	BlogAuthor      string
// 	BlogURL         string
// }

// WriteIndexHTML writes an index.html file
// func (i *IndexWriter) WriteIndexHTML(path, pageTitle, metaDescription string, content template.HTML, t *template.Template) error {
// 	filePath := filepath.Join(path, "index.html")
// 	f, err := os.Create(filePath)
// 	if err != nil {
// 		return fmt.Errorf("error creating file %s: %v", filePath, err)
// 	}
// 	defer f.Close()
// 	metaDesc := metaDescription
// 	if metaDescription == "" {
// 		metaDesc = i.BlogDescription
// 	}
// 	hlbuf := bytes.Buffer{}
// 	hlw := bufio.NewWriter(&hlbuf)
// 	formatter := html.New(html.WithClasses(true))
// 	formatter.WriteCSS(hlw, styles.MonokaiLight)
// 	hlw.Flush()
// 	w := bufio.NewWriter(f)
// 	td := IndexData{
// 		Name:            i.BlogAuthor,
// 		Year:            time.Now().Year(),
// 		HTMLTitle:       getHTMLTitle(pageTitle, i.BlogTitle),
// 		PageTitle:       pageTitle,
// 		Content:         content,
// 		CanonicalLink:   buildCanonicalLink(path, i.BlogURL),
// 		MetaDescription: metaDesc,
// 		HighlightCSS:    template.CSS(hlbuf.String()),
// 	}
// 	if err := t.Execute(w, td); err != nil {
// 		return fmt.Errorf("error executing template %s: %v", filePath, err)
// 	}
// 	if err := w.Flush(); err != nil {
// 		return fmt.Errorf("error writing file %s: %v", filePath, err)
// 	}
// 	return nil
// }

// func getHTMLTitle(pageTitle, blogTitle string) string {
// 	if pageTitle == "" {
// 		return blogTitle
// 	}
// 	return fmt.Sprintf("%s - %s", pageTitle, blogTitle)
// }

// func createTagPostsMap(entities []*Post) map[string][]*Post {
// 	result := make(map[string][]*Post)
// 	for _, entity := range entities {
// 		for _, tag := range entity.Meta.Tags {
// 			key := strings.ToLower(tag)
// 			if result[key] == nil {
// 				result[key] = []*Post{entity}
// 			} else {
// 				result[key] = append(result[key], entity)
// 			}
// 		}
// 	}
// 	return result
// }

// func getTemplate(path string) (*template.Template, error) {
// 	t, err := template.ParseFiles(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading template %s: %v", path, err)
// 	}
// 	return t, nil
// }

// func getNumOfPagesOnFrontpage(entities []*Post, numPosts int) int {
// 	if len(entities) < numPosts {
// 		return len(entities)
// 	}
// 	return numPosts
// }

// func buildCanonicalLink(path, baseURL string) string {
// 	parts := strings.Split(path, "/")
// 	if len(parts) > 1 {
// 		return fmt.Sprintf("%s/%s/index.html", baseURL, strings.Join(parts[1:], "/"))
// 	}
// 	return "/"
// }
