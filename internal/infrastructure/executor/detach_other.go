//go:build !unix

package executor

import "os/exec"

func detach(*exec.Cmd) {}
