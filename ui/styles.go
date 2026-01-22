package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	typeColors = map[string]lipgloss.Color{
		"normal":   lipgloss.Color("248"),
		"fogo":     lipgloss.Color("208"),
		"água":     lipgloss.Color("27"),
		"grama":    lipgloss.Color("82"),
		"elétrico": lipgloss.Color("226"),
		"gelo":     lipgloss.Color("45"),
		"lutador":  lipgloss.Color("160"),
		"veneno":   lipgloss.Color("153"),
		"terra":    lipgloss.Color("172"),
		"voador":   lipgloss.Color("163"),
		"psíquico": lipgloss.Color("203"),
		"inseto":   lipgloss.Color("166"),
		"pedra":    lipgloss.Color("179"),
		"fantasma": lipgloss.Color("111"),
		"dragão":   lipgloss.Color("169"),
		"sombrio":  lipgloss.Color("88"),
		"metálico": lipgloss.Color("201"),
		"fada":     lipgloss.Color("197"),
	}
)

func getTypeColor(typeName string) lipgloss.Color {
	if color, ok := typeColors[typeName]; ok {
		return color
	}
	return lipgloss.Color("255")
}

func getTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)
}

func getHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		MarginBottom(1)
}

func getLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)
}

func getValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
}

func getCursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)
}

func getNormalItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
}

func getSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)
}

func getTypeStyle(typeName string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(getTypeColor(typeName)).
		Bold(true)
}

func getHighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Background(lipgloss.Color("88"))
}

func getBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)
}

func getStatBarStyle(stat int, maxValue int) string {
	width := 15
	if maxValue <= 0 {
		maxValue = 150 // Default max for Pokemon stats
	}
	filled := int(float64(stat) / float64(maxValue) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func renderStatBar(stat int, maxValue int) string {
	bar := getStatBarStyle(stat, maxValue)
	return lipgloss.JoinHorizontal(lipgloss.Left,
		bar,
		lipgloss.NewStyle().Width(4).Render(""),
		lipgloss.NewStyle().Render(fmt.Sprintf("%3d", stat)),
	)
}

func getTypeEmoji(typeName string) string {
	emojis := map[string]string{
		"normal":   "⚪",
		"fogo":     "🔥",
		"água":     "💧",
		"grama":    "🌿",
		"elétrico": "⚡",
		"gelo":     "❄️",
		"lutador":  "👊",
		"veneno":   "☠️",
		"terra":    "🌍",
		"voador":   "🕊️",
		"psíquico": "🔮",
		"inseto":   "🐛",
		"pedra":    "🪨",
		"fantasma": "👻",
		"dragão":   "🐉",
		"sombrio":  "🌑",
		"metálico": "⚙️",
		"fada":     "🧚",
	}

	if emoji, ok := emojis[typeName]; ok {
		return emoji
	}
	return "⚪"
}

const (
	LabelPOKEDEX         = "📖 POKÉDEX"
	LabelSEARCH          = "🔍 Buscar"
	LabelBROWSE_TYPES    = "🎨 Tipos"
	LabelBROWSE_GEN      = "📚 Gerações"
	LabelFAVORITES       = "⭐ Favoritos"
	LabelDETAILS         = "📊 Detalhes"
	LabelPREV            = "◀ Anterior"
	LabelNEXT            = "Próximo ▶"
	LabelBACK            = "◀ Voltar"
	LabelSEARCH_QUERY    = "Digita o nome ou número do Pokémon:"
	LabelRESULTS         = "Resultados:"
	LabelTYPE            = "Tipo:"
	LabelHEIGHT          = "Altura:"
	LabelWEIGHT          = "Peso:"
	LabelSTATS           = "Estatísticas:"
	LabelEVOLUTION       = "Evolução:"
	LabelMOVES           = "Movimentos Característicos:"
	LabelTOTAL           = "Total:"
	LabelPOKEMON         = "Pokémon"
	LabelGENERATION      = "Geração:"
	LabelSHINY           = "Shiny ✨"
	LabelNORMAL          = "Normal"
	LabelNO_RESULTS      = "Nenhum resultado encontrado"
	LabelNO_FAVORITES    = "Nenhum favorito ainda"
	LabelPRESS_ENTER     = "Pressiona Enter para selecionar"
	LabelPRESS_Q         = "Pressiona q para voltar"
	LabelTOGGLE_FAVORITE = "⭐ Favorito"
	LabelGENERATIONS     = "Navegar por Geração"
	LabelTYPES           = "Navegar por Tipo"
)

var TypeNames = []string{
	"normal",
	"fogo",
	"água",
	"grama",
	"elétrico",
	"gelo",
	"lutador",
	"veneno",
	"terra",
	"voador",
	"psíquico",
	"inseto",
	"pedra",
	"fantasma",
	"dragão",
	"sombrio",
	"metálico",
	"fada",
}

var Generations = []struct {
	ID     int
	NamePT string
	NameEN string
	Region string
}{
	{1, "Primeira Geração", "Generation I", "Kanto"},
	{2, "Segunda Geração", "Generation II", "Johto"},
	{3, "Terceira Geração", "Generation III", "Hoenn"},
	{4, "Quarta Geração", "Generation IV", "Sinnoh"},
	{5, "Quinta Geração", "Generation V", "Unova"},
	{6, "Sexta Geração", "Generation VI", "Kalos"},
	{7, "Sétima Geração", "Generation VII", "Alola"},
	{8, "Oitava Geração", "Generation VIII", "Galar"},
	{9, "Nona Geração", "Generation IX", "Paldea"},
}
