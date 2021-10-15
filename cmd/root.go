package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "chaos-os",
	}

	initFaultCmd()
	rootCmd.AddCommand(faultCmd)
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		errMsg, _ := fmt.Printf("failed to run cmd: %s", err)
		logrus.Error(errMsg)
	}
}
