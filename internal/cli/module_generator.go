package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateModule(moduleName string) error {
	if moduleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	fmt.Println("[CLI] Generating module:", moduleName)

	// 1. Create folder structure
	if err := createModuleDirectories(moduleName); err != nil {
		return err
	}

	// 2. Generate files from templates
	if err := generateModuleFiles(moduleName); err != nil {
		return err
	}

	// 3. Auto-register module (AST-based)
	if err := AutoRegisterModule(moduleName); err != nil {
		return err
	}

	fmt.Println("[CLI] Module generated and registered successfully")
	return nil
}

/* =========================
   DIRECTORY GENERATION
   ========================= */

func createModuleDirectories(module string) error {
	basePath := filepath.Join("internal", "modules", module)

	dirs := []string{
		filepath.Join(basePath, "domain"),
		filepath.Join(basePath, "usecase"),
		filepath.Join(basePath, "infrastructure"),
		filepath.Join(basePath, "presentation"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

/* =========================
   FILE GENERATION
   ========================= */

func generateModuleFiles(module string) error {
	data := TemplateData{
		ModuleName: module,             // product
		EntityName: capitalize(module), // Product
	}

	baseTpl := "internal/cli/templates"
	baseDst := filepath.Join("internal", "modules", module)

	files := []struct {
		Tpl string
		Dst string
	}{
		{"domain_entity.go.tpl", "domain/entity.go"},
		{"domain_repository.go.tpl", "domain/repository.go"},
		{"usecase.go.tpl", "usecase/usecase.go"},
		{"infra_repository_gorm.go.tpl", "infrastructure/repository_gorm.go"},
		{"infra_routes.go.tpl", "infrastructure/routes.go"},
		{"presentation_handler.go.tpl", "presentation/handler.go"},
	}

	for _, f := range files {
		if err := writeFromTemplate(
			filepath.Join(baseDst, f.Dst),
			filepath.Join(baseTpl, f.Tpl),
			data,
		); err != nil {
			return err
		}
	}

	return nil
}
