package cmd

import (
	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/spf13/cobra"
)

var cpuCmd *cobra.Command

func initCpuCmd() {
	cpuCmd = &cobra.Command{
		Use:   "cpu",
		Short: "cpu fault",
	}

	stressCmd := &cobra.Command{
		Use:   "stress",
		Short: "stress cpu",
		Run: func(cmd *cobra.Command, args []string) {
			addFault(binding.FT_CPUSTRESS, cmd.Flags())
		},
	}

	var nCpu, load int
	var taskset, cpuMask string
	stressCmd.Flags().IntVar(&nCpu, "cpu", -1, "number of workers. -1: all cpus")
	stressCmd.Flags().IntVar(&load, "load", 100, "load for a single core. default 100")
	stressCmd.Flags().StringVar(&taskset, "taskset", "", "taskset. (0 to N-1)")
	stressCmd.Flags().StringVar(&cpuMask, "cpu-mask", "", "cpu mask. 0: all cpus.")

	cpuCmd.AddCommand(stressCmd)
}
