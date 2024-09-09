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

## Testing with tmpltree

The `tmpltree` package now includes a `NewTemplateManagerImpl` variable that holds the actual implementation of `NewTemplateManager`. This allows for easier mocking in tests. Here's an example of how to use it in your tests:

```go
import "github.com/btassone/tmpltree"

func TestYourFunction(t *testing.T) {
    // Store the original implementation
    originalNewTemplateManagerImpl := tmpltree.NewTemplateManagerImpl

    // Replace with a mock implementation
    tmpltree.NewTemplateManagerImpl = func(rootDir string, baseTemplates map[string]string) (*tmpltree.TemplateManager, error) {
        return &tmpltree.TemplateManager{}, nil
    }

    // Restore the original implementation at the end of the test
    defer func() { tmpltree.NewTemplateManagerImpl = originalNewTemplateManagerImpl }()

    // Your test code here...
}
```

This approach allows you to easily mock the `NewTemplateManager` function in your tests without modifying the package-level function directly.

## API Reference

(Keep the existing API reference section, adding information about NewTemplateManagerImpl if necessary)

## Directory Structure

(Keep the existing Directory Structure section)

## Contributing

Contributions to `tmpltree` are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.