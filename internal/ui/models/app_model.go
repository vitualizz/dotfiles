package models

import "github.com/vitualizz/dotfiles/internal/domain/entities"

type ViewState int

const (
	StateLanguageSelect ViewState = iota
	StateMainMenu
	StateProgress
	StateThanks
	StateSettings
	StateAbout
)

type AppModel struct {
	ViewState    ViewState
	CurrentLang  string
	IsLoading    bool
	MainMenuChoice int
	SettingsChoice int
	ToolChoice   int
	InDocker     bool

	ProgressTools   []entities.Tool
	ProgressIdx     int
	ProgressResults []entities.InstallResult
	ProgressStatus  map[string]bool
	ProgressLastOutput string
	IsUninstallMode bool
	ShowLog        bool
	Results        []entities.InstallResult
	LogPath        string
}

func NewAppModel() *AppModel {
	return &AppModel{
		ViewState:      StateLanguageSelect,
		CurrentLang:    "es",
		MainMenuChoice: 0,
		SettingsChoice: 0,
		ToolChoice:     0,
		ProgressStatus: make(map[string]bool),
	}
}

func (m *AppModel) StartProgress(tools []entities.Tool) {
	m.ViewState = StateProgress
	m.ProgressTools = tools
	m.ProgressIdx = 0
	m.ProgressResults = nil
	m.ProgressStatus = make(map[string]bool)
	m.IsUninstallMode = false
}

func (m *AppModel) StartUninstallProgress(tools []entities.Tool) {
	m.ViewState = StateProgress
	m.ProgressTools = tools
	m.ProgressIdx = 0
	m.ProgressResults = nil
	m.ProgressStatus = make(map[string]bool)
	m.IsUninstallMode = true
}

func (m *AppModel) UpdateProgress(tool entities.Tool, success bool, msg string) {
	m.ProgressStatus[tool.Name] = success
	m.ProgressResults = append(m.ProgressResults, entities.InstallResult{
		ToolName: tool.Name,
		Success:  success,
		Message:  msg,
	})
	m.ProgressLastOutput = msg
	m.ProgressIdx++
}

func (m *AppModel) GetCurrentProgressTool() entities.Tool {
	if m.ProgressIdx < len(m.ProgressTools) {
		return m.ProgressTools[m.ProgressIdx]
	}
	return entities.Tool{}
}

func (m *AppModel) IsProgressDone() bool {
	return m.ProgressIdx >= len(m.ProgressTools)
}

func (m *AppModel) GetProgressStats() (total, success, failed int) {
	total = len(m.ProgressTools)
	for _, r := range m.ProgressResults {
		if r.Success {
			success++
		} else {
			failed++
		}
	}
	return total, success, failed
}
