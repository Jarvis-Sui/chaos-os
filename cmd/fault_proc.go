package cmd

import (
	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/spf13/cobra"
)

var procCmd *cobra.Command

func initProcCmd() {
	procCmd = &cobra.Command{
		Use:   "process",
		Short: "process fault",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "pause the process",
		Run: func(cmd *cobra.Command, args []string) {
			addFault(binding.FT_PROCPAUSE, cmd.Flags())
		},
	}

	var pid []int
	var pattern string
	procCmd.PersistentFlags().IntSliceVar(&pid, "pid", []int{}, "process pid")
	procCmd.PersistentFlags().StringVar(&pattern, "pattern", "", "process pattern, used by grep")

	procCmd.AddCommand(pauseCmd)
}
