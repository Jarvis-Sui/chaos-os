package cmd

import (
	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/manager"
	"github.com/spf13/cobra"
)

var tcCmd *cobra.Command

func initTcCmd() {
	tcCmd = &cobra.Command{
		Use: "network",
	}

	var device, destIp, destPort string

	tcCmd.PersistentFlags().StringVarP(&device, "interface", "i", "bond0", "network interface")
	tcCmd.PersistentFlags().StringVar(&destIp, "dest-ip", "", "destination ip")
	tcCmd.PersistentFlags().StringVar(&destPort, "dest-port", "", "destination port")
	tcCmd.MarkPersistentFlagRequired("interface")
	tcCmd.MarkPersistentFlagRequired("dest-ip")
	tcCmd.MarkPersistentFlagRequired("dest-port")

	delayCmd := initTcDelayCmd()
	lossCmd := initTcLossCmd()

	tcCmd.AddCommand(delayCmd)
	tcCmd.AddCommand(lossCmd)
}

func initTcDelayCmd() *cobra.Command {
	delayCmd := &cobra.Command{
		Use: "delay",
		Run: func(cmd *cobra.Command, args []string) {
			manager.FaultCreate(binding.FT_NETDELAY, cmd.Flags())
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
			manager.FaultCreate(binding.FT_NETLOSS, cmd.Flags())
		},
	}

	var percent, correlation int
	lossCmd.Flags().IntVar(&percent, "percent", 0, "percent to loss. int value.")
	lossCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")

	lossCmd.MarkFlagRequired("percent")
	return lossCmd
}
