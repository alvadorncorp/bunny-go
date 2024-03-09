package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {

	cmd := &cobra.Command{
		Use: "api",
	}

	cmd.AddCommand(storageCmd())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
