package ui

import (
	"charm-pokemon/assets"
	"charm-pokemon/models"
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PokedexState int

const (
	StatePokedexView PokedexState = iota
	StateSearch
	StateBrowseType
	StateBrowseGeneration
	StateBrowseGenerationList
	StateFavorites
	StateDetail
)

type RenderMode int

const (
	RenderHalfBlock RenderMode = iota
	RenderSixel
)

type PokedexModel struct {
	state     PokedexState
	pokedex   *models.Pokedex
	favorites *models.FavoritesManager

	currentPokemon *models.Pokemon
	showShiny      bool

	searchQuery         string
	searchResults       []*models.Pokemon
	selectedSearchIndex int

	typeCursor           int
	generationCursor     int
	generationListCursor int
	favoritesCursor      int

	pokemonList       []*models.Pokemon
	pokemonListCursor int

	menuCursor int

	renderMode RenderMode

	// Terminal dimensions for responsive layout
	width  int
	height int
}

func NewPokedexModel(pokedex *models.Pokedex, favorites *models.FavoritesManager) PokedexModel {
	// Initialize current pokemon to the first one in the list
	var initialPokemon *models.Pokemon
	if pokedex != nil && len(pokedex.Pokemon) > 0 {
		initialPokemon = pokedex.Pokemon[0]
		if favorites != nil {
			initialPokemon.IsFavorite = favorites.IsFavorite(initialPokemon.ID)
		}
	}

	return PokedexModel{
		state:                StatePokedexView,
		pokedex:              pokedex,
		favorites:            favorites,
		currentPokemon:       initialPokemon,
		showShiny:            false,
		searchQuery:          "",
		searchResults:        make([]*models.Pokemon, 0),
		selectedSearchIndex:  0,
		typeCursor:           0,
		generationCursor:     0,
		generationListCursor: 0,
		favoritesCursor:      0,
		pokemonList:          make([]*models.Pokemon, 0),
		pokemonListCursor:    0,
		menuCursor:           0,
		width:                80, // Default width
		height:               24, // Default height
	}
}

func (m PokedexModel) Init() tea.Cmd {
	return nil
}

func (m PokedexModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch m.state {
		case StatePokedexView:
			return m.updatePokedexView(msg)
		case StateSearch:
			return m.updateSearch(msg)
		case StateBrowseType:
			return m.updateBrowseType(msg)
		case StateBrowseGeneration:
			return m.updateBrowseGeneration(msg)
		case StateBrowseGenerationList:
			return m.updateBrowseGenerationList(msg)
		case StateFavorites:
			return m.updateFavorites(msg)
		case StateDetail:
			return m.updateDetail(msg)
		}
	}
	return m, nil
}

func (m PokedexModel) View() string {
	switch m.state {
	case StatePokedexView:
		return m.viewPokedex()
	case StateSearch:
		return m.viewSearch()
	case StateBrowseType:
		return m.viewBrowseType()
	case StateBrowseGeneration:
		return m.viewBrowseGeneration()
	case StateBrowseGenerationList:
		return m.viewBrowseGenerationList()
	case StateFavorites:
		return m.viewFavorites()
	case StateDetail:
		return m.viewDetail()
	default:
		return "Estado desconhecido"
	}
}

func (m PokedexModel) viewPokedex() string {
	pokemon := m.GetCurrentPokemon()

	var s strings.Builder

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(LabelPOKEDEX),
			"",
		),
	))

	if pokemon != nil {
		art := m.loadPokemonArt(pokemon)

		// Calculate appropriate width - use terminal width or default
		artWidth := 65
		if m.width > 0 && m.width < artWidth+10 {
			artWidth = m.width - 10
		}

		// Apply type-based coloring ONLY if art is not already colored (Braille legacy)
		artStyle := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(artWidth)

		// If the art doesn't contain ANSI color codes, apply type color
		if !strings.Contains(art, "\x1b[") && len(pokemon.Types) > 0 {
			artStyle = artStyle.Foreground(getTypeColor(pokemon.Types[0]))
		}

		s.WriteString("\n\n")
		s.WriteString(artStyle.Render(art))
		s.WriteString("\n\n")

		// Show all type emojis
		typeEmojis := ""
		for _, t := range pokemon.Types {
			typeEmojis += getTypeEmoji(t) + " "
		}

		s.WriteString(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.width).
			Render(fmt.Sprintf("#%d %s %s\n", pokemon.ID, pokemon.NamePT, typeEmojis)))
		s.WriteString("\n")

		menuItems := []struct {
			label  string
			hotkey string
		}{
			{LabelSEARCH, "1"},
			{LabelBROWSE_TYPES, "2"},
			{LabelBROWSE_GEN, "3"},
			{LabelFAVORITES, "4"},
		}

		// Calculate max width for menu alignment
		maxWidth := 0
		var menuStrings []string
		for _, item := range menuItems {
			str := fmt.Sprintf("[%s] %s", item.hotkey, item.label)
			if len(str) > maxWidth {
				maxWidth = len(str)
			}
			menuStrings = append(menuStrings, str)
		}

		for _, str := range menuStrings {
			s.WriteString(lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(m.width).
				Render(fmt.Sprintf("%-*s", maxWidth, str)) + "\n")
		}

		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.width).
			Render(fmt.Sprintf("%s   Enter para detalhes   %s", LabelPREV, LabelNEXT)))
		s.WriteString("\n\n")

		s.WriteString(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.width).
			Faint(true).
			Render(LabelPRESS_Q))
	}

	return s.String()
}

