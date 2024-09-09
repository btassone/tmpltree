# tmpltree

[![Go Tests](https://github.com/btassone/tmpltree/actions/workflows/go.yml/badge.svg)](https://github.com/btassone/tmpltree/actions/workflows/go.yml)

`tmpltree` is a Go package that provides functionality to build and manage a tree structure representation of template directories. It's particularly useful for projects that need to organize and navigate template files in a hierarchical manner.

## Features

- Build a tree structure from a template directory
- Navigate the tree structure programmatically
- Print the tree structure for visualization
- Easily retrieve specific nodes in the tree
- Render templates with support for multiple base layouts
- Manage templates with a `TemplateManager`

## Installation

To use `tmpltree` in your Go project, you can install it using `go get`:

```bash
go get github.com/btassone/tmpltree
```

## Usage

Here's a quick example of how to use the `tmpltree` package:

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/btassone/tmpltree"
)

func main() {
	// Define base templates
	baseTemplates := map[string]string{
		"base":  "./templates/layouts/base.html",
		"admin": "./templates/layouts/admin.html",
	}

	// Create a new TemplateManager
	tm, err := tmpltree.NewTemplateManager("./templates", baseTemplates)
	if err != nil {
		log.Fatalf("Error creating TemplateManager: %v", err)
	}

	// Print the tree structure
	fmt.Println("Template Structure:")
	tm.Root.Print(os.Stdout, "")

	// Access specific nodes
	if usersPage, ok := tm.Root.GetNode("pages", "users"); ok {
		fmt.Println("\nUsers page files:", usersPage.Files)
	}

	if layouts, ok := tm.Root.GetNode("layouts"); ok {
		fmt.Println("Layout files:", layouts.Files)
	}

	// Render a template
	var buf bytes.Buffer
	err = tm.RenderTemplate("pages/index", "base", &buf, nil)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}
	fmt.Println("\nRendered template:\n", buf.String())
}
```

## Package Structure

The `tmpltree` package is now organized into several files for better modularity:

- `tmpltree.go`: Main package file with core types and functions
- `template_node.go`: Contains the `TemplateNode` struct and its methods
- `template_manager.go`: Contains the `TemplateManager` struct and its methods
- `build_template_tree.go`: Contains the `BuildTemplateTree` function

## Testing

The package includes comprehensive tests for each component. Test files are organized to match the package structure:

- `template_node_test.go`
- `template_manager_test.go`
- `build_template_tree_test.go`

To run the tests, use the following command:

```bash
go test ./...
```

## API Reference

### TemplateNode

- `NewTemplateNode(name string, path string) *TemplateNode`: Creates a new TemplateNode
- `(n *TemplateNode) Print(w io.Writer, indent string)`: Prints the template tree structure
- `(n *TemplateNode) GetNode(path ...string) (*TemplateNode, bool)`: Retrieves a specific node from the tree

### TemplateManager

- `NewTemplateManager(rootDir string, baseTemplates map[string]string) (*TemplateManager, error)`: Creates a new TemplateManager
- `(tm *TemplateManager) RenderTemplate(tmplPath string, baseTemplateName string, w io.Writer, data interface{}) error`: Renders a template with the given path and base template

### Utility Functions

- `BuildTemplateTree(rootDir string) (*TemplateNode, error)`: Constructs a tree structure of the template directory

## Integration with Web Frameworks

The `tmpltree` package can be easily integrated with popular Go web frameworks. Here's an example of how to use it with Gin:

```go
import (
	"github.com/gin-gonic/gin"
	"github.com/btassone/tmpltree"
)

func main() {
	r := gin.Default()
	
	tm, err := tmpltree.NewTemplateManager("./templates", map[string]string{"base": "./templates/layouts/base.html"})
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		var buf bytes.Buffer
		err := tm.RenderTemplate("pages/index", "base", &buf, gin.H{"title": "Welcome"})
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.Data(200, "text/html", buf.Bytes())
	})

	r.Run(":8080")
}
```

## Contributing

Contributions to `tmpltree` are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.