package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateData struct {
	ModuleName string // product
	EntityName string // Product
}

func renderTemplate(tplPath string, data TemplateData) (string, error) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func writeFromTemplate(dstPath, tplPath string, data TemplateData) error {
	content, err := renderTemplate(tplPath, data)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(dstPath, []byte(content), 0644)
}
