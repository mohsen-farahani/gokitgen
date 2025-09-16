package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mohsen-farahani/gokitgen/pkg/generator/model"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5F87")).
				Bold(true)
)

type item struct {
	title, desc string
	command     string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type mainModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.choice = i.command
			}
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	if m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {
	items := []list.Item{
		item{title: "🚀 Init Project", desc: "Initialize a new Go Kit project structure", command: "init"},
		item{title: "📦 Generate Model", desc: "Generate model, service, API, tests, and more", command: "model"},
		item{title: "⚡ Generate Command", desc: "Coming soon...", command: "command"},
	}

	const defaultWidth = 50
	const defaultHeight = 15

	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, defaultHeight)
	l.Title = "✨ Go Kit Generator"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = selectedItemStyle
	l.Styles.HelpStyle = itemStyle

	m := mainModel{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	mm := finalModel.(mainModel)
	if mm.choice == "" {
		return
	}

	selectedItem := mm.list.SelectedItem().(item)
	fmt.Printf("\n✅ You selected: %s\n\n", selectedItem.title)

	switch mm.choice {
	case "init":
		fmt.Println("Initializing project... (not implemented yet)")
	case "model":
		config := model.RunWizard()
		if err := model.GenerateCode(config); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Println("✅ Model generated successfully!")
		}
	case "command":
		fmt.Println("Command generator coming soon...")
	}
}