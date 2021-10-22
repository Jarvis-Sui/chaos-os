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
	"github.com/spf13/cobra"
)

var stressng = filepath.Join(util.GetExecBinPath(), "stress-ng")

func main() {
	rootCmd := &cobra.Command{
		Use: "memory",
	}

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Run: func(cmd *cobra.Command, args []string) {
			pid, _ := cmd.Flags().GetInt("pid")
			cmdArgs := fmt.Sprintf(`arr=$(ps -o pid --ppid %d | grep -v PID) && if [ ! -z "${arr}" ]; then kill -9 $arr; fi`, pid)
			exec.Command("bash", "-c", cmdArgs).Run()
		},
	}

	rootCmd.AddCommand(initCreateCmd())

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

func initCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use: "create",
	}

	var timeout int
	createCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "timeout")

	vmCmd := &cobra.Command{
		Use:   "vm",
		Short: "stress on memory. may take some time to consume the specified amount of memory",
		Run:   createVm,
	}

	var nWorker int
	var bytes string
	vmCmd.Flags().IntVar(&nWorker, "worker-num", 1, "number of workers")
	vmCmd.Flags().StringVar(&bytes, "bytes", "128M", "the size of the POSIX shared memory objects to be created. One can specify the size as % of total available memory or in units of Bytes, KBytes, MBytes and GBytes using the suffix b, k, m or g.")

	createCmd.AddCommand(vmCmd)
	return createCmd
}

func createVm(cmd *cobra.Command, args []string) {
	nWorker, _ := cmd.Flags().GetInt("worker-num")
	bytes, _ := cmd.Flags().GetString("bytes")
	timeout, _ := cmd.Flags().GetInt("timeout")

	cmdArgs := fmt.Sprintf("%s --vm %d --vm-bytes %s -t %d", stressng, nWorker, bytes, timeout)

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
