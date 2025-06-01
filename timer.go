package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	workDuration  = 30 * time.Minute
	breakDuration = 10 * time.Minute
)

type model struct {
	choice     int
	session    sessionType
	progress   progress.Model
	duration   time.Duration
	elapsed    time.Duration
	isRunning  bool
	isComplete bool
	showMenu   bool
}

type sessionType int

const (
	none sessionType = iota
	work
	breakSession
)

func (s sessionType) String() string {
	switch s {
	case work:
		return "Work"
	case breakSession:
		return "Break"
	default:
		return ""
	}
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// clearScreen clears the terminal screen
func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func initialModel() model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return model{
		progress: p,
		showMenu: true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showMenu {
			switch msg.String() {
			case "1", "w", "work":
				m.session = work
				m.duration = workDuration
				m.elapsed = 0
				m.isRunning = true
				m.isComplete = false
				m.showMenu = false
				return m, tickCmd()
			case "2", "b", "break":
				m.session = breakSession
				m.duration = breakDuration
				m.elapsed = 0
				m.isRunning = true
				m.isComplete = false
				m.showMenu = false
				return m, tickCmd()
			case "q", "quit", "ctrl+c":
				clearScreen()
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "q", "quit", "ctrl+c":
				clearScreen()
				return m, tea.Quit
			case "r", "return":
				if m.isComplete {
					clearScreen() // Clear screen when returning to menu
					m.showMenu = true
					m.isRunning = false
					m.isComplete = false
					m.elapsed = 0
				}
				return m, nil
			case " ":
				if m.isRunning {
					m.isRunning = false
				} else if !m.isComplete {
					m.isRunning = true
					return m, tickCmd()
				}
				return m, nil
			}
		}

	case tickMsg:
		if m.isRunning && !m.isComplete {
			m.elapsed += time.Second
			if m.elapsed >= m.duration {
				m.isComplete = true
				m.isRunning = false
				clearScreen() // Clear screen when session completes
			} else {
				return m, tickCmd()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 2)

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Padding(0, 2)

	sessionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	completeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	if m.showMenu {
		s += titleStyle.Render("Timer")
		s += "\n\n"
		s += menuStyle.Render("Choose an option:")
		s += "\n"
		s += menuStyle.Render("  1 / w  →  Work session (30 minutes)")
		s += "\n"
		s += menuStyle.Render("  2 / b  →  Break session (10 minutes)")
		s += "\n"
		s += menuStyle.Render("  q      →  Quit")
		s += "\n\n"
	} else {

		s += titleStyle.Render(fmt.Sprintf("%s Session", m.session))
		s += "\n\n"

		remaining := m.duration - m.elapsed
		if remaining < 0 {
			remaining = 0
		}
		mins := int(remaining.Minutes())
		secs := int(remaining.Seconds()) % 60

		s += sessionStyle.Render(fmt.Sprintf("%s Session", m.session))
		s += " - "
		s += timeStyle.Render(fmt.Sprintf("Time Left: %02d:%02d", mins, secs))
		s += "\n\n"

		// Progress bar
		progressPercent := float64(m.elapsed) / float64(m.duration)
		if progressPercent > 1 {
			progressPercent = 1
		}
		s += m.progress.ViewAs(progressPercent)
		s += "\n"

		if m.isComplete {
			s += completeStyle.Render("Session complete!")
			s += "\n\n"
			s += instructionStyle.Render("Press 'r' to return to menu, 'q' to quit")
		} else if m.isRunning {
			s += instructionStyle.Render("Press 'space' to pause, 'q' to quit")
		} else {
			s += instructionStyle.Render("Paused - Press 'space' to resume, 'q' to quit")
		}

		s += "\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
