package entities

import (
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

var DockerTools = map[string]bool{
	"kitty":          true,
	"kitty-config":   true,
	"docker":         true,
	"docker-compose": true,
	"lazydocker":     true,
	"docker-filemanager": true,
}

func IsDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := strings.ToLower(string(data))
		if strings.Contains(content, "docker") || strings.Contains(content, "containerd") || strings.Contains(content, "lxc") {
			return true
		}
	}

	if os.Getenv("container") != "" {
		return true
	}

	return false
}

type Category string

const (
	CategoryTerminal   Category = "terminal"
	CategoryShell     Category = "shell"
	CategoryEditor   Category = "editor"
	CategoryTools    Category = "tools"
	CategoryContainer Category = "container"
	CategoryFonts   Category = "fonts"
	CategoryTheme  Category = "theme"
)

func AllCategories() []Category {
	return []Category{
		CategoryTerminal,
		CategoryShell,
		CategoryEditor,
		CategoryTools,
		CategoryContainer,
		CategoryFonts,
	}
}

func (c Category) String() string { return string(c) }

func (c Category) IsValid() bool {
	return slices.Contains(AllCategories(), c)
}

type Distro string

const (
	DistroArch     Distro = "arch"
	DistroDebian   Distro = "debian"
	DistroFedora   Distro = "fedora"
	DistroSuse     Distro = "suse"
	DistroAlpine   Distro = "alpine"
	DistroBrew     Distro = "brew"
	DistroMacOS    Distro = "macos"
	DistroAll      Distro = "all"
	DistroFallback Distro = "fallback"
)

var DistroDetectionOrder = []Distro{DistroArch, DistroDebian, DistroFedora, DistroSuse, DistroAlpine, DistroBrew, DistroFallback}

var MacOSCommandOrder = []Distro{DistroMacOS, DistroBrew, DistroAll, DistroFallback}

func (d Distro) IsMacOS() bool { return d == DistroMacOS }

func commandOrder(distro Distro) []Distro {
	if distro.IsMacOS() {
		return MacOSCommandOrder
	}
	order := make([]Distro, 0, len(DistroDetectionOrder)+2)
	order = append(order, distro, DistroAll)
	for _, d := range DistroDetectionOrder {
		if d != distro && d != DistroBrew {
			order = append(order, d)
		}
	}
	return order
}

func resolveCommand(cmds map[Distro]string, distro Distro) string {
	for _, d := range commandOrder(distro) {
		if cmd := cmds[d]; cmd != "" {
			return cmd
		}
	}
	return ""
}

func releaseFileFor(distro Distro) string {
	switch distro {
	case DistroArch:
		return "/etc/arch-release"
	case DistroDebian:
		return "/etc/debian_version"
	case DistroFedora:
		return "/etc/fedora-release"
	case DistroSuse:
		return "/etc/SuSE-release"
	case DistroAlpine:
		return "/etc/alpine-release"
	}
	return ""
}

func DetectDistro() Distro {
	if runtime.GOOS == "darwin" {
		return DistroMacOS
	}

	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		lower := strings.ToLower(string(data))
		switch {
		case strings.Contains(lower, "arch"):
			return DistroArch
		case strings.Contains(lower, "debian") || strings.Contains(lower, "ubuntu"):
			return DistroDebian
		case strings.Contains(lower, "fedora") || strings.Contains(lower, "rhel") ||
			strings.Contains(lower, "rocky") || strings.Contains(lower, "almalinux"):
			return DistroFedora
		case strings.Contains(lower, "opensuse") || strings.Contains(lower, "sles"):
			return DistroSuse
		case strings.Contains(lower, "alpine"):
			return DistroAlpine
		}
	}

	cmds := []struct {
		distro Distro
		check  string
	}{
		{DistroArch, "pacman"},
		{DistroDebian, "apt-get"},
		{DistroFedora, "dnf"},
		{DistroSuse, "zypper"},
		{DistroAlpine, "apk"},
		{DistroBrew, "brew"},
	}
	for _, c := range cmds {
		if _, err := exec.LookPath(c.check); err == nil {
			if releaseFile := releaseFileFor(c.distro); releaseFile != "" {
				if _, err := os.Stat(releaseFile); err == nil {
					return c.distro
				}
				continue
			}
			return c.distro
		}
	}

	return ""
}

