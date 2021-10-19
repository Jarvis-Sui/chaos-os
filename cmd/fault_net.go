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
	reorderCmd := initTcReorderCmd()
	dupCmd := initTcDuplicateCmd()
	corruptCmd := initTcCorruptCmd()

	tcCmd.AddCommand(delayCmd)
	tcCmd.AddCommand(lossCmd)
	tcCmd.AddCommand(reorderCmd)
	tcCmd.AddCommand(dupCmd)
	tcCmd.AddCommand(corruptCmd)
}

func initTcDelayCmd() *cobra.Command {
	delayCmd := &cobra.Command{
		Use: "delay",
		Run: func(cmd *cobra.Command, args []string) {
			addNetFault(binding.FT_NETDELAY, cmd.Flags())
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
			addNetFault(binding.FT_NETLOSS, cmd.Flags())
		},
	}

	var percent, correlation int
	lossCmd.Flags().IntVar(&percent, "percent", 0, "percent to loss. int value.")
	lossCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")

	lossCmd.MarkFlagRequired("percent")
	return lossCmd
}

func initTcReorderCmd() *cobra.Command {
	reorderCmd := &cobra.Command{
		Use:   "reorder",
		Short: "reorder of packets",
		Run: func(cmd *cobra.Command, args []string) {
			addNetFault(binding.FT_NETREORDER, cmd.Flags())
		},
	}

	var delay, percent, correlation, distance int
	reorderCmd.Flags().IntVar(&delay, "delay", 0, "ms time to delay")
	reorderCmd.Flags().IntVar(&percent, "percent", 0, "percent to reorder. int value.")
	reorderCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")
	reorderCmd.Flags().IntVar(&distance, "distance", 0, "gap")

	reorderCmd.MarkFlagRequired("delay")
	reorderCmd.MarkFlagRequired("percent")

	return reorderCmd
}

func initTcDuplicateCmd() *cobra.Command {
	dupCmd := &cobra.Command{
		Use:   "duplicate",
		Short: "duplication of packets",
		Run: func(cmd *cobra.Command, args []string) {
			addNetFault(binding.FT_NETDUPLICATE, cmd.Flags())
		},
	}
	var percent, correlation int
	dupCmd.Flags().IntVar(&percent, "percent", 0, "percent to reorder. int value.")
	dupCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")
	dupCmd.MarkFlagRequired("percent")

	return dupCmd
}

func initTcCorruptCmd() *cobra.Command {
	corruptCmd := &cobra.Command{
		Use:   "corrupt",
		Short: "corruption of packets",
		Run: func(cmd *cobra.Command, args []string) {
			addNetFault(binding.FT_NETCORRUPT, cmd.Flags())
		},
	}
	var percent, correlation int
	corruptCmd.Flags().IntVar(&percent, "percent", 0, "percent to reorder. int value.")
	corruptCmd.Flags().IntVar(&correlation, "correlation", 0, "correlation between packets")
	corruptCmd.MarkFlagRequired("percent")

	return corruptCmd
}

func addNetFault(ft binding.FaultType, flags *pflag.FlagSet) {
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
