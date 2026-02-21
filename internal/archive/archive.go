package archive

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
)

// FormatVersion is the current archive format version.
const FormatVersion = "1"

// Manifest holds metadata about an exported archive.
type Manifest struct {
	Version    string    `json:"version"`
	Agent      string    `json:"agent"`
	Scope      string    `json:"scope"`
	ExportedAt time.Time `json:"exported_at"`
	CASVersion string    `json:"cas_version,omitempty"`
}

// Archive represents the contents of an export archive.
type Archive struct {
	Manifest     *Manifest
	Instructions *agent.Instruction
	Skills       []agent.Skill
}

// Write creates a ZIP archive at path from the given Archive.
func Write(path string, a *Archive) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating archive file: %w", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// Write manifest.json
	manifestData, err := json.MarshalIndent(a.Manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}
	if err := writeEntry(w, "manifest.json", manifestData); err != nil {
		return err
	}

	// Write instructions.md if present
	if a.Instructions != nil && a.Instructions.Content != "" {
		if err := writeEntry(w, "instructions.md", []byte(a.Instructions.Content)); err != nil {
			return err
		}
	}

	// Write skills
	for _, s := range a.Skills {
		entryPath := fmt.Sprintf("skills/%s/SKILL.md", s.Name)
		if err := writeEntry(w, entryPath, []byte(s.Content)); err != nil {
			return err
		}
	}

	return nil
}

// Read loads an Archive from a ZIP file at path.
func Read(path string) (*Archive, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("opening archive: %w", err)
	}
	defer r.Close()

	a := &Archive{}
	foundManifest := false

	for _, f := range r.File {
		name := f.Name

		switch {
		case name == "manifest.json":
			data, err := readEntry(f)
			if err != nil {
				return nil, fmt.Errorf("reading manifest: %w", err)
			}
			var m Manifest
			if err := json.Unmarshal(data, &m); err != nil {
				return nil, fmt.Errorf("parsing manifest: %w", err)
			}
			a.Manifest = &m
			foundManifest = true

		case name == "instructions.md":
			data, err := readEntry(f)
			if err != nil {
				return nil, fmt.Errorf("reading instructions: %w", err)
			}
			a.Instructions = &agent.Instruction{Content: string(data)}

		case strings.HasPrefix(name, "skills/") && strings.HasSuffix(name, "/SKILL.md"):
			// Extract skill name from "skills/<name>/SKILL.md"
			parts := strings.Split(name, "/")
			if len(parts) != 3 {
				continue
			}
			data, err := readEntry(f)
			if err != nil {
				return nil, fmt.Errorf("reading skill %s: %w", parts[1], err)
			}
			a.Skills = append(a.Skills, agent.Skill{
				Name:    parts[1],
				Content: string(data),
			})
		}
	}

	if !foundManifest {
		return nil, fmt.Errorf("archive missing manifest.json")
	}

	return a, nil
}

func writeEntry(w *zip.Writer, name string, data []byte) error {
	fw, err := w.Create(name)
	if err != nil {
		return fmt.Errorf("creating entry %s: %w", name, err)
	}
	if _, err := fw.Write(data); err != nil {
		return fmt.Errorf("writing entry %s: %w", name, err)
	}
	return nil
}

func readEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}
