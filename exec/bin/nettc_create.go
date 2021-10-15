package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use: "create",
	}

	var device, destIp, destPort string
	createCmd.PersistentFlags().StringVarP(&device, "interface", "i", "bond0", "network interface")
	createCmd.PersistentFlags().StringVar(&destIp, "dest-ip", "", "destination ip")
	createCmd.PersistentFlags().StringVar(&destPort, "dest-port", "", "destination port")
	createCmd.MarkPersistentFlagRequired("interface")
	createCmd.MarkPersistentFlagRequired("dest-ip")
	createCmd.MarkPersistentFlagRequired("dest-port")

	delayCmd := initTcDelayCmd()
	lossCmd := initTcLossCmd()
	createCmd.AddCommand(delayCmd)
	createCmd.AddCommand(lossCmd)
	return createCmd
}

func initTcDelayCmd() *cobra.Command {
	delayCmd := &cobra.Command{
		Use: "delay",
		Run: func(cmd *cobra.Command, args []string) {
			createTcRule(1, 1, 1, "args")
		},
	}

	var delay, jitter, correlation int
	var distribution string
	delayCmd.Flags().IntVar(&delay, "delay", 0, "ms time to delay")
	delayCmd.Flags().IntVar(&jitter, "jitter", 0, "ms time to jitter")
	delayCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")
	delayCmd.Flags().StringVar(&distribution, "distribution", "normal", "destribution of delay")

	delayCmd.MarkFlagRequired("delay")
	return delayCmd
}

func initTcLossCmd() *cobra.Command {
	lossCmd := &cobra.Command{
		Use:   "loss",
		Short: "random loss of packets",
		Run: func(cmd *cobra.Command, args []string) {
			createTcRule(0, 0, 0, "args")
		},
	}

	var percent, correlation int
	lossCmd.Flags().IntVar(&percent, "percent", 0, "percent to loss. int value.")
	lossCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")

	lossCmd.MarkFlagRequired("percent")
	return lossCmd
}

func createTcRule(qdiscId int, classId int, prio int, netem string) {
	fmt.Printf("qdisc %d, class %d, prio %d, netem %s\n", qdiscId, classId, prio, netem)
}
