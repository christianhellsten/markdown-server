package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

const (
	defaultHtmlTmpl = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{ .Title }}</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.slate.min.css">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/default.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
		<script>hljs.highlightAll();</script>
	</head>
	<body>
		<main class="container">
			<header>{{ .Menu }}</header>
			<article>
				{{ .Content | safeHTML }}
			</article>
		</main>
	</body>
	</html>
	`
	defaultMenuTmpl = `
	<details class="dropdown">
		<summary role="button" class="contrast">📁 {{ .UrlPath }}</summary>
		{{ .Menu | safeHTML }}
	</details>
	`
)

var ignorePatterns = parseIgnoreFile(".markdown-server/ignore")

var (
	htmlTmpl string
	menuTmpl string
)

func loadTemplate(filePath string, defaultContent string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("using default template as %s was not found", filePath)
		return defaultContent
	}
	return string(content)
}

func init() {
	htmlTmpl = loadTemplate(".markdown-server/index.html", defaultHtmlTmpl)
	menuTmpl = loadTemplate(".markdown-server/menu.html", defaultMenuTmpl)
}

func fileHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cleanPath := filepath.Clean(r.URL.Path)
		requestedPath := filepath.Join(baseDir, cleanPath)

		if !strings.HasPrefix(requestedPath, baseDir) {
			http.Error(w, fmt.Sprintf("Invalid path %s", requestedPath), http.StatusBadRequest)
			return
		}

		// Ignore files
		if matchPatterns(requestedPath, ignorePatterns) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if r.URL.Path == "/" {
			renderHomepage(w, requestedPath, r.URL.Path)
			return
		}

		if filepath.Ext(requestedPath) == ".md" {
			renderMarkdownFile(w, baseDir, requestedPath, r.URL.Path)
			return
		}

		if isImageFile(requestedPath) {
			http.ServeFile(w, r, requestedPath)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg", ".ico", ".webm":
		return true
	default:
		return false
	}
}

func renderMarkdownFile(w http.ResponseWriter, baseDir string, filePath string, urlPath string) {
	mdContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	var buf strings.Builder
	if err := goldmark.Convert(mdContent, &buf); err != nil {
		log.Printf("Error converting markdown: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	htmlContent := buf.String()
	menu := generateMenu(baseDir, urlPath)

	t, err := template.New("markdown").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(htmlTmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		FilePath string
		UrlPath  string
		Content  string
		Menu     template.HTML
	}{
		Title:    filepath.Base(filePath),
		UrlPath:  urlPath,
		FilePath: filePath,
		Content:  htmlContent,
		Menu:     template.HTML(menu),
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func renderHomepage(w http.ResponseWriter, dirPath string, urlPath string) {
	menu := generateMenu(dirPath, urlPath)

	t, err := template.New("directory").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(htmlTmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		FilePath string
		UrlPath  string
		Content  string
		Menu     template.HTML
	}{
		Title:    "Home",
		FilePath: "",
		UrlPath:  urlPath,
		Content:  "Hello 👋",
		Menu:     template.HTML(menu),
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func parseIgnoreFile(filename string) []*regexp.Regexp {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading ignore file: %v", err)
		return nil
	}
	lines := strings.Split(string(data), "\n")
	var patterns []*regexp.Regexp
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			pattern, err := regexp.Compile(line)
			if err != nil {
				log.Printf("Error compiling regex for pattern %s: %v", line, err)
				continue
			}
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

func glob(root string, patterns []*regexp.Regexp) []string {
	var files []string
	err := filepath.WalkDir(root, func(filePath string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(filePath) == ".md" && !matchPatterns(filePath, patterns) {
			files = append(files, filePath)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking the path %s: %v", root, err)
	}
	return files
}

// matchPatterns checks if the given path matches any of the provided regex patterns
func matchPatterns(path string, patterns []*regexp.Regexp) bool {
	relPath, err := filepath.Rel("/", path)
	if err != nil {
		log.Printf("Error getting relative path: %v", err)
		return false
	}

	for _, pattern := range patterns {
		if pattern.MatchString(relPath) {
			return true
		}
	}
	return false
}

func ignorePath(filePath string, patterns []*regexp.Regexp) bool {
	return matchPatterns(filePath, patterns)
}

func generateMenu(root string, urlPath string) string {
	files := glob(root, ignorePatterns)

	var sb strings.Builder
	sb.WriteString("<ul>")

	for _, file := range files {
		relPath, _ := filepath.Rel(root, file)
		sb.WriteString(fmt.Sprintf("<li class=\"file\"><a href=\"/%s\">📄 %s</a></li>", relPath, relPath))
	}

	sb.WriteString("</ul>")

	mdMenu := struct {
		UrlPath string
		Menu    string
	}{
		UrlPath: urlPath,
		Menu:    sb.String(),
	}

	tmpl, err := template.New("menu").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(menuTmpl)
	if err != nil {
		log.Printf("Error parsing markdown menu template: %v", err)
		return ""
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, mdMenu); err != nil {
		log.Printf("Error executing markdown menu template: %v", err)
		return ""
	}

	return result.String()
}

func main() {
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.Int("port", 8080, "Port to listen on")
	baseDir := flag.String("dir", ".", "Base directory to serve files from")
	flag.Parse()

	absBaseDir, err := filepath.Abs(*baseDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for base directory: %v", err)
	}

	http.HandleFunc("/", fileHandler(absBaseDir))

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Serving %s on http://%s\n", absBaseDir, address)
	log.Fatal(http.ListenAndServe(address, nil))
}
