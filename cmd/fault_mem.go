package cmd

import (
	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/spf13/cobra"
)

var memCmd *cobra.Command

func initMemCmd() {
	memCmd = &cobra.Command{
		Use:   "memory",
		Short: "memory fault",
	}

	stressCmd := &cobra.Command{
		Use:   "stress",
		Short: "stress memory",
		Run: func(cmd *cobra.Command, args []string) {
			addFault(binding.FT_MEMSTRESS, cmd.Flags())
		},
	}

	var nWorker int
	var bytes string
	stressCmd.Flags().IntVar(&nWorker, "worker-num", 1, "number of workers")
	stressCmd.Flags().StringVar(&bytes, "bytes", "128M", "the size of the POSIX shared memory objects to be created. One can specify the size as % of total available memory or in units of Bytes, KBytes, MBytes and GBytes using the suffix b, k, m or g.")
	memCmd.AddCommand(stressCmd)
}
