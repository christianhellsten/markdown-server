package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
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

var (
	htmlTmpl string
	menuTmpl string
)

func loadTemplate(filePath string, defaultContent string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Using default template as %s was not found", filePath)
		return defaultContent
	}
	return string(content)
}

func init() {
	htmlTmpl = loadTemplate(".index.html.template", defaultHtmlTmpl)
	menuTmpl = loadTemplate(".menu.html.template", defaultMenuTmpl)
}

// Handler serves .md files as HTML and file system contents in a menu
func Handler(w http.ResponseWriter, r *http.Request) {
	// Only handle GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate and sanitize the requested path
	requestedPath := filepath.Clean("." + r.URL.Path)
	if strings.HasPrefix(requestedPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Check if the path is a directory
	if stat, err := os.Stat(requestedPath); err == nil && stat.IsDir() {
		renderDirectory(w, requestedPath, r.URL.Path)
		return
	}

	// Check if the path is a markdown file
	if filepath.Ext(requestedPath) == ".md" {
		renderMarkdownFile(w, requestedPath, r.URL.Path)
		return
	}

	// Check if the path is an image file
	if isImageFile(requestedPath) {
		http.ServeFile(w, r, requestedPath)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// isImageFile checks if the file has an image extension
func isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg":
		return true
	default:
		return false
	}
}

// renderMarkdownFile renders the given Markdown file as HTML
func renderMarkdownFile(w http.ResponseWriter, filePath string, urlPath string) {
	// Read the Markdown file
	mdContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Convert Markdown to HTML
	htmlContent := blackfriday.Run(mdContent)

	// Generate the menu
	menu := generateMenu(".", urlPath)

	// Use template to inject the HTML content
	t, err := template.New("markdown").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(htmlTmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Render the template
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
		Content:  string(htmlContent),
		Menu:     template.HTML(menu),
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// renderDirectory renders a list of Markdown files and directories in the given directory
func renderDirectory(w http.ResponseWriter, dirPath string, urlPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", dirPath, err)
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	// Create a list of links to .md files and directories
	var mdFiles []string
	var directories []string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		} else if filepath.Ext(file.Name()) == ".md" {
			mdFiles = append(mdFiles, file.Name())
		}
	}

	// Generate the menu
	menu := generateMenu(".", urlPath)

	// Use template to render the directory listing
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
		Title:    "",
		FilePath: "",
		UrlPath:  "",
		Content:  "Hello 👋",
		Menu:     template.HTML(menu),
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Custom glob function to filter files
func glob(root string, fn func(string) bool) []string {
	var files []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if fn(s) {
			files = append(files, s)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking the path %s: %v", root, err)
	}
	return files
}

// generateMenu generates an HTML nested list of files and directories using a markdown template
func generateMenu(root string, urlPath string) string {
	// Define the filter function to include only .md files
	mdFilter := func(path string) bool {
		if info, err := os.Stat(path); err == nil && !info.IsDir() && filepath.Ext(path) == ".md" {
			return true
		}
		return false
	}

	files := glob(root, mdFilter)

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
	// Serve static files and handle Markdown requests
	http.HandleFunc("/", Handler)

	// Start the server
	port := 8080
	log.Printf("Serving on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
