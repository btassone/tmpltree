package tmpltree

// This file contains package-level declarations and imports used across the package

// NewTemplateManagerFunc is a function type for creating a new TemplateManager
type NewTemplateManagerFunc func(rootDir string, baseTemplates map[string]string) (*TemplateManager, error)

// NewTemplateManagerImpl is the actual implementation of NewTemplateManager
var NewTemplateManagerImpl NewTemplateManagerFunc = func(rootDir string, baseTemplates map[string]string) (*TemplateManager, error) {
	root, err := BuildTemplateTree(rootDir)
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		Root:          root,
		BaseTemplates: baseTemplates,
	}, nil
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager(rootDir string, baseTemplates map[string]string) (*TemplateManager, error) {
	return NewTemplateManagerImpl(rootDir, baseTemplates)
}
