package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/d2iq-labs/helm-list-images/cmd"
)

//go:generate go run github.com/d2iq-labs/helm-list-images/docs
func main() {
	commands := cmd.GetRootCommand()

	if err := doc.GenMarkdownTree(commands, "doc"); err != nil {
		log.Fatal(err)
	}
}
