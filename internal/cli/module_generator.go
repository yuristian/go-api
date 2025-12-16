package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateModule(name string) error {
	module := strings.ToLower(strings.TrimSpace(name))

	// ----------------------------
	// Validation
	// ----------------------------
	if module == "" {
		return errors.New("module name cannot be empty")
	}

	if strings.Contains(module, " ") {
		return errors.New("module name must not contain spaces")
	}

	basePath := filepath.Join("internal", "modules", module)

	if _, err := os.Stat(basePath); err == nil {
		return fmt.Errorf("module '%s' already exists", module)
	}

	fmt.Println("[CLI] Creating module:", module)

	// ----------------------------
	// Folder structure
	// ----------------------------
	dirs := []string{
		"domain",
		"infrastructure",
		"presentation",
		"usecase",
	}

	for _, dir := range dirs {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	// ----------------------------
	// Template data
	// ----------------------------
	data := TemplateData{
		ModuleName: module,
		EntityName: capitalize(module),
	}

	// ----------------------------
	// Generate files from templates
	// ----------------------------
	err := writeFromTemplate(
		filepath.Join(basePath, "domain", "entity.go"),
		"internal/cli/templates/domain_entity.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	err = writeFromTemplate(
		filepath.Join(basePath, "domain", "repository.go"),
		"internal/cli/templates/domain_repository.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	err = writeFromTemplate(
		filepath.Join(basePath, "usecase", "usecase.go"),
		"internal/cli/templates/usecase.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	err = writeFromTemplate(
		filepath.Join(basePath, "infrastructure", "repository_gorm.go"),
		"internal/cli/templates/infra_repository_gorm.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	err = writeFromTemplate(
		filepath.Join(basePath, "presentation", "handler.go"),
		"internal/cli/templates/presentation_handler.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	err = writeFromTemplate(
		filepath.Join(basePath, "infrastructure", "routes.go"),
		"internal/cli/templates/infra_routes.go.tpl",
		data,
	)
	if err != nil {
		return err
	}

	fmt.Println("[CLI] Module generated successfully ✔")
	return nil
}

// ----------------------------
// Helpers
// ----------------------------
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
