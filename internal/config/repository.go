package config

import (
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	"github.com/vitualizz/vitualizz-devstack/internal/domain/entities"
	"github.com/vitualizz/vitualizz-devstack/internal/domain/interfaces"
)

type ToolConfig struct {
	Name         string            `yaml:"name"`
	Category    string           `yaml:"category"`
	Description string           `yaml:"description"`
	DependsOn   []string         `yaml:"depends_on"`
	Install     map[string]string `yaml:"install"`
	Uninstall   map[string]string `yaml:"uninstall"`
	Check       string           `yaml:"check"`
	SourceURL   string           `yaml:"source_url"`
	Version     string           `yaml:"version"`
	Alternatives []string        `yaml:"alternatives"`
	Enabled     bool             `yaml:"enabled"`
	Required    bool             `yaml:"required"`
}

type ThemeColorsConfig struct {
	Background  string `yaml:"background"`
	Foreground string `yaml:"foreground"`
	Normal     struct {
		Black   string `yaml:"black"`
		Red     string `yaml:"red"`
		Green   string `yaml:"green"`
		Yellow  string `yaml:"yellow"`
		Blue    string `yaml:"blue"`
		Magenta string `yaml:"magenta"`
		Cyan    string `yaml:"cyan"`
		White   string `yaml:"white"`
	} `yaml:"normal"`
	Bright struct {
		Black   string `yaml:"black"`
		Red     string `yaml:"red"`
		Green   string `yaml:"green"`
		Yellow  string `yaml:"yellow"`
		Blue    string `yaml:"blue"`
		Magenta string `yaml:"magenta"`
		Cyan    string `yaml:"cyan"`
		White   string `yaml:"white"`
	} `yaml:"bright"`
}

type ThemeConfig struct {
	Name        string            `yaml:"name"`
	DisplayName string            `yaml:"display_name"`
	Category   string            `yaml:"category"`
	Colors     ThemeColorsConfig `yaml:"colors"`
}

type ConfigFile struct {
	Tools  []ToolConfig  `yaml:"tools"`
	Themes []ThemeConfig `yaml:"themes"`
}

type ToolRepository struct {
	tools  []entities.Tool
	themes []entities.Theme
}

func NewToolRepository(path string) (*ToolRepository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	tools := make([]entities.Tool, 0, len(cfg.Tools))
	for _, t := range cfg.Tools {
		tools = append(tools, configToTool(t))
	}

	themes := make([]entities.Theme, 0, len(cfg.Themes))
	for _, t := range cfg.Themes {
		themes = append(themes, configToTheme(t))
	}

	return &ToolRepository{tools: tools, themes: themes}, nil
}

func configToTool(cfg ToolConfig) entities.Tool {
	install := make(map[entities.Distro]string, len(cfg.Install))
	for distro, cmd := range cfg.Install {
		install[entities.Distro(distro)] = cmd
	}

	uninstall := make(map[entities.Distro]string, len(cfg.Uninstall))
	for distro, cmd := range cfg.Uninstall {
		uninstall[entities.Distro(distro)] = cmd
	}

	var dependsOn []string
	for _, d := range cfg.DependsOn {
		if d != "" {
			dependsOn = append(dependsOn, d)
		}
	}

	var alt []string
	for _, a := range cfg.Alternatives {
		if a != "" {
			alt = append(alt, a)
		}
	}

	return entities.Tool{
		Name:         cfg.Name,
		Category:    entities.Category(cfg.Category),
		Description:  cfg.Description,
		Install:     install,
		Uninstall:   uninstall,
		Check:       cfg.Check,
		DependsOn:   dependsOn,
		SourceURL:   cfg.SourceURL,
		Version:     cfg.Version,
		Alternatives: alt,
		Enabled:    cfg.Enabled,
		Required:   cfg.Required,
	}
}

func configToTheme(cfg ThemeConfig) entities.Theme {
	return entities.Theme{
		Name:        cfg.Name,
		DisplayName: cfg.DisplayName,
		Category:   entities.Category(cfg.Category),
		Colors: entities.ThemeColors{
			Background:  cfg.Colors.Background,
			Foreground: cfg.Colors.Foreground,
			Normal: entities.ThemeColorSet{
				Black:   cfg.Colors.Normal.Black,
				Red:     cfg.Colors.Normal.Red,
				Green:   cfg.Colors.Normal.Green,
				Yellow:  cfg.Colors.Normal.Yellow,
				Blue:    cfg.Colors.Normal.Blue,
				Magenta: cfg.Colors.Normal.Magenta,
				Cyan:    cfg.Colors.Normal.Cyan,
				White:   cfg.Colors.Normal.White,
			},
			Bright: entities.ThemeColorSet{
				Black:   cfg.Colors.Bright.Black,
				Red:     cfg.Colors.Bright.Red,
				Green:   cfg.Colors.Bright.Green,
				Yellow:  cfg.Colors.Bright.Yellow,
				Blue:    cfg.Colors.Bright.Blue,
				Magenta: cfg.Colors.Bright.Magenta,
				Cyan:    cfg.Colors.Bright.Cyan,
				White:   cfg.Colors.Bright.White,
			},
		},
	}
}

