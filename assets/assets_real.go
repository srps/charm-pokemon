//go:build realdata

package assets

import "embed"

// EmbedFS contains the minified API data and ASCII art for embedding
// This results in a much smaller binary (~25MB vs ~700MB with all assets)
//
//go:embed embed/api_data/*.json
//go:embed embed/art/*.ascii
var EmbedFS embed.FS
