package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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
	checkTcExists()

	netem := netemFn[netemType](flags)
	device, _ := flags.GetString("interface")
	destIp, _ := flags.GetString("dest-ip")
	destPort, _ := flags.GetString("dest-port")

	createRootQdiscIfNotExist(device)
	classMinor := getNextClassMinor(device)

	addClass(device, classMinor)
	addSubQdisc(device, classMinor*10, classMinor, netem)
	addFilter(device, classMinor, classMinor, destIp, destPort)
	// class minor == filter prio == sub qdisc handle / 10, so caller can use it to find the newly added rule
	fmt.Println(classMinor)
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

func createRootQdiscIfNotExist(device string) {
	if !isRootQdiscExist(device) {
		args := fmt.Sprintf("tc qdisc add dev %s root handle 1: %s", device, qdisc)
		execCmd(args)
	}
}

func getNextClassMinor(device string) int {
	args := fmt.Sprintf("tc class ls dev %s | awk '{ print $3 }' | awk -F':' '{ print $2 }' | sort", device)
	out := execCmd(args)
	str := strings.Trim(out, "\n")
	if str == "" {
		return 1
	}
	existClasses := strings.Split(str, "\n")
	for i, v := range existClasses {
		intv, _ := strconv.ParseInt(v, 10, 64)
		if i+1 < int(intv) {
			return i + 1
		}
	}
	return len(existClasses) + 1
}

func addClass(device string, classMinor int) {
	args := fmt.Sprintf("tc class add dev %s classid 1:%d parent 1: htb rate 10000Mbps", device, classMinor)
	execCmd(args)
	logrus.WithFields(logrus.Fields{"interface": device, "class minor": classMinor}).Info("add new class")
}

func addSubQdisc(device string, handle int, parentClassMinor int, netem string) {
	args := fmt.Sprintf("tc qdisc add dev %s handle %d: parent 1:%d %s", device, handle, parentClassMinor, netem)
	execCmd(args)
	logrus.WithFields(
		logrus.Fields{"interface": device, "handle": handle, "class minor": parentClassMinor, "netem": netem},
	).Info("add new sub qdisc")
}

func addFilter(device string, flowId int, prio int, destIp string, destPort string) {
	args := fmt.Sprintf("tc filter add dev %s parent 1: prio %d protocol ip u32 match ip dst %s match ip dport %s 0xffff flowid 1:%d", device, prio, destIp, destPort, flowId)
	execCmd(args)
	logrus.WithFields(
		logrus.Fields{"interface": device, "flow id": flowId, "prio": prio, "dest ip": destIp, "dest port": destPort},
	).Info("add new filter")
}

func checkTcExists() {
	execCmd("which tc")
}
