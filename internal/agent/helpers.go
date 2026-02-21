package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// readFile reads the contents of a file, returning empty string and nil if the file doesn't exist.
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// writeFile writes content to a file, creating parent directories as needed.
func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// readSkillsFromDir reads all SKILL.md files from subdirectories of dir.
func readSkillsFromDir(dir string) ([]Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var skills []Skill
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name(), "SKILL.md")
		content, err := readFile(path)
		if err != nil {
			return nil, err
		}
		if content != "" {
			skills = append(skills, Skill{
				Name:    e.Name(),
				Content: content,
			})
		}
	}
	return skills, nil
}

// writeSkillsToDir writes skills to subdirectories of dir.
func writeSkillsToDir(dir string, skills []Skill) error {
	for _, s := range skills {
		if err := validateSkillDirName(s.Name); err != nil {
			return fmt.Errorf("invalid skill name %q: %w", s.Name, err)
		}
		path := filepath.Join(dir, s.Name, "SKILL.md")
		if err := writeFile(path, s.Content); err != nil {
			return err
		}
	}
	return nil
}

func validateSkillDirName(name string) error {
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("skill name must be a non-empty directory name")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("skill name must not contain path separators")
	}
	if strings.ContainsRune(name, '\x00') {
		return fmt.Errorf("skill name must not contain NUL")
	}
	return nil
}
