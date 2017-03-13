package util

import (
	"html/template"
	"path/filepath"
	"os"
	"strings"
)

// LoadTemplates gets and loads all the application views
func LoadTemplates(path string) (*template.Template, error) {
	tpl := template.New("clientList")

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if !strings.HasSuffix(info.Name(), ".html") {
			return nil
		}

		if _, err := tpl.ParseFiles(path); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tpl, nil
}