package ui

import (
	"fmt"
	"strings"

	"memex/config"
	"memex/search"
	"memex/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Config      config.Config
	Storage     *storage.Storage
	Input       string
	Cursor      int
	Suggestions []search.Match
	Selected    int
	Width       int
	Height      int
	Quitting    bool
	OutputCmd   string
}

func NewModel(cfg config.Config, store *storage.Storage, initialInput string) Model {
	return Model{
		Config:  cfg,
		Storage: store,
		Input:   initialInput,
		Cursor:  len(initialInput),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if len(m.Suggestions) > 0 {
				m.OutputCmd = m.Suggestions[m.Selected].Command.Cmd
				// Update frequency
				m.Storage.AddCommand(m.OutputCmd, nil)
				m.Storage.Save()
			} else if m.Input != "" {
				// Allow running what's typed even if not in list
				m.OutputCmd = m.Input
				m.Storage.AddCommand(m.OutputCmd, nil)
				m.Storage.Save()
			}
			return m, tea.Quit

		case tea.KeyTab:
			if len(m.Suggestions) > 0 {
				// Autocomplete with selected
				m.Input = m.Suggestions[m.Selected].Command.Cmd
				m.Cursor = len(m.Input)
				// Reset selection to top after autocomplete? Or keep it?
				// Let's keep it simple: update input, re-search will happen
			}

		case tea.KeyUp:
			if m.Selected > 0 {
				m.Selected--
			}

		case tea.KeyDown:
			if m.Selected < len(m.Suggestions)-1 && m.Selected < m.Config.SuggestionLimit-1 {
				m.Selected++
			}

		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.Input) > 0 && m.Cursor > 0 {
				m.Input = m.Input[:m.Cursor-1] + m.Input[m.Cursor:]
				m.Cursor--
				m.Selected = 0 // Reset selection
			}

		case tea.KeyLeft:
			if m.Cursor > 0 {
				m.Cursor--
			}

		case tea.KeyRight:
			if m.Cursor < len(m.Input) {
				m.Cursor++
			}

		case tea.KeySpace:
			m.Input = m.Input[:m.Cursor] + " " + m.Input[m.Cursor:]
			m.Cursor++
			m.Selected = 0

		case tea.KeyRunes:
			m.Input = m.Input[:m.Cursor] + string(msg.Runes) + m.Input[m.Cursor:]
			m.Cursor += len(msg.Runes)
			m.Selected = 0 // Reset selection
		}

		// Update suggestions on any key press that might change state
		m.Suggestions = search.Search(m.Input, m.Storage.GetCommands())

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var s strings.Builder

	// 1. Input Line
	// Construct ghost text
	ghost := ""
	if m.Config.GhostText && len(m.Suggestions) > 0 && m.Input != "" {
		topCmd := m.Suggestions[m.Selected].Command.Cmd
		if strings.HasPrefix(topCmd, m.Input) {
			ghost = topCmd[len(m.Input):]
		}
	}

	// Render Input
	s.WriteString(InputStyle.Render("> " + m.Input))
	s.WriteString(GhostStyle.Render(ghost))
	s.WriteString("\n")

	// 2. Suggestions
	limit := m.Config.SuggestionLimit
	if limit > len(m.Suggestions) {
		limit = len(m.Suggestions)
	}

	for i := 0; i < limit; i++ {
		match := m.Suggestions[i]
		line := fmt.Sprintf("  %s", match.Command.Cmd)

		if i == m.Selected {
			s.WriteString(SelectedItemStyle.Render(line))
		} else {
			s.WriteString(ItemStyle.Render(line))
		}

		// Render tags if any
		if len(match.Command.Tags) > 0 {
			tags := fmt.Sprintf(" [%s]", strings.Join(match.Command.Tags, ", "))
			s.WriteString(TagStyle.Render(tags))
		}
		s.WriteString("\n")
	}

	return s.String()
}
