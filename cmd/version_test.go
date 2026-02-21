package cmd

import (
	"bytes"
	"testing"
)

func TestVersionCommandUsesVersionVariable(t *testing.T) {
	prev := Version
	Version = "1.2.3"
	defer func() {
		Version = prev
	}()

	cmd := newVersionCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(nil)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if got, want := out.String(), "cas v1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}
