package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestGlob(t *testing.T) {
	// Setup a temporary directory structure
	root, err := os.MkdirTemp("", "globtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root) // Clean up

	// Create files and directories
	filesToCreate := []string{
		"node_modules/package.json",
		"dist/app.js",
		".env",
		".env.production",
		"src/app.js",
		"src/components/App.jsx",
		"src/components/App.md",
		"node_modules/src/components/Ignore.md",
		"dist/Ignore.md",
		"src/styles/app.css",
		"test/app.test.js",
		".vscode/settings.json",
		"build/README.md",
		"build/output.txt",
		"README.md",
	}

	for _, file := range filesToCreate {
		filePath := filepath.Join(root, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Create(filePath); err != nil {
			t.Fatal(err)
		}
	}

	// Define patterns to ignore and compile them into regexps
	patternsStr := []string{"node_modules", "dist/.*", ".env", ".env\\..*", "build/.*"}
	var patterns []*regexp.Regexp
	for _, patternStr := range patternsStr {
		pattern, err := regexp.Compile(patternStr)
		if err != nil {
			t.Fatalf("Error compiling regex for pattern %s: %v", patternStr, err)
		}
		patterns = append(patterns, pattern)
	}

	// Expected files to be returned by glob
	expectedFiles := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "src/components/App.md"),
	}

	gotFiles := glob(root, patterns)

	if len(gotFiles) != len(expectedFiles) {
		t.Fatalf("Expected %d files, but got %d files: %v", len(expectedFiles), len(gotFiles), gotFiles)
	}

	// Convert slices to maps for easier comparison
	gotFileMap := make(map[string]bool)
	for _, file := range gotFiles {
		gotFileMap[file] = true
	}

	for _, expectedFile := range expectedFiles {
		if !gotFileMap[expectedFile] {
			t.Errorf("Expected file %s not found in glob results", expectedFile)
		}
	}
}
