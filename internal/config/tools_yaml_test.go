package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vitualizz/vitualizz-devstack/internal/config"
	"github.com/vitualizz/vitualizz-devstack/internal/domain/entities"
	"github.com/vitualizz/vitualizz-devstack/internal/infrastructure/installers"
)

const shippedToolsYAML = "../../cmd/vitualizz-devstack/config/tools.yaml"

var linuxOnly = regexp.MustCompile(`\b(sudo|apt-get|apt|dpkg|pacman|yay|dnf|zypper|apk|fc-cache)\b|_linux|Linux_|linux-`)

func loadShippedTools(t *testing.T) []entities.Tool {
	t.Helper()
	repo, err := config.NewToolRepository(shippedToolsYAML)
	if err != nil {
		t.Fatalf("loading %s: %v", shippedToolsYAML, err)
	}
	tools := repo.GetAll()
	if len(tools) == 0 {
		t.Fatal("shipped tools.yaml has no tools")
	}
	return tools
}

func TestShippedTools_MacOSCommandsArePortable(t *testing.T) {
	for _, tool := range loadShippedTools(t) {
		if !tool.HasInstallCommand() {
			continue
		}
		t.Run(tool.Name, func(t *testing.T) {
			install := tool.GetInstallCmd(entities.DistroMacOS)
			if install == "" {
				t.Fatal("no install command resolves on macOS")
			}
			if m := linuxOnly.FindString(install); m != "" {
				t.Errorf("macOS install command uses Linux-only %q:\n%s", m, install)
			}
			if uninstall := tool.GetUninstallCmd(entities.DistroMacOS); uninstall != "" {
				if m := linuxOnly.FindString(uninstall); m != "" {
					t.Errorf("macOS uninstall command uses Linux-only %q:\n%s", m, uninstall)
				}
			}
		})
	}
}

func TestShippedTools_LinuxNeverResolvesMacOSCommands(t *testing.T) {
	distros := []entities.Distro{
		entities.DistroArch, entities.DistroDebian, entities.DistroFedora,
		entities.DistroSuse, entities.DistroAlpine,
	}
	for _, tool := range loadShippedTools(t) {
		for _, d := range distros {
			for _, cmd := range []string{tool.GetInstallCmd(d), tool.GetUninstallCmd(d)} {
				if regexp.MustCompile(`--cask|xcode-select|\bbrew\b`).MatchString(cmd) {
					t.Errorf("%s on %s resolves to a macOS command:\n%s", tool.Name, d, cmd)
				}
			}
		}
	}
}

func TestShippedTools_ConfigIsComposedNotReplaced(t *testing.T) {
	configDir, err := filepath.Abs(filepath.Dir(shippedToolsYAML))
	if err != nil {
		t.Fatal(err)
	}
	repo, err := config.NewToolRepository(shippedToolsYAML)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		tool   string
		file   string
		marker string
	}{
		{"zsh-config", ".zshrc", "# vitualizz-devstack"},
		{"kitty-config", ".config/kitty/kitty.conf", "include vitualizz/kitty.conf"},
	}

	for _, tc := range cases {
		for _, symlinked := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/symlink=%v", tc.tool, symlinked), func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("HOME", home)

				target := filepath.Join(home, tc.file)
				if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
					t.Fatal(err)
				}
				original := "user setting 1\nuser setting 2\n"
				realFile := target
				if symlinked {
					realFile = filepath.Join(home, "dotfiles", filepath.Base(tc.file))
					if err := os.MkdirAll(filepath.Dir(realFile), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(realFile, target); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.WriteFile(realFile, []byte(original), 0o644); err != nil {
					t.Fatal(err)
				}

				tool := repo.GetByID(tc.tool)
				inst := installers.NewToolInstaller()
				inst.SetConfigDir(configDir)

				for i := 0; i < 2; i++ {
					if res, err := inst.Install(tool); err != nil {
						t.Fatalf("install #%d: %v (%s)", i+1, err, res.Message)
					}
				}
				if ok, _ := inst.IsInstalled(tool); !ok {
					t.Error("check fails after install")
				}

				got := readFile(t, target)
				if !strings.Contains(got, original) {
					t.Errorf("user config lost after install:\n%s", got)
				}
				if n := strings.Count(got, tc.marker); n != 1 {
					t.Errorf("marker appears %d times after installing twice, want 1:\n%s", n, got)
				}

				if res, err := inst.Uninstall(tool); err != nil {
					t.Fatalf("uninstall: %v (%s)", err, res.Message)
				}
				if got := readFile(t, target); got != original {
					t.Errorf("uninstall did not restore the original file:\ngot:  %q\nwant: %q", got, original)
				}
				if fi, err := os.Lstat(target); err != nil || (fi.Mode()&os.ModeSymlink != 0) != symlinked {
					t.Errorf("symlink state changed (symlinked=%v)", symlinked)
				}
			})
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
