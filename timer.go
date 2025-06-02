package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	workDuration  = 30 * time.Minute
	breakDuration = 10 * time.Minute
)

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

func (s sessionType) duration() time.Duration {
	switch s {
	case work:
		return workDuration
	case breakSession:
		return breakDuration
	default:
		return 0
	}
}

func (s sessionType) next() sessionType {
	switch s {
	case work:
		return breakSession
	case breakSession:
		return work
	default:
		return none
	}
}

type model struct {
	session         sessionType
	progress        progress.Model
	duration        time.Duration
	elapsed         time.Duration
	isRunning       bool
	isComplete      bool
	showMenu        bool
	autoTransition  bool
	waitingForSound bool
}

type soundCompleteMsg struct{}

func soundCompleteCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second) // Wait for sound to finish
		return soundCompleteMsg{}
	}
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func initialModel() model {
	p := progress.New(
		progress.WithScaledGradient("#282b59", "#3c2859"),
		progress.WithWidth(40),
	)

	return model{
		progress: p,
		showMenu: true,
	}
}

func (m *model) startSession(session sessionType) {
	m.session = session
	m.duration = session.duration()
	m.elapsed = 0
	m.isRunning = true
	m.isComplete = false
	m.showMenu = false
}

func (m *model) completeSession() {
	m.isComplete = true
	m.isRunning = false

	m.playSessionEndSound()

	if m.autoTransition {
		m.waitingForSound = true
	}
}

func (m *model) resetToMenu() {
	clearScreen()
	m.showMenu = true
	m.isRunning = false
	m.isComplete = false
	m.elapsed = 0
	m.autoTransition = false
}

func (m *model) togglePause() {
	if m.isComplete {
		return
	}
	m.isRunning = !m.isRunning
}

func (m *model) playSessionEndSound() {
	go func() {
		// Play system sound with lower volume
		exec.Command("paplay", "--volume=19661", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
	}()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showMenu {
			return m.handleMenuInput(msg)
		}
		return m.handleSessionInput(msg)

	case tickMsg:
		if m.isRunning && !m.isComplete {
			m.elapsed += time.Second
			if m.elapsed >= m.duration {
				m.completeSession()
				if !m.autoTransition {
					clearScreen()
				} else {
					return m, soundCompleteCmd()
				}
			} else {
				return m, tickCmd()
			}
		}

	case soundCompleteMsg:
		if m.waitingForSound && m.autoTransition {
			m.waitingForSound = false
			nextSession := m.session.next()
			if nextSession != none {
				m.startSession(nextSession)
				m.isComplete = false
				return m, tickCmd()
			}
		}
	}

	if m.isRunning && !m.isComplete {
		return m, tickCmd()
	}
	return m, nil
}

func (m model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1", "w", "work":
		m.startSession(work)
		return m, tickCmd()
	case "2", "b", "break":
		m.startSession(breakSession)
		return m, tickCmd()
	case "3", "a", "auto":
		m.autoTransition = true
		m.startSession(work)
		return m, tickCmd()
	case "q", "quit", "ctrl+c":
		clearScreen()
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleSessionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "quit", "ctrl+c":
		clearScreen()
		return m, tea.Quit
	case "r", "return":
		if m.isComplete {
			m.resetToMenu()
		}
		return m, nil
	case " ":
		m.togglePause()
		return m, nil
	case "s", "skip":
		m.completeSession()
		if !m.autoTransition {
			clearScreen()
		}
		return m, nil
	}
	return m, nil
}

func (m model) renderMenu() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8e9ee8")).
		Bold(true).
		Padding(1, 2)

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#317f85")).
		Padding(0, 2)

	var s string
	s += titleStyle.Render("Pomodoro Timer")
	s += "\n\n"
	s += menuStyle.Render("Choose an option:")
	s += "\n"
	s += menuStyle.Render("  1 / w  →  Work session (30 minutes)")
	s += "\n"
	s += menuStyle.Render("  2 / b  →  Break session (10 minutes)")
	s += "\n"
	s += menuStyle.Render("  3 / a  →  Auto cycle (work → break → work...)")
	s += "\n"
	s += menuStyle.Render("  q      →  Quit")
	s += "\n\n"

	return s
}

func (m model) renderSession() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8e9ee8")).
		Bold(true).
		Padding(1, 2)

	sessionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a38ee8")).
		Bold(true)

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6e4fb3")).
		Bold(true)

	completeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6269a8")).
		Bold(true)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	var s string

	title := fmt.Sprintf("%s Session", m.session)
	if m.autoTransition {
		title += " (Auto Mode)"
	}
	s += titleStyle.Render(title)
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

	progressPercent := float64(m.elapsed) / float64(m.duration)
	if progressPercent > 1 {
		progressPercent = 1
	}
	s += m.progress.ViewAs(progressPercent)
	s += "\n"

	if m.isComplete && !m.autoTransition {
		s += completeStyle.Render("Session complete!")
		s += "\n\n"
		s += instructionStyle.Render("Press 'r' to return to menu, 'q' to quit")
	} else if m.isRunning {
		instructions := "Press 'space' to pause, 's' to skip, 'q' to quit"
		s += instructionStyle.Render(instructions)
	} else if !m.isComplete {
		instructions := "Paused - Press 'space' to resume, 's' to skip, 'q' to quit"
		s += instructionStyle.Render(instructions)
	}

	s += "\n"
	return s
}

func (m model) View() string {
	if m.showMenu {
		return m.renderMenu()
	}
	return m.renderSession()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
