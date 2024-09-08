package tmpltree

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewTemplateNode(t *testing.T) {
	node := NewTemplateNode("test", "/path/to/test")
	if node.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", node.Name)
	}
	if node.Path != "/path/to/test" {
		t.Errorf("Expected path '/path/to/test', got '%s'", node.Path)
	}
	if len(node.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(node.Children))
	}
	if len(node.Files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(node.Files))
	}
}

func TestBuildTemplateTree(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tmpltree_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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

func TestPrint(t *testing.T) {
	root := NewTemplateNode("root", "/root")
	child1 := NewTemplateNode("child1", "/root/child1")
	child2 := NewTemplateNode("child2", "/root/child2")
	root.Children["child1"] = child1
	root.Children["child2"] = child2
	root.Files = []string{"file1.txt", "file2.txt"}
	child1.Files = []string{"child1file.txt"}

	var buf bytes.Buffer
	root.Print(&buf, "")
	actual := buf.String()

	expected := `root/
  file1.txt
  file2.txt
  child1/
    child1file.txt
  child2/
`

	if actual != expected {
		t.Errorf("Print output doesn't match expected.\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestGetNode(t *testing.T) {
	root := NewTemplateNode("root", "/root")
	child1 := NewTemplateNode("child1", "/root/child1")
	child2 := NewTemplateNode("child2", "/root/child2")
	grandchild := NewTemplateNode("grandchild", "/root/child1/grandchild")
	root.Children["child1"] = child1
	root.Children["child2"] = child2
	child1.Children["grandchild"] = grandchild

	tests := []struct {
		path     []string
		expected *TemplateNode
		ok       bool
	}{
		{[]string{"child1"}, child1, true},
		{[]string{"child2"}, child2, true},
		{[]string{"child1", "grandchild"}, grandchild, true},
		{[]string{"nonexistent"}, nil, false},
		{[]string{"child1", "nonexistent"}, nil, false},
	}

	for _, tt := range tests {
		result, ok := root.GetNode(tt.path...)
		if ok != tt.ok {
			t.Errorf("GetNode(%v) ok = %v, want %v", tt.path, ok, tt.ok)
		}
		if result != tt.expected {
			t.Errorf("GetNode(%v) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

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
