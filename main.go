package main

import (
	"charm-pokemon/data"
	"charm-pokemon/models"
	"charm-pokemon/ui"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pikachu ASCII Art - using raw string literal to preserve formatting
const pikachuArt = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡾⠋⠉⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⣠⠖⠲⢤⡖⠒⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⡏⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⡏⠀⠀⠀⠀⠀⢀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⣀⠀⠀⢸⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠸⣄⠀⠁⣠⠞⠉⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣠⣤⣤⣤⣤⠀⠀
⠀⡞⠉⠻⠁⢹⠀⠀⡏⠀⠀⠀⠀⢸⠃⠀⠀⠀⠀⠀⠀⠀⠀⠹⣶⠋⠀⠀⠀⠀⣀⡤⠴⠒⠊⠉⠉⠀⠀⣿⣿⣿⠿⠋⠀⠀
⠀⠳⢤⡀⠀⡞⠁⠀⡇⠀⠀⢀⡠⠼⠴⠒⠒⠒⠒⠦⠤⠤⣄⣀⠀⢀⣠⠴⠚⠉⠀⠀⠀⠀⠀⠀⠀⠀⣼⠿⠋⠁⠀⠀⠀⠀
⠀⠀⠀⠈⠷⡏⠀⠀⣇⠔⠂⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢨⠿⠋⠀⠀⠀⠀⠀⠀⠀⠀⣀⡤⠖⠋⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢰⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⠤⠒⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢀⡟⠀⣠⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⢻⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣤⣤⡤⠤⢴
⠀⠀⠀⠀⠀⠀⣸⠁⣾⣿⣀⣽⡆⠀⠀⠀⠀⠀⠀⠀⢠⣾⠉⢿⣦⠀⠀⠀⢸⡀⠀⠀⢀⣠⠤⠔⠒⠋⠉⠉⠀⠀⠀⠀⢀⡞
⠀⠀⠀⠀⠀⢀⡏⠀⠹⠿⠿⠟⠁⠀⠰⠦⠀⠀⠀⠀⠸⣿⣿⣿⡿⠀⠀⠀⢘⡧⠖⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡼⠀
⠀⠀⠀⠀⠀⣼⠦⣄⠀⠀⢠⣀⣀⣴⠟⠶⣄⡀⠀⠀⡀⠀⠉⠁⠀⠀⠀⠀⢸⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⠁⠀
⠀⠀⠀⠀⢰⡇⠀⠈⡇⠀⠀⠸⡾⠁⠀⠀⠀⠉⠉⡏⠀⠀⠀⣠⠖⠉⠓⢤⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⠃⠀⠀
⠀⠀⠀⠀⠀⢧⣀⡼⠃⠀⠀⠀⢧⠀⠀⠀⠀⠀⢸⠃⠀⠀⠀⣧⠀⠀⠀⣸⢹⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡰⠃⠀⠀⠀
⠀⠀⠀⠀⠀⠈⢧⡀⠀⠀⠀⠀⠘⣆⠀⠀⠀⢠⠏⠀⠀⠀⠀⠈⠳⠤⠖⠃⡟⠀⠀⠀⢾⠛⠛⠛⠛⠛⠛⠛⠛⠁⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠙⣆⠀⠀⠀⠀⠈⠦⣀⡴⠋⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⠙⢦⠀⠀⠘⡇⠀⠀⠀⠀⠀⠀⢀⣀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢠⡇⠙⠦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⠴⠋⠸⡇⠈⢳⡀⠀⢹⡀⠀⠀⠀⢀⡞⠁⠉⣇⣀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⡼⣀⠀⠀⠈⠙⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⣷⠴⠚⠁⠀⣀⣷⠀⠀⠀⢠⠇⠀⠀⠀⠀⠀⣳
⠀⠀⠀⠀⠀⠀⡴⠁⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣆⡴⠚⠉⠉⠀⠀⠀⠀⢸⠃⣀⣠⠤⠤⠖⠋
⣼⢷⡆⠀⣠⡴⠧⣄⣇⠀⠀⠀⠀⡴⠚⠙⠲⠞⠛⠙⡆⠀⠀⠀⠀⠀⢀⡇⣠⣽⢦⣄⢀⣴⣶⠀⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀
⡿⣼⣽⡞⠁⠀⠀⠀⢹⡀⠀⠀⠀⢹⠀⠀⠀⠀⠀⠀⣸⠀⠀⠀⠀⠀⣼⠉⠁⠀⠀⢠⢟⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⣷⠉⠁⢳⠀⠀⠀⠀⠈⣧⠀⠀⠀⠀⠙⢦⠀⠀⠀⡠⠁⠀⠀⠀⠀⣰⠃⠀⠀⠀⠀⠏⠀⠀⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠹⡆⠀⠈⡇⠀⠀⠀⠀⠘⣆⠀⠀⠀⠀⠀⠹⣧⠞⠁⠀⠀⠀⠀⣰⠃⠀⠀⠀⠀⠀⠀⠀⣸⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢳⡀⠀⠙⠀⠀⠀⠀⠀⠘⣆⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⣰⠃⠀⠀⠀⠀⢀⡄⠀⢠⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢳⡀⣰⣀⣀⣀⠀⠀⠀⠘⣦⣀⠀⠀⠀⡇⠀⠀⠀⢀⡴⠃⠀⠀⠀⠀⠀⢸⡇⢠⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠉⠉⠀⠀⠈⠉⠉⠉⠙⠻⠿⠾⠾⠻⠓⢦⠦⡶⡶⠿⠛⠛⠓⠒⠒⠚⠛⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

const (
	stateMainMenu = iota
	stateApps
	stateShutdown
	statePokedex
)

type model struct {
	choices      []string         // menu items
	cursor       int              // which menu item the cursor is pointing at
	selected     map[int]struct{} // which menu items are selected
	state        int              // current screen state
	appsChoices  []string         // apps menu items
	appsCursor   int              // cursor for the apps menu
	shutdownPerc int              // percentage for shutdown animation
	pokedex      *models.Pokedex
	favorites    *models.FavoritesManager
	pokedexModel ui.PokedexModel
}

func initialModel() model {
	pokedex := data.GetPokedex()
	favorites := models.NewFavoritesManager()

	return model{
		choices:      []string{"Pokedex", "Iniciar Apps", "Fechar o terminal"},
		selected:     make(map[int]struct{}),
		state:        stateMainMenu,
		appsChoices:  []string{"Browser (MS Edge)", "Bloco de Notas", "Voltar ao Menu Principal"},
		appsCursor:   0,
		shutdownPerc: 100,
		pokedex:      pokedex,
		favorites:    favorites,
		pokedexModel: ui.NewPokedexModel(pokedex, favorites),
	}
}

func (m model) renderPikachu() string {
	// Style for Pikachu - yellow color
	pikachuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")) // Bright Yellow

	artLines := strings.Split(pikachuArt, "\n")
	maxWidth := 0
	for _, line := range artLines {
		// Use lipgloss.Width for proper Unicode/Braille character visual width calculation
		lineWidth := lipgloss.Width(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	// Style for the welcome message
	welcomeMessage := "Olá Minês!"
	messageStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")). // Nice blue/cyan color
		Width(maxWidth).                  // Set an appropriate width for centering
		Align(lipgloss.Center)

	// Combine the styled elements vertically - Pikachu first, then welcome message
	return lipgloss.JoinVertical(lipgloss.Center,
		pikachuStyle.Render(pikachuArt),
		messageStyle.Render(welcomeMessage),
	)
}

func (m model) Init() tea.Cmd {

	return tea.SetWindowTitle("Pikachu Terminal")

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateMainMenu:
			return m.updateMainMenu(msg)
		case stateApps:
			return m.updateAppsMenu(msg)
		case stateShutdown:
			return m.updateShutdown(msg)
		case statePokedex:
			pokedexModel, cmd := m.pokedexModel.Update(msg)
			m.pokedexModel = pokedexModel.(ui.PokedexModel)

			if m.pokedexModel.GetCurrentPokemon() != nil && !m.pokedexModel.GetCurrentPokemon().IsFavorite {
				m.pokedexModel.GetCurrentPokemon().IsFavorite = m.favorites.IsFavorite(m.pokedexModel.GetCurrentPokemon().ID)
			}

			return m, cmd
		}

	case tickMsg:
		if m.state == stateShutdown {
			if m.shutdownPerc <= 0 {
				return m, tea.Quit
			}
			m.shutdownPerc -= 10
			return m, tick()
		}
	}

	return m, nil
}

func (m model) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	case "enter", " ":
		switch m.cursor {
		case 0: // Pokedex
			m.state = statePokedex
			m.pokedexModel = ui.NewPokedexModel(m.pokedex, m.favorites)
		case 1: // Open Apps
			m.state = stateApps
			m.appsCursor = 0
		case 2: // Power Off
			m.state = stateShutdown
			return m, tick()
		}
	}
	return m, nil
}

