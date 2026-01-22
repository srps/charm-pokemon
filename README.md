# Ϟ Charm Pokemon Ϟ

A beautiful, high-performance Pokemon terminal application (TUI) built with Go and the [Charm](https://charm.sh/) libraries.

![Pokemon TUI](https://raw.githubusercontent.com/srps/charm-pokemon/main/.github/screenshot.png)

## ✨ Features

- **Full Pokedex**: Information on all 1,025 Pokemon from Generation 1 to 9.
- **Rich Graphics**:
  - **Half-block ASCII**: 24-bit color representations that work in any modern terminal.
  - **Sixel Support**: Pixel-perfect graphics for terminals that support the Sixel protocol (toggle with `v`).
- **Multilingual**: Comprehensive data in both Portuguese (PT-BR) and English.
- **Live Search**: Find Pokemon instantly by name or ID.
- **Smart Filters**: Browse by Type, Generation, or Region.
- **Favorites**: Mark and persist your favorite Pokemon.
- **Optimized Binary**: Advanced minification and embedding techniques keep the standalone binary size under 30MB.
- **App Launcher**: Integrated shortcuts to common system tools.

## 🚀 Getting Started

### 📦 Download Binaries

The easiest way to use charm-pokemon is to download a pre-built binary for your system from the **[Releases](https://github.com/yourusername/charm-pokemon/releases)** page.

1. Download the version for your OS (Windows, Linux, or macOS).
2. (Linux/macOS) Make it executable: `chmod +x charm-pokemon-*`
3. Run it!

### 🛠️ Building from Source

If you prefer to build it yourself, ensure you have **Go 1.25+** installed.

### Prerequisites

- Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/yourusername/charm-pokemon.git
cd charm-pokemon

# Build the optimized executable (embeds ~30MB of data)
go build -tags realdata -ldflags="-s -w" -o pokemon.exe .

# Run it!
./pokemon.exe
```

## 🎮 Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate menus and lists |
| `←/→` or `h/l` | Browse Pokemon in Pokedex |
| `Enter` | Select / View details |
| `1` | Open Search |
| `2` | Browse by Type |
| `3` | Browse by Generation |
| `4` | View Favorites |
| `s` | Toggle Normal/Shiny sprite (in detail view) |
| `v` | Toggle ASCII/Sixel rendering |
| `f` | Toggle favorite status |
| `q` / `Esc` | Back / Exit |

## 🛠️ Data & Optimization

The project uses a sophisticated data pipeline to minimize binary size while maintaining high quality:

1. **Downloader**: Fetches latest data from [PokeAPI](https://pokeapi.co/).
2. **Sprite Converter**: Generates high-fidelity ASCII and Sixel art.
3. **Data Minifier**: Strips unused API fields (movesets, URLs) to reduce JSON size by ~80%.
4. **Build Tags**: Uses `-tags realdata` to switch between sample development data and the full embedded dataset.

To rebuild the data from scratch:

```bash
go run tools/download_data/main.go
go run tools/convert_sprites/main.go
go run tools/minify_data/main.go
```

## 📦 Tech Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: The TUI framework.
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**: Terminal styling and layouts.
- **[Go-Sixel](https://github.com/mattn/go-sixel)**: Sixel graphics encoding.

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.

---

*Gotta catch 'em all... in your terminal!*
