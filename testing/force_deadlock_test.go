package testing

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestForceDeadlock(t *testing.T) {
	exe := filepath.Join(t.TempDir(), "deadlock.exe")

	compile := exec.Command("go", "build", "-o", exe, "-v", "./testdata/deadlock")
	data, err := compile.CombinedOutput()
	s := string(data)
	t.Log(s)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	run := exec.CommandContext(ctx, exe)
	out, err := run.CombinedOutput()
	s = string(out)
	t.Log(s)
	fmt.Println(s)
	if !strings.Contains(s, "deadlock") &&
		!strings.Contains(s, "Deadlock") {
		t.Error("expected output to mention a deadlock")
	}
	if err == nil {
		t.Error("expected process to fail")
	}
}
