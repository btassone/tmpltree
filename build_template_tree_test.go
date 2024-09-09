package tmpltree

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildTemplateTree(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tmpltree_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
			return
		}
	}(tempDir)

	// Create test directory structure
	createTestDirStructure(t, tempDir)

	// Build the template tree
	root, err := BuildTemplateTree(tempDir)
	if err != nil {
		t.Fatalf("BuildTemplateTree failed: %v", err)
	}

	// Verify the structure
	if len(root.Children) != 3 {
		t.Errorf("Expected 3 top-level children, got %d", len(root.Children))
	}

	expectedStructure := map[string][]string{
		"layouts":  {"base.html"},
		"pages":    {"about.html", "contact.html", "index.html"},
		"partials": {},
	}

	for folder, expectedFiles := range expectedStructure {
		node, ok := root.Children[folder]
		if !ok {
			t.Errorf("Expected folder '%s' not found", folder)
			continue
		}
		if !reflect.DeepEqual(node.Files, expectedFiles) {
			t.Errorf("For folder '%s', expected files %v, got %v", folder, expectedFiles, node.Files)
		}
	}

	// Check for nested structure in pages
	pagesNode := root.Children["pages"]
	if usersNode, ok := pagesNode.Children["users"]; ok {
		if !reflect.DeepEqual(usersNode.Files, []string{"index.html"}) {
			t.Errorf("Expected users/index.html, got %v", usersNode.Files)
		}
	} else {
		t.Error("Expected 'users' folder in pages, not found")
	}
}

// Helper function to create test directory structure
func createTestDirStructure(t *testing.T, root string) {
	dirs := []string{
		"layouts",
		"pages",
		"pages/users",
		"partials",
	}
	files := []string{
		"layouts/base.html",
		"pages/about.html",
		"pages/contact.html",
		"pages/index.html",
		"pages/users/index.html",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(root, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	for _, file := range files {
		err := os.WriteFile(filepath.Join(root, file), []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}
}
