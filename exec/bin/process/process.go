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
	"github.com/spf13/pflag"
)

func main() {
	f, err := os.OpenFile(util.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", f)
	}

	defer f.Close()
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(f)

	rootCmd := &cobra.Command{
		Use: "process",
	}

	rootCmd.AddCommand(initCreateCmd())
	rootCmd.AddCommand(initDestroyCmd())

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
		return nil
	})
	rootCmd.Execute()
}

func initCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use: "create",
		PreRun: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			pid, _ := flags.GetIntSlice("pid")
			pattern, _ := flags.GetString("pattern")
			if len(pid) == 0 && pattern == "" {
				fmt.Println("pid/pattern not set")
				os.Exit(1)
			}
		},
	}

	var pid []int
	var pattern string
	createCmd.PersistentFlags().IntSliceVar(&pid, "pid", []int{}, "process pid")
	createCmd.PersistentFlags().StringVar(&pattern, "pattern", "", "process pattern")

	createPauseCmd := &cobra.Command{
		Use: "pause",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd.Flags(), "STOP")
		},
	}

	createCmd.AddCommand(createPauseCmd)
	return createCmd

}

func initDestroyCmd() *cobra.Command {
	destroyCmd := &cobra.Command{
		Use: "destroy",
		PreRun: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			pid, _ := flags.GetIntSlice("pid")
			pattern, _ := flags.GetString("pattern")
			if len(pid) == 0 && pattern == "" {
				fmt.Println("pid/pattern not set")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd.Flags(), "CONT")
		},
	}

	var pid []int
	var pattern string
	destroyCmd.PersistentFlags().IntSliceVar(&pid, "pid", []int{}, "process pid")
	destroyCmd.PersistentFlags().StringVar(&pattern, "pattern", "", "process pattern")

	return destroyCmd

}

func run(flags *pflag.FlagSet, signal string) {
	pids, _ := flags.GetIntSlice("pid")

	if len(pids) == 0 {
		pattern, _ := flags.GetString("pattern")
		pids = getPids(pattern)
	}

	strPids := []string{}
	for _, p := range pids {
		strPids = append(strPids, fmt.Sprintf("%d", p))
	}

	logrus.WithFields(logrus.Fields{"pids": pids, "signal": signal}).Info("kill process")

	execCmd(fmt.Sprintf("arr=(%s) && for i in ${arr[@]}; do kill -%s $i; done", strings.Join(strPids, " "), signal))
	fmt.Println(strings.Join(strPids, ","))
}

func getPids(pattern string) []int {
	pid := os.Getpid()
	ppid := os.Getppid()
	out := execCmd(fmt.Sprintf("ps -ef | grep -v grep | grep -v %d | grep -v %d | grep -- '%s' | awk '{ print $2 }'", pid, ppid, pattern))

	pids := []int{}
	for _, v := range strings.Split(out, "\n") {
		if i, err := strconv.ParseInt(strings.TrimRight(v, "\n"), 10, 64); err == nil {
			pids = append(pids, int(i))
		}
	}

	return pids
}

func execCmd(args string) string {
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err != nil {
		logrus.WithFields(logrus.Fields{"out": string(out), "err": err, "cmd": args}).Error("failed to exec command")
		os.Exit(1)
	} else {
		return string(out)
	}

	return ""
}
