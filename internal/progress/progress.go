package progress

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
)

const (
	padding  = 16
	maxWidth = 160
)

// Model to represent the state of the progress bar and stopwatch.
type model struct {
	totalItems  int
	currentItem int
	stopwatch   stopwatch.Model
	progressbar progress.Model
}

// initialize the stopwatch
func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

// Updates the progress bar and stopwatch.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Resize the progress bar dynamically
		m.progressbar.PercentageStyle.PaddingLeft(padding)
		m.progressbar.Width = msg.Width - padding - 4
		if m.progressbar.Width > maxWidth {
			m.progressbar.Width = maxWidth
		}
		return m, nil
	case int:
		// Update progress
		m.currentItem += msg

		if m.currentItem >= m.totalItems {
			return m, tea.Quit
		}
		return m, nil
	case progress.FrameMsg:
		progressModel, _ := m.progressbar.Update(msg)
		m.progressbar = progressModel.(progress.Model)
		return m, nil
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

// Returns a string representation of the progress bar and stopwatch.
func (m model) View() string {
	if m.currentItem >= m.totalItems {
		// Print summary
		m.stopwatch.Stop()
		return fmt.Sprintf("Done after %s\n", m.stopwatch.View())
	}

	// Calculate progress percentage
	percent := float64(m.currentItem) / float64(m.totalItems)
	out := fmt.Sprintf("\n %s Time Elapsed: %s\n", m.progressbar.ViewAs(percent), m.stopwatch.View())
	return out
}

// Initialize a new progress bar model.
func InitBar(num_items int) *tea.Program {
	progressbar := progress.New(progress.WithDefaultGradient())

	p := tea.NewProgram(model{
		totalItems:  num_items,
		stopwatch:   stopwatch.NewWithInterval(time.Second),
		currentItem: 0,
		progressbar: progressbar,
	})

	return p
}
