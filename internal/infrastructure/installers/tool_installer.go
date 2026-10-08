package installers

import (
	"strings"
	"time"

	"github.com/vitualizz/vitualizz-devstack/internal/domain/entities"
	"github.com/vitualizz/vitualizz-devstack/internal/infrastructure/executor"
	"github.com/vitualizz/vitualizz-devstack/internal/infrastructure/logger"
)

type ToolInstaller struct {
	exec      *executor.ShellExecutor
	distro    entities.Distro
	configDir string
	log       *logger.InstallLogger
}

func NewToolInstaller() *ToolInstaller {
	return &ToolInstaller{
		exec:   executor.NewShellExecutor(),
		distro: entities.DetectDistro(),
	}
}

func NewToolInstallerWithDistro(distro entities.Distro) *ToolInstaller {
	return &ToolInstaller{
		exec:   executor.NewShellExecutor(),
		distro: distro,
	}
}

func (i *ToolInstaller) SetLogger(l *logger.InstallLogger) {
	i.log = l
	i.exec.LogFunc = func(toolName, command string, output string, err error, duration time.Duration) {
		if l == nil {
			return
		}
		if err != nil {
			l.LogError(toolName, command, err, output)
		} else {
			l.LogSuccess(toolName, duration)
		}
	}
}

func (i *ToolInstaller) SetConfigDir(dir string) {
	i.configDir = dir
	if dir != "" {
		i.exec.EnvVars = []string{"DEVSTACK_CONFIG=" + dir}
	}
}

func (i *ToolInstaller) SetDetached(detached bool) {
	i.exec.Detached = detached
}

func (i *ToolInstaller) Distro() entities.Distro {
	return i.distro
}

func (i *ToolInstaller) Install(tool *entities.Tool) (*entities.InstallResult, error) {
	start := time.Now()

	result := &entities.InstallResult{
		ToolName: tool.Name,
		Distro:   i.distro,
	}

	cmd := tool.GetInstallCmd(i.distro)
	if cmd == "" {
		result.Success = false
		result.Message = "no install command available for distro: " + string(i.distro)
		return result, nil
	}

	i.exec.ToolName = tool.Name

	if i.log != nil {
		i.log.LogCommand(tool.Name, cmd)
	}

	output, err := i.exec.ExecuteWithOutput(cmd)
	result.DurationMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.Message = formatError(output, err)
		return result, err
	}

	result.Success = true
	result.Message = cleanOutput(output)
	return result, nil
}

func (i *ToolInstaller) Uninstall(tool *entities.Tool) (*entities.InstallResult, error) {
	start := time.Now()

	result := &entities.InstallResult{
		ToolName: tool.Name,
		Distro:   i.distro,
	}

	cmd := tool.GetUninstallCmd(i.distro)
	if cmd == "" {
		result.Success = false
		result.Message = "no uninstall command available for distro: " + string(i.distro)
		return result, nil
	}

	i.exec.ToolName = tool.Name

	if i.log != nil {
		i.log.LogCommand(tool.Name, cmd)
	}

	output, err := i.exec.ExecuteWithOutput(cmd)
	result.DurationMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.Message = formatError(output, err)
		return result, err
	}

	result.Success = true
	result.Message = cleanOutput(output)
	return result, nil
}

func (i *ToolInstaller) IsInstalled(tool *entities.Tool) (bool, error) {
	if tool.Check == "" {
		return false, nil
	}

	i.exec.ToolName = tool.Name

	_, err := i.exec.ExecuteWithOutput(tool.Check)
	return err == nil, nil
}

func formatError(output string, err error) string {
	if output != "" {
		return strings.TrimSpace(output)
	}
	return err.Error()
}

func cleanOutput(output string) string {
	lines := strings.Split(output, "\n")
	var clean []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.Contains(line, "\r") {
			line = strings.SplitN(line, "\r", 2)[1]
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		clean = append(clean, line)
	}

	if len(clean) == 0 {
		return ""
	}

	return strings.Join(clean, "\n")
}

var _ entities.Installer = (*ToolInstaller)(nil)
