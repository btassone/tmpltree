package tmpltree

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tmpltree_render_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(tempDir)

	// Create test directory structure
	createTestDirStructure(t, tempDir)

	// Create a base template
	baseTemplatePath := filepath.Join(tempDir, "layouts", "base.html")
	baseTemplateContent := `
<!DOCTYPE html>
<html>
<head>
    <title>{{template "title" .}}</title>
</head>
<body>
    {{template "content" .}}
</body>
</html>
`
	err = os.WriteFile(baseTemplatePath, []byte(baseTemplateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create base template: %v", err)
	}

	// Create a test page template
	testPagePath := filepath.Join(tempDir, "pages", "test.html")
	testPageContent := `
{{define "title"}}Test Page{{end}}
{{define "content"}}
<h1>Test Page</h1>
<p>Hello, {{.Name}}!</p>
{{end}}
`
	err = os.WriteFile(testPagePath, []byte(testPageContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test page template: %v", err)
	}

	// Build the template tree
	root, err := BuildTemplateTree(tempDir)
	if err != nil {
		t.Fatalf("BuildTemplateTree failed: %v", err)
	}

	// Test rendering the template
	var buf bytes.Buffer
	err = root.RenderTemplate("pages/test", baseTemplatePath, &buf, struct{ Name string }{"World"})
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    
<h1>Test Page</h1>
<p>Hello, World!</p>

</body>
</html>
`

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Errorf("RenderTemplate output doesn't match expected.\nGot:\n%s\nExpected:\n%s", buf.String(), expected)
	}
}

func TestRenderTemplateErrors(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tmpltree_render_error_test")
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

	tests := []struct {
		name             string
		tmplPath         string
		baseTemplatePath string
		expectedError    string
	}{
		{
			name:             "Non-existent template",
			tmplPath:         "pages/nonexistent",
			baseTemplatePath: filepath.Join(tempDir, "layouts", "base.html"),
			expectedError:    "template file not found: nonexistent.html",
		},
		{
			name:             "Non-existent base template",
			tmplPath:         "pages/index",
			baseTemplatePath: filepath.Join(tempDir, "layouts", "nonexistent.html"),
			expectedError:    "error parsing template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := root.RenderTemplate(tt.tmplPath, tt.baseTemplatePath, &buf, nil)
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			} else if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', but got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}

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

	// Split the output into lines and sort them
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	sort.Strings(actualLines)

	expected := `
root/
  file1.txt
  file2.txt
  child1/
    child1file.txt
  child2/
`
	// Split the expected output into lines and sort them
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")
	sort.Strings(expectedLines)

	// Compare the sorted lines
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Errorf("Print output doesn't match expected.\nGot:\n%s\nExpected:\n%s",
			strings.Join(actualLines, "\n"),
			strings.Join(expectedLines, "\n"))
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