func (m PokedexModel) viewSearch() string {
	var s strings.Builder

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(LabelSEARCH),
			"",
		),
	))

	s.WriteString("\n")
	s.WriteString(getLabelStyle().Render(LabelSEARCH_QUERY))
	s.WriteString("\n")

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render("> " + m.searchQuery + "_"))
	s.WriteString("\n\n")

	s.WriteString(getLabelStyle().Render(LabelRESULTS))
	s.WriteString("\n\n")

	if len(m.searchResults) > 0 {
		maxResults := 8
		startIdx := 0
		if m.selectedSearchIndex >= maxResults {
			startIdx = m.selectedSearchIndex - maxResults + 1
		}

		for i := startIdx; i < len(m.searchResults) && i < startIdx+maxResults; i++ {
			pokemon := m.searchResults[i]
			cursor := " "
			if i == m.selectedSearchIndex {
				cursor = ">"
			}

			style := getNormalItemStyle()
			if i == m.selectedSearchIndex {
				style = getCursorStyle()
			}

			typeEmoji := ""
			if len(pokemon.Types) > 0 {
				typeEmoji = getTypeEmoji(pokemon.Types[0])
			}

			s.WriteString(style.Render(fmt.Sprintf("%s #%4d %-20s %s\n", cursor, pokemon.ID, pokemon.NamePT, typeEmoji)))
		}
	} else if m.searchQuery != "" {
		s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelNO_RESULTS))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelPRESS_ENTER))

	return s.String()
}

func (m PokedexModel) viewBrowseType() string {
	var s strings.Builder

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(LabelTYPES),
			"",
		),
	))

	s.WriteString("\n\n")

	for i, typeName := range TypeNames {
		cursor := " "
		if i == m.typeCursor {
			cursor = ">"
		}

		style := getNormalItemStyle()
		if i == m.typeCursor {
			style = getCursorStyle()
		}

		typeEmoji := getTypeEmoji(typeName)
		count := len(m.pokedex.GetPokemonByType(typeName))

		s.WriteString(style.Render(fmt.Sprintf("%s %s %-12s - %3d %s\n", cursor, typeEmoji, typeName, count, LabelPOKEMON)))
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelPRESS_ENTER))

	return s.String()
}

func (m PokedexModel) viewBrowseGeneration() string {
	var s strings.Builder

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(LabelGENERATIONS),
			"",
		),
	))

	s.WriteString("\n\n")

	for i, gen := range Generations {
		cursor := " "
		if i == m.generationCursor {
			cursor = ">"
		}

		style := getNormalItemStyle()
		if i == m.generationCursor {
			style = getCursorStyle()
		}

		count := len(m.pokedex.GetPokemonByGeneration(gen.ID))

		s.WriteString(style.Render(fmt.Sprintf("%s %-20s (%-10s) - %3d %s\n", cursor, gen.NamePT, gen.Region, count, LabelPOKEMON)))
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelPRESS_ENTER))

	return s.String()
}

func (m PokedexModel) viewBrowseGenerationList() string {
	var s strings.Builder

	currentGen := Generations[m.generationCursor]

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(fmt.Sprintf("%s - %s", currentGen.NamePT, currentGen.Region)),
			"",
		),
	))

	s.WriteString("\n\n")

	maxResults := 10
	startIdx := 0
	if m.generationListCursor >= maxResults {
		startIdx = m.generationListCursor - maxResults + 1
	}

	for i := startIdx; i < len(m.pokemonList) && i < startIdx+maxResults; i++ {
		pokemon := m.pokemonList[i]
		cursor := " "
		if i == m.generationListCursor {
			cursor = ">"
		}

		style := getNormalItemStyle()
		if i == m.generationListCursor {
			style = getCursorStyle()
		}

		typeEmoji := ""
		if len(pokemon.Types) > 0 {
			typeEmoji = getTypeEmoji(pokemon.Types[0])
		}

		s.WriteString(style.Render(fmt.Sprintf("%s #%4d %-20s %s\n", cursor, pokemon.ID, pokemon.NamePT, typeEmoji)))
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelPRESS_ENTER))

	return s.String()
}

