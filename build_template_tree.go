package tmpltree

import (
	"os"
	"path/filepath"
	"strings"
)

// BuildTemplateTree constructs a tree structure of the template directory
func BuildTemplateTree(rootDir string) (*TemplateNode, error) {
	root := NewTemplateNode("templates", rootDir)
	requiredFolders := []string{"layouts", "pages", "partials"}

	for _, folder := range requiredFolders {
		root.Children[folder] = NewTemplateNode(folder, filepath.Join(rootDir, folder))
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) == 1 {
			return nil // Skip the root directory
		}

		current := root.Children[parts[0]]
		if current == nil {
			return nil // Skip if not in a required folder
		}

		for i, part := range parts[1:] {
			if i == len(parts)-2 { // If this is the last part and it's a file
				if !info.IsDir() {
					current.Files = append(current.Files, part)
				}
				break
			}
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = NewTemplateNode(part, filepath.Join(current.Path, part))
			}
			current = current.Children[part]
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return root, nil
}
