package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ShellExecutor struct {
	EnvVars  []string
	ToolName string
	LogFunc  func(toolName, command string, output string, err error, duration time.Duration)
	Detached bool
}

func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

func toolPaths() string {
	return "$HOME/.cargo/bin:$HOME/.local/bin:$HOME/.mise/bin:$HOME/go/bin:$HOME/.opencode/bin:" +
		"/opt/homebrew/bin:/opt/homebrew/sbin:/home/linuxbrew/.linuxbrew/bin:/usr/local/bin"
}

const baseEnv = "export NONINTERACTIVE=1 HOMEBREW_NO_ENV_HINTS=1 && "

func (e *ShellExecutor) wrapCmd(cmd string) string {
	wrapped := fmt.Sprintf("export PATH=%s:$PATH && ", toolPaths()) + baseEnv

	if len(e.EnvVars) > 0 {
		hasRef := false
		for _, env := range e.EnvVars {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && strings.Contains(cmd, "$"+parts[0]) {
				hasRef = true
				break
			}
		}
		if hasRef {
			var exports []string
			for _, env := range e.EnvVars {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					exports = append(exports, "export "+env)
				}
			}
			wrapped += strings.Join(exports, " && ") + " && "
		}
	}

	return wrapped + cmd
}

func (e *ShellExecutor) Execute(cmd string) (string, error) {
	start := time.Now()

	sh, c := shellArgs()
	command := exec.Command(sh, c, e.wrapCmd(cmd))
	if e.Detached {
		detach(command)
	}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	duration := time.Since(start)

	out := stdout.String()
	errOut := stderr.String()

	if e.LogFunc != nil {
		var combinedErr error
		if err != nil {
			combinedErr = fmt.Errorf("%s (stderr: %s)", err, errOut)
		}
		e.LogFunc(e.ToolName, cmd, out+"\n"+errOut, combinedErr, duration)
	}

	if err != nil {
		return "", fmt.Errorf("command failed: %w\n--- stderr ---\n%s", err, errOut)
	}

	return fmt.Sprintf("[%dms] %s", duration.Milliseconds(), out), nil
}

const installTimeout = 15 * time.Minute

func (e *ShellExecutor) ExecuteWithOutput(cmd string) (string, error) {
	return e.run(cmd, true)
}

func (e *ShellExecutor) Check(cmd string) error {
	_, err := e.run(cmd, false)
	return err
}

func (e *ShellExecutor) run(cmd string, logged bool) (string, error) {
	start := time.Now()
	sh, c := shellArgs()

	ctx, cancel := context.WithTimeout(context.Background(), installTimeout)
	defer cancel()

	command := exec.CommandContext(ctx, sh, c, e.wrapCmd(cmd))
	if e.Detached {
		detach(command)
	}
	output, err := command.CombinedOutput()
	duration := time.Since(start)

	if logged && e.LogFunc != nil {
		e.LogFunc(e.ToolName, cmd, string(output), err, duration)
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return string(output), fmt.Errorf("command timed out after %v", installTimeout)
		}
		return string(output), err
	}
	return string(output), nil
}

func shellArgs() (string, string) {
	if isZshAvailable() {
		return "/bin/zsh", "-c"
	}
	return "/bin/sh", "-c"
}

func isZshAvailable() bool {
	_, err := exec.LookPath("zsh")
	return err == nil
}
