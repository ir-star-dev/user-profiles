package view

import (
	"html/template"
	"path/filepath"
)

// LoadTemplate парсит набор файлов в один template instance
func LoadTemplate(files ...string) (*template.Template, error) {
	t := template.New("")

	for _, file := range files {
		// поддержка glob (например templates/partials/*.tmpl)
		matches, err := filepath.Glob(file)
		if err != nil {
			return nil, err
		}

		if len(matches) == 0 {
			// если это не glob — пробуем как обычный файл
			matches = []string{file}
		}

		_, err = t.ParseFiles(matches...)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}