func (m PokedexModel) viewFavorites() string {
	var s strings.Builder

	s.WriteString(getBoxStyle().Render(
		lipgloss.JoinVertical(lipgloss.Center,
			getTitleStyle().Render(LabelFAVORITES),
			"",
		),
	))

	s.WriteString("\n\n")

	if len(m.pokemonList) > 0 {
		maxResults := 10
		startIdx := 0
		if m.favoritesCursor >= maxResults {
			startIdx = m.favoritesCursor - maxResults + 1
		}

		for i := startIdx; i < len(m.pokemonList) && i < startIdx+maxResults; i++ {
			pokemon := m.pokemonList[i]
			cursor := " "
			if i == m.favoritesCursor {
				cursor = ">"
			}

			style := getNormalItemStyle()
			if i == m.favoritesCursor {
				style = getCursorStyle()
			}

			typeEmoji := ""
			if len(pokemon.Types) > 0 {
				typeEmoji = getTypeEmoji(pokemon.Types[0])
			}

			s.WriteString(style.Render(fmt.Sprintf("%s #%d %s %s\n", cursor, pokemon.ID, pokemon.NamePT, typeEmoji)))
		}

		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("%s: %d %s\n", LabelTOTAL, len(m.pokemonList), LabelPOKEMON)))
	} else {
		s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelNO_FAVORITES))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(LabelPRESS_ENTER))

	return s.String()
}

func (m PokedexModel) viewDetail() string {
	pokemon := m.GetCurrentPokemon()
	if pokemon == nil {
		return "Nenhum Pokémon selecionado"
	}

	var s strings.Builder

	art := m.loadPokemonArt(pokemon)

	typeEmojis := ""
	for _, t := range pokemon.Types {
		typeEmojis += getTypeEmoji(t) + " "
	}

	s.WriteString(getHeaderStyle().Render(fmt.Sprintf("#%d %s %s", pokemon.ID, pokemon.NamePT, typeEmojis)))
	s.WriteString("\n")

	// Calculate appropriate width - use terminal width or default
	artWidth := 65
	if m.width > 0 && m.width < artWidth+10 {
		artWidth = m.width - 10
	}

	// Apply type-based coloring ONLY if art is not already colored
	artStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(artWidth)
	if !strings.Contains(art, "\x1b[") && len(pokemon.Types) > 0 {
		artStyle = artStyle.Foreground(getTypeColor(pokemon.Types[0]))
	}

	s.WriteString(artStyle.Render(art))
	s.WriteString("\n\n")

	s.WriteString(lipgloss.NewStyle().Render(fmt.Sprintf("%s %.1fm   %s %.1fkg", LabelHEIGHT, pokemon.Height/10.0, LabelWEIGHT, pokemon.Weight/10.0)))

	favStatus := ""
	if pokemon.IsFavorite {
		favStatus = " ⭐"
	}

	s.WriteString("  ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(LabelTOGGLE_FAVORITE))
	s.WriteString(favStatus)
	s.WriteString("\n\n")

	// Fixed shiny toggle - show which mode is active
	normalStyle := lipgloss.NewStyle()
	shinyStyle := lipgloss.NewStyle()
	if m.showShiny {
		shinyStyle = shinyStyle.Bold(true).Foreground(lipgloss.Color("226"))
		s.WriteString(fmt.Sprintf("[ %s ]  [%s %s] ◄", normalStyle.Render(LabelNORMAL), shinyStyle.Render(LabelSHINY), ""))
	} else {
		normalStyle = normalStyle.Bold(true).Foreground(lipgloss.Color("39"))
		s.WriteString(fmt.Sprintf("◄ [%s]  [ %s ]", normalStyle.Render(LabelNORMAL), shinyStyle.Render(LabelSHINY)))
	}
	s.WriteString("\n\n")

	s.WriteString(getLabelStyle().Render(LabelSTATS))
	s.WriteString("\n")

	stats := []struct {
		name  string
		value int
	}{
		{"HP", pokemon.Stats.HP},
		{"Ataque", pokemon.Stats.Attack},
		{"Defesa", pokemon.Stats.Defense},
		{"Sp.Atk", pokemon.Stats.SpAtk},
		{"Sp.Def", pokemon.Stats.SpDef},
		{"Veloc.", pokemon.Stats.Speed},
	}

	for _, stat := range stats {
		s.WriteString(fmt.Sprintf("  %-10s ", stat.name))
		s.WriteString(renderStatBar(stat.value, 150))
		s.WriteString("\n")
	}

	if pokemon.Evolution != nil {
		s.WriteString("\n")
		s.WriteString(getLabelStyle().Render(LabelEVOLUTION))
		s.WriteString("\n")

		evolutionStr := pokemon.Evolution.Base.Name
		currentStage := pokemon.Evolution.FindStage(pokemon.ID)

		for i, evo := range pokemon.Evolution.Evolution {
			trigger := ""
			if evo.MinLevel > 0 {
				trigger = fmt.Sprintf("(Lv%d)", evo.MinLevel)
			} else if evo.Item != "" {
				trigger = fmt.Sprintf("(%s)", evo.Item)
			}

			evolutionStr += fmt.Sprintf(" %s %s", trigger, evo.Name)

			if i == currentStage-1 {
				evolutionStr += " ←"
			}
		}

		s.WriteString(lipgloss.NewStyle().Render(fmt.Sprintf("  %s\n", evolutionStr)))
	}

	if len(pokemon.SignatureMoves) > 0 {
		s.WriteString("\n")
		s.WriteString(getLabelStyle().Render(LabelMOVES))
		s.WriteString("\n")

		for _, move := range pokemon.SignatureMoves {
			typeEmoji := getTypeEmoji(move.Type)
			s.WriteString(fmt.Sprintf("  • %s (%s %s) - %d poder\n", move.NamePT, typeEmoji, move.Type, move.Power))
		}
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("[s] Alternar Shiny   [f] Favorito   [%s / %s] Navegar   [q] Voltar", LabelPREV, LabelNEXT)))

	return s.String()
}

