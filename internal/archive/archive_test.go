package archive

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
)

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zip")

	original := &Archive{
		Manifest: &Manifest{
			Version:    FormatVersion,
			Agent:      "claude",
			Scope:      "global",
			ExportedAt: time.Now().Truncate(time.Second),
			CASVersion: "0.1.0",
		},
		Instructions: &agent.Instruction{Content: "# My Instructions\nDo stuff."},
		Skills: []agent.Skill{
			{Name: "my-skill", Content: "# My Skill\nSkill content."},
			{Name: "other-skill", Content: "# Other\nMore content."},
		},
	}

	if err := Write(path, original); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if got.Manifest.Version != original.Manifest.Version {
		t.Errorf("manifest version = %q, want %q", got.Manifest.Version, original.Manifest.Version)
	}
	if got.Manifest.Agent != original.Manifest.Agent {
		t.Errorf("manifest agent = %q, want %q", got.Manifest.Agent, original.Manifest.Agent)
	}
	if got.Manifest.Scope != original.Manifest.Scope {
		t.Errorf("manifest scope = %q, want %q", got.Manifest.Scope, original.Manifest.Scope)
	}
	if got.Manifest.CASVersion != original.Manifest.CASVersion {
		t.Errorf("manifest cas_version = %q, want %q", got.Manifest.CASVersion, original.Manifest.CASVersion)
	}
	if !got.Manifest.ExportedAt.Equal(original.Manifest.ExportedAt) {
		t.Errorf("manifest exported_at = %v, want %v", got.Manifest.ExportedAt, original.Manifest.ExportedAt)
	}

	if got.Instructions == nil {
		t.Fatal("expected instructions, got nil")
	}
	if got.Instructions.Content != original.Instructions.Content {
		t.Errorf("instructions = %q, want %q", got.Instructions.Content, original.Instructions.Content)
	}

	if len(got.Skills) != len(original.Skills) {
		t.Fatalf("skills count = %d, want %d", len(got.Skills), len(original.Skills))
	}
	for i, s := range got.Skills {
		if s.Name != original.Skills[i].Name {
			t.Errorf("skill[%d].Name = %q, want %q", i, s.Name, original.Skills[i].Name)
		}
		if s.Content != original.Skills[i].Content {
			t.Errorf("skill[%d].Content = %q, want %q", i, s.Content, original.Skills[i].Content)
		}
	}
}

func TestRoundTripNoInstructions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zip")

	original := &Archive{
		Manifest: &Manifest{
			Version:    FormatVersion,
			Agent:      "copilot",
			Scope:      "local",
			ExportedAt: time.Now().Truncate(time.Second),
		},
		Skills: []agent.Skill{
			{Name: "a-skill", Content: "content"},
		},
	}

	if err := Write(path, original); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if got.Instructions != nil {
		t.Errorf("expected nil instructions, got %+v", got.Instructions)
	}
	if len(got.Skills) != 1 {
		t.Fatalf("skills count = %d, want 1", len(got.Skills))
	}
}

func TestRoundTripNoSkills(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zip")

	original := &Archive{
		Manifest: &Manifest{
			Version:    FormatVersion,
			Agent:      "claude",
			Scope:      "global",
			ExportedAt: time.Now().Truncate(time.Second),
		},
		Instructions: &agent.Instruction{Content: "instructions only"},
	}

	if err := Write(path, original); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if got.Instructions == nil || got.Instructions.Content != "instructions only" {
		t.Errorf("unexpected instructions: %+v", got.Instructions)
	}
	if len(got.Skills) != 0 {
		t.Errorf("expected no skills, got %d", len(got.Skills))
	}
}

func TestReadMissingManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.zip")

	// Create a ZIP with only instructions.md (no manifest)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	fw, err := w.Create("instructions.md")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte("hello"))
	w.Close()
	f.Close()

	_, err = Read(path)
	if err == nil {
		t.Fatal("expected error for missing manifest, got nil")
	}
}

func TestReadUnsupportedManifestVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad-version.zip")

	original := &Archive{
		Manifest: &Manifest{
			Version:    "999",
			Agent:      "claude",
			Scope:      "local",
			ExportedAt: time.Now().Truncate(time.Second),
		},
	}

	if err := Write(path, original); err != nil {
		t.Fatalf("Write: %v", err)
	}

	_, err := Read(path)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
	if !strings.Contains(err.Error(), "unsupported archive format version") {
		t.Fatalf("expected unsupported version error, got %v", err)
	}
}

func TestReadInvalidSkillPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad-skill-path.zip")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)

	manifest, err := w.Create("manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manifest.Write([]byte(`{"version":"1","agent":"claude","scope":"local","exported_at":"2026-01-01T00:00:00Z"}`)); err != nil {
		t.Fatal(err)
	}

	badSkill, err := w.Create("skills/../SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := badSkill.Write([]byte("malicious")); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Read(path)
	if err == nil {
		t.Fatal("expected error for invalid skill path")
	}
	if !strings.Contains(err.Error(), "invalid skill path") {
		t.Fatalf("expected invalid skill path error, got %v", err)
	}
}

func TestWriteInvalidSkillName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad-skill-name.zip")

	err := Write(path, &Archive{
		Manifest: &Manifest{
			Version:    FormatVersion,
			Agent:      "claude",
			Scope:      "local",
			ExportedAt: time.Now().Truncate(time.Second),
		},
		Skills: []agent.Skill{{Name: "../escape", Content: "x"}},
	})
	if err == nil {
		t.Fatal("expected error for invalid skill name")
	}
	if !strings.Contains(err.Error(), "invalid skill name") {
		t.Fatalf("expected invalid skill name error, got %v", err)
	}
}
