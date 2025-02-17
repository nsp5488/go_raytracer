package progress

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
)

type model struct {
	totalItems  int
	currentItem int
	stopwatch   stopwatch.Model
	progressbar progress.Model
}

func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.stopwatch.Running() {
		m.stopwatch.Start()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case int:
		// Update progress
		m.currentItem = msg
		if m.currentItem >= m.totalItems {
			return m, tea.Quit
		}
		return m, nil
	}
	var cmd tea.Cmd

	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}
func (m model) View() string {
	if m.currentItem >= m.totalItems {
		return "Done!\n"
	}

	// Calculate progress percentage
	percent := float64(m.currentItem) / float64(m.totalItems)

	return fmt.Sprintf("\n %s Time Elapsed: %s\n",
		m.progressbar.ViewAs(percent), m.stopwatch.View(),
	)
}

func InitBar(num_items int) *tea.Program {
	progressbar := progress.New(progress.WithDefaultGradient())

	p := tea.NewProgram(model{
		totalItems:  num_items,
		stopwatch:   stopwatch.NewWithInterval(time.Millisecond),
		currentItem: 0,
		progressbar: progressbar,
	})

	return p
}
