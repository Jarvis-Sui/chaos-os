package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jarvis-Sui/chaos-os/handler"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var serverCmd *cobra.Command

func initServerCmd() {
	serverCmd = &cobra.Command{
		Use: "server",
	}

	startCmd := &cobra.Command{
		Use: "start",
		Run: startServer,
	}

	var port int
	var nohup bool
	startCmd.Flags().IntVar(&port, "port", 9530, "server port")
	startCmd.Flags().BoolVar(&nohup, "background", false, "run in background")
	startCmd.MarkFlagRequired("port")

	stopCmd := &cobra.Command{
		Use: "stop",
		Run: stopServer,
	}

	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(stopCmd)
}

func startServer(cmd *cobra.Command, args []string) {
	if isServerRunning() {
		fmt.Println("server already running")
		return
	}

	port, _ := cmd.Flags().GetInt("port")
	nohup, _ := cmd.Flags().GetBool("background")
	if nohup {
		args := fmt.Sprintf("nohup %s server start --port %d >/dev/null 2>&1 &", util.GetBinPath(), port)
		cmd := exec.Command("bash", "-c", args)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("%s. %s\n", string(out), err)
		} else {
			fmt.Println("successfully started")
		}
		return
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/fault/status", handler.GetFaults)
	e.PUT("/fault", handler.AddFault)
	e.DELETE("/fault", handler.DestroyFault)
	logOutput, _ := os.OpenFile(filepath.Join(util.GetRootPath(), "echo.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	e.Logger.SetOutput(logOutput)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func stopServer(cmd *cobra.Command, args []string) {
	if !isServerRunning() {
		return
	}

	pid := os.Getpid()
	shArgs := fmt.Sprintf("ps -ef | grep -v grep | grep -v %d | grep 'chaos-os server start' | awk '{ print $2 }'", pid)
	shCmd := exec.Command("bash", "-c", shArgs)
	if out, err := shCmd.CombinedOutput(); err == nil {
		pid, _ := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64)
		if out, err = exec.Command("bash", "-c", fmt.Sprintf("kill -9 %d", pid)).CombinedOutput(); err != nil {
			fmt.Printf("%s. %s\n", string(out), err)
		} else {
			fmt.Println("successfully stopped")
		}
	}
}

func isServerRunning() bool {
	pid := os.Getpid()
	args := fmt.Sprintf("ps -ef | grep -v grep | grep -v %d | grep 'chaos-os server start' | wc -l", pid)
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err == nil {
		intv, _ := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64)
		return intv != 0
	} else {
		fmt.Printf("%s. %s\n", string(out), err)
		os.Exit(1)
	}
	return false
}
