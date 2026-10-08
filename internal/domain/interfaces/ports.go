package interfaces

import "github.com/vitualizz/dotfiles/internal/domain/entities"

type ToolRepository interface {
	GetAll() []entities.Tool
	GetMainTools() []entities.Tool
	GetByCategory(category entities.Category) []entities.Tool
	GetByID(name string) *entities.Tool
	GetDependencies(name string) []entities.Tool
	GetDependents(name string) []entities.Tool
	GetPackages() []Package
	GetThemes() []entities.Theme
	GetThemeByName(name string) *entities.Theme
	Save(tool entities.Tool)
}

type Package struct {
	Name             string
	Label            string
	Icon             string
	Description      string
	Tools           []entities.Tool
	DefaultSelected bool
	Selected        bool
}

type SudoPort interface {
	NeedsSudo(tools []entities.Tool, uninstall bool) bool
	SudoReady(ok bool)
}

type InstallerPort interface {
	Install(tool *entities.Tool) (*entities.InstallResult, error)
	Uninstall(tool *entities.Tool) (*entities.InstallResult, error)
	IsInstalled(tool *entities.Tool) (bool, error)
}

type ExecutorPort interface {
	Execute(cmd string) (string, error)
	ExecuteWithOutput(cmd string) (string, error)
}

type I18nPort interface {
	Get(key string, lang string) string
	SetLanguage(lang string)
	GetAvailableLanguages() []string
}