type Theme struct {
	Name        string      `json:"name" yaml:"name"`
	DisplayName string     `json:"display_name" yaml:"display_name"`
	Category   Category   `json:"category" yaml:"category"`
	Colors    ThemeColors `json:"colors" yaml:"colors"`
}

type ThemeColors struct {
	Background  string `json:"background" yaml:"background"`
	Foreground string `json:"foreground" yaml:"foreground"`
	Normal     ThemeColorSet `json:"normal" yaml:"normal"`
	Bright    ThemeColorSet `json:"bright" yaml:"bright"`
}

type ThemeColorSet struct {
	Black   string `json:"black" yaml:"black"`
	Red     string `json:"red" yaml:"red"`
	Green   string `json:"green" yaml:"green"`
	Yellow  string `json:"yellow" yaml:"yellow"`
	Blue    string `json:"blue" yaml:"blue"`
	Magenta string `json:"magenta" yaml:"magenta"`
	Cyan    string `json:"cyan" yaml:"cyan"`
	White   string `json:"white" yaml:"white"`
}

type Tool struct {
	Name        string   `json:"name" yaml:"name"`
	Category   Category `json:"category" yaml:"category"`
	Description string  `json:"description" yaml:"description"`

	Install   map[Distro]string `json:"install" yaml:"install"`
	Uninstall map[Distro]string `json:"uninstall" yaml:"uninstall"`
	Check     string           `json:"check" yaml:"check"`

	DependsOn []string `json:"depends_on" yaml:"depends_on"`

	SourceURL    string   `json:"source_url" yaml:"source_url"`
	Version     string   `json:"version" yaml:"version"`
	Alternatives []string `json:"alternatives" yaml:"alternatives"`

	Enabled bool `json:"enabled" yaml:"enabled"`
	Required bool `json:"required" yaml:"required"`
}

func (t *Tool) GetInstallCmd(distro Distro) string {
	return resolveCommand(t.Install, distro)
}

func (t *Tool) GetUninstallCmd(distro Distro) string {
	return resolveCommand(t.Uninstall, distro)
}

func FailedDependency(tool Tool, status map[string]bool) string {
	for _, dep := range tool.DependsOn {
		if ok, seen := status[dep]; seen && !ok {
			return dep
		}
	}
	return ""
}

func (t *Tool) HasInstallCommand() bool {
	for _, cmd := range t.Install {
		if cmd != "" {
			return true
		}
	}
	return false
}

func (t *Tool) HasUninstallCommand() bool {
	for _, cmd := range t.Uninstall {
		if cmd != "" {
			return true
		}
	}
	return false
}

func (t *Tool) IsBundle() bool {
	return !t.HasInstallCommand() && len(t.DependsOn) == 0
}

func (t *Tool) IsDockerIncompatible() bool {
	return DockerTools[t.Name]
}

func FilterDockerIncompatible(tools []Tool) []Tool {
	if !IsDocker() {
		return tools
	}
	var filtered []Tool
	for _, t := range tools {
		if !DockerTools[t.Name] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

type InstallResult struct {
	ToolName   string `json:"tool_name"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	DurationMs int64  `json:"duration_ms"`
	Distro    Distro `json:"distro,omitempty"`
}

type Installer interface {
	Install(tool *Tool) (*InstallResult, error)
	Uninstall(tool *Tool) (*InstallResult, error)
	IsInstalled(tool *Tool) (bool, error)
}

type Executor interface {
	Execute(cmd string) (string, error)
	ExecuteWithOutput(cmd string) (string, error)
}