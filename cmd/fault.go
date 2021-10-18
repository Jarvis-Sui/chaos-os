package cmd

import (
	"github.com/Jarvis-Sui/chaos-os/manager"
	"github.com/spf13/cobra"
)

var faultCmd *cobra.Command
var faultType string
var timeout int

func initFaultCmd() {
	faultCmd = &cobra.Command{
		Use: "fault",
	}

	faultCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "create fault",
	}

	faultCreateCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "timeout")
	faultCreateCmd.MarkPersistentFlagRequired("timeout")

	initTcCmd()
	faultCreateCmd.AddCommand(tcCmd)

	var id string
	faultDestroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy fault by id",
		Run:   destroy,
	}
	faultDestroyCmd.Flags().StringVar(&id, "id", "", "fault id")
	faultDestroyCmd.MarkFlagRequired("id")

	var limit int
	var state string
	faultStatusCmd := &cobra.Command{
		Use: "status",
		Run: status,
	}
	faultStatusCmd.Flags().StringVar(&id, "id", "", "fault id")
	faultStatusCmd.Flags().StringVar(&state, "status", "", "status of faults to return. Ready | Running | Error | Destroyed")
	faultStatusCmd.Flags().IntVar(&limit, "limit", 100, "maximum number of faults returned")

	faultCmd.AddCommand(faultCreateCmd)
	faultCmd.AddCommand(faultDestroyCmd)
	faultCmd.AddCommand(faultStatusCmd)
}

func destroy(cmd *cobra.Command, args []string) {
	manager.FaultDestroy(cmd.Flags())
}

func status(cmd *cobra.Command, args []string) {
	manager.FaultStatus(cmd.Flags())
}