func (m PokedexModel) updatePokedexView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit

	case "left", "h":
		if m.currentPokemon != nil {
			m.currentPokemon = m.pokedex.GetPrevPokemon(m.currentPokemon.ID)
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}

	case "right", "l":
		if m.currentPokemon != nil {
			m.currentPokemon = m.pokedex.GetNextPokemon(m.currentPokemon.ID)
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}

	case "1":
		m.state = StateSearch
		m.searchQuery = ""
		m.searchResults = make([]*models.Pokemon, 0)
		return m, nil

	case "v":
		if m.renderMode == RenderHalfBlock {
			m.renderMode = RenderSixel
		} else {
			m.renderMode = RenderHalfBlock
		}
		return m, nil

	case "2":
		m.state = StateBrowseType
		m.typeCursor = 0
		return m, nil

	case "3":
		m.state = StateBrowseGeneration
		m.generationCursor = 0
		return m, nil

	case "4":
		m.state = StateFavorites
		favIDs := m.favorites.GetAllFavorites()
		sort.Ints(favIDs) // Sort favorites by ID
		m.pokemonList = make([]*models.Pokemon, 0)
		for _, id := range favIDs {
			if pokemon := m.pokedex.GetByID(id); pokemon != nil {
				m.pokemonList = append(m.pokemonList, pokemon)
			}
		}
		m.favoritesCursor = 0
		return m, nil

	case "enter", " ":
		if m.currentPokemon != nil {
			m.state = StateDetail
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}
	}
	return m, nil
}

func (m PokedexModel) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StatePokedexView
		return m, nil

	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.updateSearchResults()
		}

	case "enter":
		if len(m.searchResults) > 0 && m.selectedSearchIndex < len(m.searchResults) {
			m.currentPokemon = m.searchResults[m.selectedSearchIndex]
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
			m.state = StatePokedexView
		}

	case "up", "k":
		if m.selectedSearchIndex > 0 {
			m.selectedSearchIndex--
		}

	case "down", "j":
		if m.selectedSearchIndex < len(m.searchResults)-1 {
			m.selectedSearchIndex++
		}

	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
			m.updateSearchResults()
		}
	}
	return m, nil
}

func (m *PokedexModel) updateSearchResults() {
	if m.searchQuery == "" {
		m.searchResults = make([]*models.Pokemon, 0)
		m.selectedSearchIndex = 0
		return
	}

	filter := models.PokemonFilter{
		Query: m.searchQuery,
	}
	m.searchResults = m.pokedex.Search(filter)
	m.selectedSearchIndex = 0
}