func (r *ToolRepository) GetAll() []entities.Tool {
	return r.tools
}

func (r *ToolRepository) GetMainTools() []entities.Tool {
	depended := make(map[string]bool)
	for _, t := range r.tools {
		for _, dep := range t.DependsOn {
			depended[dep] = true
		}
	}

	var main []entities.Tool
	for _, t := range r.tools {
		if !depended[t.Name] {
			main = append(main, t)
		}
	}
	return main
}

func (r *ToolRepository) GetDependents(name string) []entities.Tool {
	var result []entities.Tool
	for _, t := range r.tools {
		for _, dep := range t.DependsOn {
			if dep == name {
				result = append(result, t)
				break
			}
		}
	}
	return result
}

func (r *ToolRepository) GetByCategory(category entities.Category) []entities.Tool {
	var result []entities.Tool
	for _, t := range r.tools {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

func (r *ToolRepository) GetByID(name string) *entities.Tool {
	for _, t := range r.tools {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (r *ToolRepository) GetDependencies(name string) []entities.Tool {
	tool := r.GetByID(name)
	if tool == nil {
		return nil
	}

	var deps []entities.Tool
	for _, depName := range tool.DependsOn {
		if dep := r.GetByID(depName); dep != nil {
			deps = append(deps, *dep)
		}
	}
	return deps
}

func (r *ToolRepository) GetAllDependencies(name string) []entities.Tool {
	var visited []string
	var result []entities.Tool

	var collect func(toolName string)
	collect = func(toolName string) {
		if slicesContains(visited, toolName) {
			return
		}
		visited = append(visited, toolName)

		deps := r.GetDependencies(toolName)
		for _, dep := range deps {
			if !slicesContains(visited, dep.Name) {
				result = append(result, dep)
				collect(dep.Name)
			}
		}
	}

	collect(name)
	return result
}

func (r *ToolRepository) Save(tool entities.Tool) {
	for i, t := range r.tools {
		if t.Name == tool.Name {
			r.tools[i] = tool
			return
		}
	}
	r.tools = append(r.tools, tool)
}

func (r *ToolRepository) GetThemes() []entities.Theme {
	return r.themes
}

func (r *ToolRepository) GetThemeByName(name string) *entities.Theme {
	for _, t := range r.themes {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (r *ToolRepository) GetPackages() []interfaces.Package {
	order := []entities.Category{
		entities.CategoryTerminal,
		entities.CategoryShell,
		entities.CategoryEditor,
		entities.CategoryTools,
		entities.CategoryContainer,
		entities.CategoryFonts,
	}

	icons := map[entities.Category]string{
		entities.CategoryTerminal:   "🖥️",
		entities.CategoryShell:       "🐚",
		entities.CategoryEditor:   "📝",
		entities.CategoryTools:     "🛠️",
		entities.CategoryContainer:  "📦",
		entities.CategoryFonts:      "🔤",
	}

	descriptions := map[entities.Category]string{
		entities.CategoryTerminal:   "Kitty + Tokyo Night",
		entities.CategoryShell:       "Zsh + Starship + autosuggestions + syntax-highlighting",
		entities.CategoryEditor:   "Neovim",
		entities.CategoryTools:     "fzf · bat · eza · yazi · jq · lazygit · gh · uv",
		entities.CategoryContainer:  "Docker + Docker Compose",
		entities.CategoryFonts:      "Hack Nerd Font",
	}

	byCategory := make(map[entities.Category][]entities.Tool)
	for _, t := range r.tools {
		if !t.IsBundle() {
			byCategory[t.Category] = append(byCategory[t.Category], t)
		}
	}

	var pkgs []interfaces.Package
	for _, cat := range order {
		tools, ok := byCategory[cat]
		if !ok || len(tools) == 0 {
			continue
		}

		if cat == entities.CategoryFonts {
			pkgs = append(pkgs, interfaces.Package{
				Name:             string(cat),
				Label:            cases.Title(language.Und).String(string(cat)),
				Icon:             icons[cat],
				Description:      descriptions[cat],
				Tools:            tools,
				DefaultSelected:   true,
				Selected:        true,
			})
			continue
		}

		pkgs = append(pkgs, interfaces.Package{
			Name:             string(cat),
			Label:            cases.Title(language.Und).String(string(cat)),
			Icon:             icons[cat],
			Description:      descriptions[cat],
			Tools:            tools,
			DefaultSelected:   false,
			Selected:        false,
		})
	}

	return pkgs
}

func slicesContains[T comparable](s []T, v T) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}