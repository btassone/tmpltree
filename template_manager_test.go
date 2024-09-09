package tmpltree

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tmpltree_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestDirStructure(t, tempDir)

	baseTemplates := map[string]string{
		"base":  filepath.Join(tempDir, "layouts", "base.html"),
		"admin": filepath.Join(tempDir, "layouts", "admin.html"),
	}

	// Store the original implementation
	originalNewTemplateManagerImpl := NewTemplateManagerImpl

	// Restore the original implementation at the end of the test
	defer func() { NewTemplateManagerImpl = originalNewTemplateManagerImpl }()

	// Test the actual implementation
	tm, err := NewTemplateManager(tempDir, baseTemplates)
	if err != nil {
		t.Fatalf("NewTemplateManager failed: %v", err)
	}

	if tm.Root == nil {
		t.Error("Root node is nil")
	}

	if !reflect.DeepEqual(tm.BaseTemplates, baseTemplates) {
		t.Errorf("BaseTemplates don't match. Got %v, want %v", tm.BaseTemplates, baseTemplates)
	}

	// Test with a mock implementation
	mockCalled := false
	NewTemplateManagerImpl = func(rootDir string, baseTemplates map[string]string) (*TemplateManager, error) {
		mockCalled = true
		return &TemplateManager{}, nil
	}

	_, err = NewTemplateManager(tempDir, baseTemplates)
	if err != nil {
		t.Fatalf("Mock NewTemplateManager failed: %v", err)
	}

	if !mockCalled {
		t.Error("Mock implementation was not called")
	}
}

func TestTemplateManagerRenderTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tmpltree_render_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestDirStructure(t, tempDir)

	// Create base templates
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
	adminTemplateContent := `
<!DOCTYPE html>
<html>
<head>
    <title>Admin - {{template "title" .}}</title>
</head>
<body>
    <h1>Admin Panel</h1>
    {{template "content" .}}
</body>
</html>
`
	baseTemplatePath := filepath.Join(tempDir, "layouts", "base.html")
	adminTemplatePath := filepath.Join(tempDir, "layouts", "admin.html")

	err = os.WriteFile(baseTemplatePath, []byte(baseTemplateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create base template: %v", err)
	}

	err = os.WriteFile(adminTemplatePath, []byte(adminTemplateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create admin template: %v", err)
	}

	// Create a test page template
	testPagePath := filepath.Join(tempDir, "pages", "test.html")
	testPageContent := `
{{define "title"}}Test Page{{end}}
{{define "content"}}
<h2>Test Page</h2>
<p>Hello, {{.Name}}!</p>
{{end}}
`
	err = os.WriteFile(testPagePath, []byte(testPageContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test page template: %v", err)
	}

	baseTemplates := map[string]string{
		"base":  baseTemplatePath,
		"admin": adminTemplatePath,
	}

	tm, err := NewTemplateManager(tempDir, baseTemplates)
	if err != nil {
		t.Fatalf("NewTemplateManager failed: %v", err)
	}

	tests := []struct {
		name           string
		tmplPath       string
		baseTemplate   string
		data           interface{}
		expectedOutput string
	}{
		{
			name:         "Render with base template",
			tmplPath:     "pages/test",
			baseTemplate: "base",
			data:         struct{ Name string }{"World"},
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    
<h2>Test Page</h2>
<p>Hello, World!</p>

</body>
</html>
`,
		},
		{
			name:         "Render with admin template",
			tmplPath:     "pages/test",
			baseTemplate: "admin",
			data:         struct{ Name string }{"Admin"},
			expectedOutput: `
<!DOCTYPE html>
<html>
<head>
    <title>Admin - Test Page</title>
</head>
<body>
    <h1>Admin Panel</h1>
    
<h2>Test Page</h2>
<p>Hello, Admin!</p>

</body>
</html>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tm.RenderTemplate(tt.tmplPath, tt.baseTemplate, &buf, tt.data)
			if err != nil {
				t.Fatalf("RenderTemplate failed: %v", err)
			}

			if strings.TrimSpace(buf.String()) != strings.TrimSpace(tt.expectedOutput) {
				t.Errorf("RenderTemplate output doesn't match expected.\nGot:\n%s\nExpected:\n%s", buf.String(), tt.expectedOutput)
			}
		})
	}
}

func TestTemplateManagerRenderTemplateErrors(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tmpltree_render_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	createTestDirStructure(t, tempDir)

	baseTemplates := map[string]string{
		"base": filepath.Join(tempDir, "layouts", "base.html"),
	}

	tm, err := NewTemplateManager(tempDir, baseTemplates)
	if err != nil {
		t.Fatalf("NewTemplateManager failed: %v", err)
	}

	tests := []struct {
		name          string
		tmplPath      string
		baseTemplate  string
		expectedError string
	}{
		{
			name:          "Non-existent template",
			tmplPath:      "pages/nonexistent",
			baseTemplate:  "base",
			expectedError: "template file not found: nonexistent.html",
		},
		{
			name:          "Non-existent base template",
			tmplPath:      "pages/index",
			baseTemplate:  "nonexistent",
			expectedError: "base template not found: nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tm.RenderTemplate(tt.tmplPath, tt.baseTemplate, &buf, nil)
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			} else if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', but got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}
