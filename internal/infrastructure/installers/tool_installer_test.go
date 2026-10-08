package installers

import (
	"os"
	"testing"

	"github.com/vitualizz/dotfiles/internal/domain/entities"
)

func TestToolInstaller(t *testing.T) {
	installer := NewToolInstaller()

	t.Run("creates installer", func(t *testing.T) {
		if installer == nil {
			t.Error("expected installer to not be nil")
		}
	})

	t.Run("detects distro", func(t *testing.T) {
		distro := installer.Distro()
		if distro == "" {
			t.Log("distro detection returned empty (expected in test environment)")
		}
	})
}

func BenchmarkToolInstaller(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewToolInstaller()
	}
}
func TestNeedsSudo(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root never needs sudo")
	}
	tool := func(install, uninstall, check string) entities.Tool {
		return entities.Tool{
			Name:      "t",
			Install:   map[entities.Distro]string{entities.DistroAll: install},
			Uninstall: map[entities.Distro]string{entities.DistroAll: uninstall},
			Check:     check,
		}
	}

	tests := []struct {
		name      string
		tool      entities.Tool
		uninstall bool
		want      bool
	}{
		{"missing tool with sudo install", tool("sudo apt-get install -y t", "", "false"), false, true},
		{"missing tool with cask install", tool("brew install --cask t", "", "false"), false, true},
		{"installed tool is skipped", tool("sudo apt-get install -y t", "", "true"), false, false},
		{"no sudo in command", tool("brew install t", "", "false"), false, false},
		{"uninstall of installed tool with sudo", tool("", "sudo apt-get remove -y t", "true"), true, true},
		{"uninstall of missing tool is skipped", tool("", "sudo apt-get remove -y t", "false"), true, false},
	}

	installer := NewToolInstaller()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := installer.NeedsSudo([]entities.Tool{tt.tool}, tt.uninstall); got != tt.want {
				t.Errorf("NeedsSudo() = %v, want %v", got, tt.want)
			}
		})
	}
}