func (m model) updateAppsMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.appsCursor > 0 {
			m.appsCursor--
		}

	case "down", "j":
		if m.appsCursor < len(m.appsChoices)-1 {
			m.appsCursor++
		}

	case "enter", " ":
		switch m.appsCursor {
		case 0: // Open Browser (Edge)
			_ = openBrowser() // Ignore error - app launch is best effort
			// Remain in the same state
		case 1: // Open Notepad
			_ = openNotepad() // Ignore error - app launch is best effort
			// Remain in the same state
		case 2: // Back to main menu
			m.state = stateMainMenu
		}
	}
	return m, nil
}

func (m model) updateShutdown(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" || msg.String() == "q" {
		return m, tea.Quit
	}
	return m, nil
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func openBrowser() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", "msedge")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", "-a", "Microsoft Edge")
	} else {
		cmd = exec.Command("microsoft-edge")
	}
	return cmd.Start()
}

func openNotepad() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("notepad")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", "-a", "TextEdit")
	} else {
		cmd = exec.Command("gedit")
	}
	return cmd.Start()
}

func (m model) View() string {
	switch m.state {
	case stateMainMenu:
		return m.renderPikachu() + "\n" + m.mainMenuView()
	case stateApps:
		return m.renderPikachu() + "\n" + m.appsMenuView()
	case stateShutdown:
		return m.renderPikachu() + "\n" + m.shutdownView()
	case statePokedex:
		return m.pokedexModel.View()
	default:
		return "Error: unknown state"
	}
}

func (m model) mainMenuView() string {
	s := "Bem vinda ao Terminal Pikachu!\n\n"
	s += "Usa as setas para navegar, Enter para selecionar\n\n"

	// Iterate over choices
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(choice)
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPressiona q para sair\n"
	return s
}

func (m model) appsMenuView() string {
	s := "Apps disponíveis\n\n"

	// Iterate over app choices
	for i, choice := range m.appsChoices {
		cursor := " "
		if m.appsCursor == i {
			cursor = ">"
			choice = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(choice)
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPressiona q para sair\n"
	return s
}

func (m model) shutdownView() string {
	s := "A desligar...\n\n"

	// Create progress bar
	width := 50
	filled := int(float64(m.shutdownPerc) / 100.0 * float64(width))

	progress := ""
	for i := 0; i < width; i++ {
		if i < filled {
			progress += "█"
		} else {
			progress += "░"
		}
	}

	s += fmt.Sprintf("[%s] %d%%\n", progress, m.shutdownPerc)
	s += "\nPressiona Ctrl+C para fechar\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops: %v", err)
		os.Exit(1)
	}
}
