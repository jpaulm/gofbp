package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestForceDeadlock(t *testing.T) {
	exe := filepath.Join(t.TempDir(), "deadlock.exe")

	compile := exec.Command("go", "build", "-o", exe, "./testdata/deadlock")
	data, err := compile.CombinedOutput()
	t.Log(string(data))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	run := exec.CommandContext(ctx, exe)
	out, err := run.CombinedOutput()
	t.Log(string(out))
	if !strings.Contains(string(out), "deadlock") &&
		!strings.Contains(string(out), "Deadlock") {
		t.Error("expected output to mention a deadlock")
	}
	if err == nil {
		t.Error("expected process to fail")
	}
}
