package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func destroyTc(cmd *cobra.Command, args []string) {
	fmt.Println("destroyTc")
}

func listTc(cmd *cobra.Command, args []string) {
	fmt.Println("listTc")
}

func main() {
	rootCmd := &cobra.Command{
		Use: "nettc",
	}

	createCmd := initCreateCmd()

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Run: destroyTc,
	}

	listCmd := &cobra.Command{
		Use: "ls",
		Run: listTc,
	}

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(listCmd)

	rootCmd.Execute()
}
