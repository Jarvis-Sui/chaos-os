package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var netemFn = map[string]func(*pflag.FlagSet) string{
	"delay": netemDelay,
	"loss":  netemLoss,
}

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
			createTcRule(cmd.Flags(), "delay")
		},
	}

	var delay, jitter, correlation int
	var distribution string
	delayCmd.Flags().IntVar(&delay, "delay", 0, "ms time to delay")
	delayCmd.Flags().IntVar(&jitter, "jitter", 0, "ms time to jitter")
	delayCmd.Flags().IntVar(&correlation, "correlation", 0, "%% correlation between packets, range [0-100]")
	delayCmd.Flags().StringVar(&distribution, "distribution", "normal", "destribution of delay, uniform | normal | pareto |  paretonormal")

	delayCmd.MarkFlagRequired("delay")
	return delayCmd
}

func initTcLossCmd() *cobra.Command {
	lossCmd := &cobra.Command{
		Use:   "loss",
		Short: "random loss of packets",
		Run: func(cmd *cobra.Command, args []string) {
			createTcRule(cmd.Flags(), "loss")
		},
	}

	var percent, correlation int
	lossCmd.Flags().IntVar(&percent, "percent", 0, "percent to loss. range [0-100]")
	lossCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")

	lossCmd.MarkFlagRequired("percent")
	return lossCmd
}

func createTcRule(flags *pflag.FlagSet, netemType string) {
	// 0. build netem parameters
	// 1. check root qdisc exists
	// 2. get avail next classid
	// 3. add class, sub qdisc, filter
	//
	netem := netemFn[netemType](flags)
	fmt.Println(netem)
}

func netemDelay(flags *pflag.FlagSet) string {
	delay, _ := flags.GetInt("delay")
	jitter, _ := flags.GetInt("jitter")
	correlation, _ := flags.GetInt("correlation")
	distribution, _ := flags.GetString("distribution")
	return fmt.Sprintf("netem delay %dms %dms %d%% distribution %s", delay, jitter, correlation, distribution)
}

func netemLoss(flags *pflag.FlagSet) string {
	percent, _ := flags.GetInt("percent")
	correlation, _ := flags.GetInt("correlation")
	return fmt.Sprintf("netem loss random %d%% %d%%", percent, correlation)
}