func (m PokedexModel) updateBrowseType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StatePokedexView
		return m, nil

	case "up", "k":
		if m.typeCursor > 0 {
			m.typeCursor--
		}

	case "down", "j":
		if m.typeCursor < len(TypeNames)-1 {
			m.typeCursor++
		}

	case "enter", " ":
		selectedType := TypeNames[m.typeCursor]
		m.pokemonList = m.pokedex.GetPokemonByType(selectedType)
		m.pokemonListCursor = 0
		if len(m.pokemonList) > 0 {
			m.currentPokemon = m.pokemonList[0]
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}
		m.state = StatePokedexView
	}
	return m, nil
}

func (m PokedexModel) updateBrowseGeneration(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StatePokedexView
		return m, nil

	case "up", "k":
		if m.generationCursor > 0 {
			m.generationCursor--
		}

	case "down", "j":
		if m.generationCursor < len(Generations)-1 {
			m.generationCursor++
		}

	case "enter", " ":
		selectedGen := Generations[m.generationCursor]
		m.pokemonList = m.pokedex.GetPokemonByGeneration(selectedGen.ID)
		m.generationListCursor = 0
		if len(m.pokemonList) > 0 {
			m.state = StateBrowseGenerationList
		}
	}
	return m, nil
}

func (m PokedexModel) updateBrowseGenerationList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StateBrowseGeneration
		return m, nil

	case "up", "k":
		if m.generationListCursor > 0 {
			m.generationListCursor--
		}

	case "down", "j":
		if m.generationListCursor < len(m.pokemonList)-1 {
			m.generationListCursor++
		}

	case "enter", " ":
		if m.generationListCursor < len(m.pokemonList) {
			m.currentPokemon = m.pokemonList[m.generationListCursor]
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
			m.state = StatePokedexView
		}
	}
	return m, nil
}

func (m PokedexModel) updateFavorites(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StatePokedexView
		return m, nil

	case "up", "k":
		if m.favoritesCursor > 0 {
			m.favoritesCursor--
		}

	case "down", "j":
		if m.favoritesCursor < len(m.pokemonList)-1 {
			m.favoritesCursor++
		}

	case "enter", " ":
		if m.favoritesCursor < len(m.pokemonList) {
			m.currentPokemon = m.pokemonList[m.favoritesCursor]
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
			m.state = StatePokedexView
		}
	}
	return m, nil
}

func (m PokedexModel) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StatePokedexView
		return m, nil

	case "s":
		m.showShiny = !m.showShiny

	case "f":
		if m.currentPokemon != nil {
			isFav, _ := m.favorites.ToggleFavorite(m.currentPokemon.ID)
			m.currentPokemon.IsFavorite = isFav
		}

	case "left", "h":
		if m.currentPokemon != nil {
			m.currentPokemon = m.pokedex.GetPrevPokemon(m.currentPokemon.ID)
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}

	case "right", "l":
		if m.currentPokemon != nil {
			m.currentPokemon = m.pokedex.GetNextPokemon(m.currentPokemon.ID)
			m.currentPokemon.IsFavorite = m.favorites.IsFavorite(m.currentPokemon.ID)
		}
	}
	return m, nil
}

func (m PokedexModel) GetCurrentPokemon() *models.Pokemon {
	if m.currentPokemon == nil && m.pokedex != nil && len(m.pokedex.Pokemon) > 0 {
		m.currentPokemon = m.pokedex.Pokemon[0]
	}
	return m.currentPokemon
}

func (m PokedexModel) loadPokemonArt(pokemon *models.Pokemon) string {
	suffix := ""
	if m.showShiny {
		suffix = "_shiny"
	}

	// For ASCII mode, try embedded FS first (optimized binary)
	if m.renderMode == RenderHalfBlock {
		embeddedPath := fmt.Sprintf("embed/art/%d%s.ascii", pokemon.ID, suffix)
		data, err := assets.EmbedFS.ReadFile(embeddedPath)
		if err == nil {
			return string(data)
		}
	}

	// Determine extension based on render mode
	ext := ".ascii"
	if m.renderMode == RenderSixel {
		ext = ".sixel"
	}

	// Try to load from disk (for Sixel mode or fallback)
	path := fmt.Sprintf("assets/art/%d%s%s", pokemon.ID, suffix, ext)
	data, err := os.ReadFile(path)
	if err == nil {
		return string(data)
	}

	// Fallback to legacy hardcoded art
	if m.showShiny {
		return pokemon.ArtShiny
	}
	return pokemon.ArtStandard
}
