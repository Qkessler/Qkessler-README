package main

import (
	"embed"

	"github.com/Qkessler/Qkessler-README/cmd"
)

//go:embed assets/static-description.md
var staticDescriptionContent embed.FS

func main() {
	cmd.Execute(staticDescriptionContent)
}
