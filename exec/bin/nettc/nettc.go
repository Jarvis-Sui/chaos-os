package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var qdisc = "htb"

func destroyTc(cmd *cobra.Command, args []string) {
	classMinor, _ := cmd.Flags().GetInt("class-minor")
	device, _ := cmd.Flags().GetString("interface")

	execCmd(fmt.Sprintf("tc filter del dev %s parent 1: prio %d", device, classMinor))
	logrus.WithField("prio", classMinor).Info("deleted tc filter")

	execCmd(fmt.Sprintf("tc qdisc del dev %s handle %d: parent 1:%d", device, classMinor*10, classMinor))
	logrus.WithFields(logrus.Fields{"handle": classMinor * 10, "parent class": classMinor}).Info("deleted qdisc")

	execCmd(fmt.Sprintf("tc class del dev %s parent 1: classid 1:%d", device, classMinor))
	logrus.WithField("class minor", classMinor).Info("deleted tc class")

	if isAllFiltersDeleted(device) {
		execCmd(fmt.Sprintf("tc qdisc del dev %s root", device))
		logrus.Info("cleared all tc qdisc")
	}
}

func isRootQdiscExist(device string) bool {
	args := fmt.Sprintf("tc qdisc ls dev %s | grep %s | wc -l", device, qdisc)
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "%s. %s\n", out, err)
		os.Exit(1)
	} else {
		if v, _ := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64); v != 0 {
			return true
		}
	}
	return false
}

func isAllFiltersDeleted(device string) bool {
	args := fmt.Sprintf("tc filter ls dev %s | wc -l", device)
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "%s. %s\n", out, err)
		os.Exit(1)
	} else {
		if v, _ := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64); v == 0 {
			return true
		}
	}
	return false
}

func main() {
	f, err := os.OpenFile(util.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", f)
	}

	defer f.Close()
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(f)

	rootCmd := &cobra.Command{
		Use: "nettc",
	}

	createCmd := initTCCreateCmd()

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Run: destroyTc,
	}

	var classMinor int
	var device string
	destroyCmd.Flags().IntVar(&classMinor, "class-minor", 0, "minor id of class")
	destroyCmd.Flags().StringVarP(&device, "interface", "i", "", "interface name")
	destroyCmd.MarkFlagRequired("class-minor")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)

	rootCmd.Execute()
}
