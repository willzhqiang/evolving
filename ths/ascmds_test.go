package ths

import "testing"

func TestRunPlainAppleScript(t *testing.T) {
	out, err := Run(`return "ok"`)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out != "ok" {
		t.Fatalf("Run returned %q, want ok", out)
	}
}

func TestRunLegacyWrappedAppleScript(t *testing.T) {
	out, err := Run(`osascript -e 'return "ok"'`)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out != "ok" {
		t.Fatalf("Run returned %q, want ok", out)
	}
}
