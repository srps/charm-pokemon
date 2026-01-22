//go:build !realdata

package assets

import "embed"

// EmbedFS is empty when realdata is not present
var EmbedFS embed.FS
