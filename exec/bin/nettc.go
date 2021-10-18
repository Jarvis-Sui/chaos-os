package main

import (
	"fmt"
	"os"

	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func destroyTc(cmd *cobra.Command, args []string) {
	classMinor, _ := cmd.Flags().GetInt("class-minor")

	execCmd(fmt.Sprintf("tc filter del dev bond0 parent 1: prio %d", classMinor))
	logrus.WithField("prio", classMinor).Info("deleted tc filter")

	execCmd(fmt.Sprintf("tc qdisc del dev bond0 handle %d: parent 1:%d", classMinor*10, classMinor))
	logrus.WithFields(logrus.Fields{"handle": classMinor * 10, "parent class": classMinor}).Info("deleted qdisc")

	execCmd(fmt.Sprintf("tc class del dev bond0 parent 1: classid 1:%d", classMinor))
	logrus.WithField("class minor", classMinor).Info("deleted tc class")
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

	createCmd := initCreateCmd()

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Run: destroyTc,
	}

	var classMinor int
	destroyCmd.Flags().IntVar(&classMinor, "class-minor", 0, "minor id of class")
	destroyCmd.MarkFlagRequired("class-minor")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)

	rootCmd.Execute()
}
