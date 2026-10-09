package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractEmbeddedConfig(t *testing.T) {
	dir, err := extractEmbeddedConfig()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, f := range []string{"tools.yaml", "zsh/zshrc", "kitty/kitty.conf", "kitty/tokyonight_night.conf", "nvim/init.lua", "git/gitconfig", "nvim/lua/vitualizz/health.lua"} {
		info, err := os.Stat(filepath.Join(dir, f))
		if err != nil {
			t.Errorf("%s not extracted: %v", f, err)
			continue
		}
		if info.Mode().Perm()&0o200 == 0 {
			t.Errorf("%s is read-only (%v); seeded configs must stay editable", f, info.Mode().Perm())
		}
	}
}
