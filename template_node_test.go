package tmpltree

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
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
