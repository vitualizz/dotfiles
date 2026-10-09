package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vitualizz/dotfiles/internal/config"
	"github.com/vitualizz/dotfiles/internal/domain/entities"
	"github.com/vitualizz/dotfiles/internal/infrastructure/installers"
)

const shippedToolsYAML = "../../cmd/dotfiles/config/tools.yaml"

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

func TestShippedTools_SeedsConfigWithoutOverwriting(t *testing.T) {
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
		append bool
	}{
		{"zsh-config", ".zshrc", true},
		{"kitty-config", ".config/kitty/kitty.conf", false},
		{"starship-config", ".config/starship.toml", false},
		{"nvim-config", ".config/nvim/init.lua", false},
		{"git-config", ".gitconfig", false},
	}

	for _, tc := range cases {
		for _, existing := range []string{"none", "file", "symlink"} {
			t.Run(fmt.Sprintf("%s/%s", tc.tool, existing), func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("HOME", home)
				target := filepath.Join(home, tc.file)
				original := "user setting 1\nuser setting 2\n"

				switch existing {
				case "file":
					writeFile(t, target, original)
				case "symlink":
					real := filepath.Join(home, "dotfiles", filepath.Base(tc.file))
					writeFile(t, real, original)
					if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(real, target); err != nil {
						t.Fatal(err)
					}
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
				switch {
				case existing == "none" && got == "":
					t.Error("config was not seeded")
				case existing != "none" && !tc.append && got != original:
					t.Errorf("existing config was modified:\n%s", got)
				case existing != "none" && tc.append && !strings.HasPrefix(got, original):
					t.Errorf("existing config lost:\n%s", got)
				}
				if tc.append {
					if n := strings.Count(got, "# ── vitualizz dotfiles"); n != 1 {
						t.Errorf("block appears %d times after installing twice, want 1", n)
					}
				}
				if fi, err := os.Lstat(target); err != nil || (fi.Mode()&os.ModeSymlink != 0) != (existing == "symlink") {
					t.Error("symlink state changed")
				}

				if _, err := inst.Uninstall(tool); err != nil {
					t.Fatalf("uninstall: %v", err)
				}
				if after := readFile(t, target); after != got {
					t.Error("uninstall modified a config that belongs to the user")
				}
			})
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
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
