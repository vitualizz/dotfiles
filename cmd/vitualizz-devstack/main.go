package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vitualizz/vitualizz-devstack/internal/config"
	"github.com/vitualizz/vitualizz-devstack/internal/domain/entities"
	"github.com/vitualizz/vitualizz-devstack/internal/infrastructure/installers"
	"github.com/vitualizz/vitualizz-devstack/internal/infrastructure/logger"
	"github.com/vitualizz/vitualizz-devstack/internal/ui/components"
	"github.com/vitualizz/vitualizz-devstack/i18n/locales"
)

//go:embed all:config
var embeddedConfig embed.FS

func main() {
	ciMode := false
	selfDestruct := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--ci":
			ciMode = true
		case "--self-destruct":
			selfDestruct = true
		}
	}

	binPath := ""
	if selfDestruct {
		var err error
		binPath, err = os.Executable()
		if err != nil {
			binPath = ""
		}
		defer func() {
			if binPath != "" {
				os.Remove(binPath)
			}
		}()
	}

	if err := preflight(entities.DetectDistro()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if binPath != "" {
			_ = os.Remove(binPath)
		}
		os.Exit(1)
	}

	var configPath string
	var configDir string

	if envPath := os.Getenv("DEVSTACK_CONFIG"); envPath != "" {
		configPath = envPath
		configDir = filepath.Dir(envPath)
	} else {
		var err error
		configDir, err = extractEmbeddedConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting config: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(configDir)
		configPath = filepath.Join(configDir, "tools.yaml")
	}

	repo, err := config.NewToolRepository(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.NewInstallLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create log file: %v\n", err)
	} else {
		defer log.Close()
	}

	installer := installers.NewToolInstaller()
	installer.SetConfigDir(configDir)
	installer.SetLogger(log)

	i18n := locales.NewI18nSimple()
	inDocker := entities.IsDocker()

	if ciMode {
		runCI(repo, installer, inDocker, log)
		return
	}

	installer.SetDetached(true)
	runTUI(repo, installer, i18n, log)
}

func preflight(distro entities.Distro) error {
	if os.Geteuid() == 0 && (distro.IsMacOS() || os.Getenv("SUDO_USER") != "") {
		return errors.New("do not run with sudo: config would be installed for root instead of your user. " +
			"Re-run without sudo; DevStack asks for your password when a tool needs it")
	}
	if !distro.IsMacOS() {
		return nil
	}
	if _, err := exec.LookPath("brew"); err != nil {
		for _, p := range []string{"/opt/homebrew/bin/brew", "/usr/local/bin/brew"} {
			if _, statErr := os.Stat(p); statErr == nil {
				return nil
			}
		}
		return errors.New("homebrew is required on macOS. Install it first:\n" +
			`  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	}
	return nil
}

func extractEmbeddedConfig() (string, error) {
	tmpDir, err := os.MkdirTemp("", "dotfiles-config-*")
	if err != nil {
		return "", err
	}
	root, err := fs.Sub(embeddedConfig, "config")
	if err != nil {
		return "", err
	}
	return tmpDir, os.CopyFS(tmpDir, root)
}

func runCI(repo *config.ToolRepository, installer *installers.ToolInstaller, inDocker bool, log *logger.InstallLogger) {
	fmt.Println("┌─────────────────────────────────────────────┐")
	fmt.Println("│  Vitualizz DevStack — CI Mode               │")
	fmt.Println("└─────────────────────────────────────────────┘")
	fmt.Println()

	if inDocker {
		fmt.Println("  🐳 Docker environment detected")
		fmt.Println("  Skipping display tools (kitty)")
		fmt.Println()
	}

	if log != nil {
		fmt.Printf("  📝 Log: %s\n\n", log.LogPath())
	}

	allTools := repo.GetAll()
	if inDocker {
		allTools = entities.FilterDockerIncompatible(allTools)
	}

	total := len(allTools)
	fmt.Printf("  Installing %d tools...\n\n", total)

	var success, failed int
	status := make(map[string]bool, total)
	for _, tool := range allTools {
		if !tool.HasInstallCommand() {
			fmt.Printf("  ⊘ %s (bundle)\n", tool.Name)
			continue
		}

		installed, _ := installer.IsInstalled(&tool)
		if installed {
			success++
			status[tool.Name] = true
			fmt.Printf("  ✓ %s (already installed)\n", tool.Name)
			continue
		}

		if dep := entities.FailedDependency(tool, status); dep != "" {
			failed++
			status[tool.Name] = false
			fmt.Printf("  ⊘ %s (skipped: dependency %s failed)\n", tool.Name, dep)
			continue
		}

		fmt.Printf("  → %s... ", tool.Name)
		result, err := installer.Install(&tool)
		ok := err == nil && result.Success
		status[tool.Name] = ok
		if !ok {
			failed++
			fmt.Println("✗")
			if result != nil && result.Message != "" {
				fmt.Printf("    └─ %s\n", truncate(result.Message, 80))
			}
		} else {
			success++
			fmt.Println("✓")
		}
	}

	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────┐")
	fmt.Printf("│  Total: %d  |  ✓ %d  |  ✗ %d           \n", total, success, failed)
	fmt.Println("└─────────────────────────────────────────────┘")

	if log != nil {
		fmt.Printf("\n  📝 Full log: %s\n", log.LogPath())
	}

	if failed > 0 {
		fmt.Println()
		fmt.Println("Failed tools:")
		for _, tool := range allTools {
			if !tool.HasInstallCommand() {
				continue
			}
			installed, _ := installer.IsInstalled(&tool)
			if !installed {
				fmt.Printf("  ✗ %s\n", tool.Name)
			}
		}
		os.Exit(1)
	}
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func runTUI(repo *config.ToolRepository, installer *installers.ToolInstaller, i18n *locales.I18nSimple, log *logger.InstallLogger) {
	logPath := ""
	if log != nil {
		logPath = log.LogPath()
	}

	app := components.NewApp(repo, installer, installer, i18n, logPath)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
