# tmpltree

[![Go Tests](https://github.com/btassone/tmpltree/actions/workflows/go.yml/badge.svg)](https://github.com/btassone/tmpltree/actions/workflows/go.yml)

`tmpltree` is a Go package that provides functionality to build and manage a tree structure representation of template directories. It's particularly useful for projects that need to organize and navigate template files in a hierarchical manner.

## Features

- Build a tree structure from a template directory
- Navigate the tree structure programmatically
- Print the tree structure for visualization
- Easily retrieve specific nodes in the tree
- Render templates with support for base layouts

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
    "fmt"
    "log"
    "os"

    "github.com/btassone/tmpltree"
)

func main() {
    // Build the template tree
    root, err := tmpltree.BuildTemplateTree("./templates")
    if err != nil {
        log.Fatalf("Error building template tree: %v", err)
    }

    // Print the tree structure
    fmt.Println("Template Structure:")
    root.Print(os.Stdout, "")

    // Access specific nodes
    if usersPage, ok := root.GetNode("pages", "users"); ok {
        fmt.Println("\nUsers page files:", usersPage.Files)
    }

    if layouts, ok := root.GetNode("layouts"); ok {
        fmt.Println("Layout files:", layouts.Files)
    }

    // Render a template
    var buf bytes.Buffer
    err = root.RenderTemplate("pages/index", "./templates/layouts/base.html", &buf, nil)
    if err != nil {
        log.Fatalf("Error rendering template: %v", err)
    }
    fmt.Println("\nRendered template:\n", buf.String())
}
```

## API Reference

### Types

#### `TemplateNode`

Represents a node in the template tree structure.

```go
type TemplateNode struct {
    Name     string
    Path     string
    Children map[string]*TemplateNode
    Files    []string
}
```

### Functions

#### `NewTemplateNode(name string, path string) *TemplateNode`

Creates a new `TemplateNode`.

#### `BuildTemplateTree(rootDir string) (*TemplateNode, error)`

Constructs a tree structure of the template directory.

### Methods

#### `(n *TemplateNode) Print(w io.Writer, indent string)`

Prints the template tree structure to the provided writer.

#### `(n *TemplateNode) GetNode(path ...string) (*TemplateNode, bool)`

Retrieves a specific node from the template tree.

#### `(n *TemplateNode) RenderTemplate(tmplPath string, baseTemplatePath string, w io.Writer, data interface{}) error`

Renders a template with the given path, using a base template.

## Directory Structure

The package expects the following directory structure for templates:

```
templates/
  ├── layouts/
  ├── pages/
  └── partials/
```

These are the required top-level folders. You can have additional subdirectories and files within these folders.

## Contributing

Contributions to `tmpltree` are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.