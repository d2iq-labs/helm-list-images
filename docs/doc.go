package main

import (
	"log"

	"github.com/d2iq-labs/helm-list-images/cmd"
	"github.com/spf13/cobra/doc"
)

//go:generate go run github.com/d2iq-labs/helm-list-images/docs
func main() {
	commands := cmd.SetListImagesCommands()

	if err := doc.GenMarkdownTree(commands, "doc"); err != nil {
		log.Fatal(err)
	}
}
