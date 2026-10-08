package executor_test

import (
	"strings"
	"testing"
	"time"

	"github.com/vitualizz/dotfiles/internal/infrastructure/executor"
)

func TestExecuteWithOutput_Success(t *testing.T) {
	e := executor.NewShellExecutor()

	tests := []struct {
		name      string
		cmd       string
		wantInOut string
		wantErr   bool
	}{
		{
			name:      "echo returns output",
			cmd:       "echo hello",
			wantInOut: "hello",
			wantErr:   false,
		},
		{
			name:      "multiline output",
			cmd:       "printf 'line1\nline2'",
			wantInOut: "line1",
			wantErr:   false,
		},
		{
			name:    "exit 1 returns error",
			cmd:     "exit 1",
			wantErr: true,
		},
		{
			name:    "nonexistent command returns error",
			cmd:     "command-that-does-not-exist-xyz",
			wantErr: true,
		},
		{
			name:      "stderr is captured in output on error",
			cmd:       "echo errout >&2; exit 1",
			wantInOut: "errout",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := e.ExecuteWithOutput(tt.cmd)
			if tt.wantErr && err == nil {
				t.Errorf("ExecuteWithOutput(%q) error = nil, want error", tt.cmd)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ExecuteWithOutput(%q) error = %v, want nil", tt.cmd, err)
			}
			if tt.wantInOut != "" && !strings.Contains(out, tt.wantInOut) {
				t.Errorf("ExecuteWithOutput(%q) output = %q, want to contain %q", tt.cmd, out, tt.wantInOut)
			}
		})
	}
}

func TestExecute_Success(t *testing.T) {
	e := executor.NewShellExecutor()

	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "successful command",
			cmd:     "echo hello",
			wantErr: false,
		},
		{
			name:    "failing command returns error",
			cmd:     "exit 1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Execute(tt.cmd)
			if tt.wantErr && err == nil {
				t.Errorf("Execute(%q) error = nil, want error", tt.cmd)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Execute(%q) error = %v, want nil", tt.cmd, err)
			}
		})
	}
}

func TestExecute_OutputContainsTime(t *testing.T) {
	e := executor.NewShellExecutor()
	out, err := e.Execute("echo hello")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out, "ms") {
		t.Errorf("Execute() output = %q, want to contain timing [Xms]", out)
	}
}

func TestExecuteWithOutput_ExportsNonInteractive(t *testing.T) {
	e := executor.NewShellExecutor()
	out, err := e.ExecuteWithOutput("echo $NONINTERACTIVE")
	if err != nil {
		t.Fatalf("ExecuteWithOutput() error = %v", err)
	}
	if strings.TrimSpace(out) != "1" {
		t.Errorf("NONINTERACTIVE = %q, want %q", strings.TrimSpace(out), "1")
	}
}

func TestExecuteWithOutput_PathIncludesHomebrew(t *testing.T) {
	e := executor.NewShellExecutor()
	out, err := e.ExecuteWithOutput("echo $PATH")
	if err != nil {
		t.Fatalf("ExecuteWithOutput() error = %v", err)
	}
	if !strings.Contains(out, "/opt/homebrew/bin") {
		t.Errorf("PATH = %q, want to contain /opt/homebrew/bin", out)
	}
}

func TestExecuteWithOutput_DetachedHasNoControllingTTY(t *testing.T) {
	e := executor.NewShellExecutor()
	e.Detached = true
	if _, err := e.ExecuteWithOutput("exec 3</dev/tty"); err == nil {
		t.Error("detached command could open /dev/tty, want failure")
	}
}

func TestCheck_DoesNotLog(t *testing.T) {
	e := executor.NewShellExecutor()
	logged := 0
	e.LogFunc = func(string, string, string, error, time.Duration) { logged++ }

	if err := e.Check("exit 1"); err == nil {
		t.Error("Check(exit 1) error = nil, want error")
	}
	if err := e.Check("true"); err != nil {
		t.Errorf("Check(true) error = %v, want nil", err)
	}
	if logged != 0 {
		t.Errorf("Check logged %d times, want 0", logged)
	}

	if _, err := e.ExecuteWithOutput("true"); err != nil {
		t.Fatal(err)
	}
	if logged != 1 {
		t.Errorf("ExecuteWithOutput logged %d times, want 1", logged)
	}
}
