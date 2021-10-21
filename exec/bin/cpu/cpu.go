package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stressng = filepath.Join(util.GetExecBinPath(), "stress-ng")

func main() {
	f, err := os.OpenFile(util.GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", f)
	}

	defer f.Close()
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(f)

	rootCmd := &cobra.Command{
		Use: "cpu",
	}

	createCmd := &cobra.Command{
		Use: "create",
		Run: create,
	}

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Run: func(cmd *cobra.Command, args []string) {
			pid, _ := cmd.Flags().GetInt("pid")
			cmdArgs := fmt.Sprintf(`arr=$(ps -o pid --ppid %d | grep -v PID) && if [ ! -z "${arr}" ]; then kill -9 $arr; fi`, pid)
			exec.Command("bash", "-c", cmdArgs).Run()
		},
	}

	var nCpu, timeout, load int
	var taskset, cpuMask string
	createCmd.Flags().IntVar(&nCpu, "cpu", -1, "number of workers. -1: all cpus")
	createCmd.Flags().IntVar(&load, "load", 100, "load for a single core. default 100")
	createCmd.Flags().IntVar(&timeout, "timeout", 0, "timeout. 0: forever")
	createCmd.Flags().StringVar(&taskset, "taskset", "", "taskset. (0 to N-1)")
	createCmd.Flags().StringVar(&cpuMask, "cpu-mask", "", "cpu mask. 0: all cpus.")
	rootCmd.AddCommand(createCmd)

	var pid int
	destroyCmd.Flags().IntVar(&pid, "pid", -1, "pid")
	rootCmd.AddCommand(destroyCmd)

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
		return nil
	})
	rootCmd.Execute()
}

func create(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()

	nCpu, _ := flags.GetInt("cpu")
	load, _ := flags.GetInt("load")
	timeout, _ := flags.GetInt("timeout")

	cmdArgs := fmt.Sprintf("%s --cpu %d -l %d -t %d", stressng, nCpu, load, timeout)

	cpuMask, _ := flags.GetString("cpu-mask")
	taskset, _ := flags.GetString("taskset")

	if taskset != "" {
		cmdArgs += fmt.Sprintf(" --taskset %s", taskset)
	} else if cpuMask != "" {
		cpuList := maskToCpuList(cpuMask)
		cpuListStr := []string{}
		for _, v := range cpuList {
			cpuListStr = append(cpuListStr, fmt.Sprintf("%d", v))
		}
		cmdArgs += fmt.Sprintf(" --taskset %s", strings.Join(cpuListStr, ","))
	}

	now := time.Now().UnixNano()
	cmdArgs = fmt.Sprintf("nohup %s >%d 2>&1 &", cmdArgs, now)

	if out, err := exec.Command("bash", "-c", cmdArgs).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "%s. %s\n", string(out), err)
		os.Exit(1)
	} else {
		pid := getStressNgPid(now)
		fmt.Println(pid)
	}
}

func maskToCpuList(mask string) []int {
	mask = strings.ToLower(mask)
	mask = strings.Replace(mask, "0x", "", -1)

	bitStr := ""
	for _, v := range mask {
		intV, err := strconv.ParseInt(string(v), 16, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse mask: %s\n", err)
		}
		bitStr += strconv.FormatInt(intV, 2)
	}
	n := len(bitStr)
	ret := []int{}

	for i, v := range bitStr {
		if v == '1' {
			ret = append(ret, n-1-i)
		}
	}
	return ret
}

func getStressNgPid(file int64) int64 {
	var out []byte
	for {
		out, _ = exec.Command("bash", "-c", fmt.Sprintf("cat %d", file)).CombinedOutput()
		if len(strings.Split(string(out), "\n")) >= 2 {
			break
		}
		time.Sleep(100 * time.Microsecond)
	}
	exec.Command("bash", "-c", fmt.Sprintf("rm -f %d", file)).Run()
	firstLine := strings.Split(string(out), "\n")[0]
	r := regexp.MustCompile(`stress-ng: info:  \[(?P<pid>\d+)\]`)

	matches := r.FindStringSubmatch(firstLine)

	if v, err := strconv.ParseInt(matches[r.SubexpIndex("pid")], 10, 64); err == nil {
		return v
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return -1
}
