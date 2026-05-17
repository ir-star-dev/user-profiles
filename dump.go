package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ignoreDirs = map[string]bool{
	".git":          true,
	"node_modules":  true,
	"vendor":        true,
	".idea":         true,
	".vscode":       true,
	"bin":           true,
	"dist":          true,
	"postgres-data": true,
	"migrate":       true,
	"seeds":         true,
	"ui":			 true,
}

var ignoreExt = map[string]bool{
	".yml":          true,
	".md":           true,
	".mod":          true,
	".sum":          true,
	".env":          true,
	".gitignore":    true,
	".editorconfig": true,
	".env.example":  true,
	".example":      true,
	".tmpl": 		 true,
}

var output strings.Builder

func main() {
	root, _ := os.Getwd()

	// 👉 сначала дерево
	err := buildTree(root, "")
	if err != nil {
		fmt.Println("Tree error:", err)
		return
	}

	output.WriteString("\n\n============================\n")
	output.WriteString("PROJECT DUMP:\n\n")

	err = walk(root, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = os.WriteFile("project_dump.txt", []byte(output.String()), 0644)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("Done → project_dump.txt")
}

func walk(path string, indent string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		name := e.Name()
		fullPath := filepath.Join(path, name)

		if e.IsDir() {
			if ignoreDirs[name] {
				continue
			}

			output.WriteString(fmt.Sprintf("%s📁 %s/\n", indent, name))

			err := walk(fullPath, indent+"    ")
			if err != nil {
				return err
			}
		} else {
			ext := filepath.Ext(name)
			if ignoreExt[ext] {
				continue
			}

			output.WriteString(fmt.Sprintf("%s📄 %s\n", indent, name))
			appendFile(fullPath, indent)
		}
	}

	return nil
}

func appendFile(path string, indent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		output.WriteString(indent + "    [ERROR READING FILE]\n")
		return
	}

	// можно ограничить большие файлы
	if len(content) > 500_000 {
		content = content[:500_000]
	}

	output.WriteString(indent + "    ----- CONTENT START -----\n")
	output.WriteString(indent + string(content) + "\n")
	output.WriteString(indent + "    ----- CONTENT END -----\n\n")
}

func buildTree(path string, indent string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		name := e.Name()
		fullPath := filepath.Join(path, name)

		if e.IsDir() {
			if ignoreDirs[name] {
				continue
			}

			output.WriteString(fmt.Sprintf("%s📁 %s/\n", indent, name))

			err := buildTree(fullPath, indent+"--")
			if err != nil {
				return err
			}
		} else {
			ext := filepath.Ext(name)
			if ignoreExt[ext] {
				continue
			}

			output.WriteString(fmt.Sprintf("%s📄 %s\n", indent, name))
		}
	}

	return nil
}
