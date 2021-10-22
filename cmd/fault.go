package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	initProcCmd()
	initCpuCmd()
	initMemCmd()

	faultCreateCmd.AddCommand(tcCmd)
	faultCreateCmd.AddCommand(procCmd)
	faultCreateCmd.AddCommand(cpuCmd)
	faultCreateCmd.AddCommand(memCmd)

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
	if err := manager.DestroyFault(cmd.Flags()); err != nil {
		fmt.Printf("failed to destroy: %s\n", err)
	}
}

func status(cmd *cobra.Command, args []string) {
	faults := manager.FaultStatus(cmd.Flags())
	for _, fault := range faults {
		s, _ := json.Marshal(fault)
		fmt.Printf("%s\n", s)
	}
}

func addFault(ft binding.FaultType, flags *pflag.FlagSet) {
	if fault, err := manager.InitFault(ft, flags); err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "fault": fault}).Error("failed to add fault")
		fmt.Printf("failed to add: %s\n", err)
	} else {
		if err := manager.CreateFault(fault); err != nil {
			logrus.WithFields(logrus.Fields{"err": err, "fault": fault}).Error("failed to add fault")
			fmt.Printf("failed to add: %s\n", err)
		} else {
			s, _ := json.Marshal(fault)
			fmt.Println(string(s))
		}
	}
}
