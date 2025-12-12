package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/avirooppal/gosysutil/monitor"
)

type tickMsg time.Time

type model struct {
	lastStats    *monitor.SystemStats
	currentStats *monitor.SystemStats
	err          error
	width        int
	height       int
}

func initialModel() model {
	// Fetch initial stats immediately
	stats, err := monitor.GetSystemStats()
	return model{
		currentStats: stats, // lastStats will be nil initially
		err:          err,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		newStats, err := monitor.GetSystemStats()
		m.err = err
		if err == nil {
			m.lastStats = m.currentStats
			m.currentStats = newStats
		}
		return m, tickCmd()
	}

	return m, nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
