package cmd

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/LaneBirmingham/coding-agent-sync/internal/agent"
	"github.com/LaneBirmingham/coding-agent-sync/internal/archive"
)

func withCmdGlobals(root string, verbose bool, fn func()) {
	prevRoot := flagRoot
	prevVerbose := flagVerbose
	flagRoot = root
	flagVerbose = verbose
	defer func() {
		flagRoot = prevRoot
		flagVerbose = prevVerbose
	}()
	fn()
}

func TestDoExportInvalidAgent(t *testing.T) {
	withCmdGlobals(t.TempDir(), false, func() {
		err := doExport("not-an-agent", "local", "", true)
		if err == nil {
			t.Fatal("expected error for invalid agent")
		}
	})
}

func TestDoExportInvalidScope(t *testing.T) {
	withCmdGlobals(t.TempDir(), false, func() {
		err := doExport("claude", "wrong", "", true)
		if err == nil {
			t.Fatal("expected error for invalid scope")
		}
	})
}

func TestDoExportDryRun(t *testing.T) {
	root := t.TempDir()
	withCmdGlobals(root, true, func() {
		err := doExport("claude", "local", "", true)
		if err != nil {
			t.Fatalf("expected dry-run export to succeed, got %v", err)
		}
	})
}

func TestDoImportInvalidScope(t *testing.T) {
	withCmdGlobals(t.TempDir(), false, func() {
		err := doImport("claude", "wrong", "input.zip", true)
		if err == nil {
			t.Fatal("expected error for invalid scope")
		}
	})
}

func TestDoImportNoTargets(t *testing.T) {
	withCmdGlobals(t.TempDir(), false, func() {
		err := doImport(" , ", "local", "input.zip", true)
		if err == nil {
			t.Fatal("expected error for empty targets")
		}
	})
}

func TestDoImportInvalidTarget(t *testing.T) {
	withCmdGlobals(t.TempDir(), false, func() {
		err := doImport("unknown", "local", "input.zip", true)
		if err == nil {
			t.Fatal("expected error for invalid target")
		}
	})
}

func TestDoImportDryRun(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(t.TempDir(), "in.zip")

	if err := archive.Write(input, &archive.Archive{
		Manifest: &archive.Manifest{
			Version:    archive.FormatVersion,
			Agent:      "claude",
			Scope:      "local",
			ExportedAt: time.Now().UTC(),
		},
		Instructions: &agent.Instruction{Content: "# A"},
	}); err != nil {
		t.Fatal(err)
	}

	withCmdGlobals(root, true, func() {
		err := doImport("copilot", "local", input, true)
		if err != nil {
			t.Fatalf("expected dry-run import to succeed, got %v", err)
		}
	})
}
