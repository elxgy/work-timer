# work-timer
Another project to pretend you're productive

## Features

- Work sessions: 30 minutes
- Break sessions: 10 minutes
- Sound feedback at session end
- Pause/resume functionality

## Prerequisites

- Go 1.19 or later
- Terminal with color support

## Installation

### 1. Clone or Download

```bash
git clone https://github.com/elxgy/work-timer
cd timer
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Build the Executable

```bash
go build -o timer timer.go
```

### 4. Set Up Terminal Alias

Choose one of the following methods:

#### Option A: Global Installation (Recommended)

```bash
# Move to system PATH
sudo mv timer /usr/local/bin/timer
sudo chmod +x /usr/local/bin/timer
```

Now you can use `timer` from anywhere in your terminal!

#### Option B: Shell Alias

**For Bash** (add to `~/.bashrc` or `~/.bash_profile`):
```bash
echo 'alias timer="/path/to/your/timer/timer"' >> ~/.bashrc
source ~/.bashrc
```

**For Zsh** (add to `~/.zshrc`):
```bash
echo 'alias timer="/path/to/your/timer/timer"' >> ~/.zshrc
source ~/.zshrc
```

**For Fish** (add to `~/.config/fish/config.fish`):
```bash
echo 'alias timer="/path/to/your/timer/timer"' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

#### Option C: Direct Go Run Alias

```bash
# Add to your shell config file
alias timer="cd /path/to/your/timer/directory && go run timer.go"
```

## Usage

### Starting the Timer

Simply run:
```bash
timer
```

### Main Menu Options

When you start the timer, you'll see:

```
Timer

Choose an option:
  1 / w  →  Work session (30 minutes)
  2 / b  →  Break session (10 minutes)
  q      →  Quit
```

### Keyboard Controls

#### In Main Menu:
- `1` or `w` - Start work session
- `2` or `b` - Start break session
- `3` or `a` - Auto cycle (work → break → work...)  
- `q` - Quit application

#### During Timer Session:
- `Space` - Pause/resume timer
- `q` - Quit application
- `r` - Return to main menu (when session is complete)

### Example Workflow

1. Run `timer` in your terminal
2. Press `1` or `w` to start a 30-minute work session
3. Use `Space` to pause if needed
4. When complete, press `r` to return to menu
5. Press `2` or `b` for a 10-minute break
6. Press `q` to quit when done

## Customization

You can modify the session durations in `timer.go`:

```go
const (
    workDuration  = 30 * time.Minute  // Change work duration
    breakDuration = 10 * time.Minute  // Change break duration
)
```

After making changes, rebuild:
```bash
go build -o timer timer.go
```

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

### Building from Source
```bash
go mod download
go build -o timer timer.go
```

## License

MIT License - feel free to use and modify as needed.