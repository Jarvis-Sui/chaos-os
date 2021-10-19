package manager

import (
	"fmt"
	"strings"
	"time"

	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func createNetworkDelay(flags *pflag.FlagSet) *b.Fault {
	timeout, _ := flags.GetInt("timeout")
	commonArgs := fmt.Sprintf("%s create delay %s", netTcBinFile, buildNetCommonArgs(flags))

	args := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "delay" {
			args += fmt.Sprintf("--delay %v ", f.Value)
		} else if f.Name == "jitter" {
			args += fmt.Sprintf("--jitter %v ", f.Value)
		} else if f.Name == "correlation" {
			args += fmt.Sprintf("--correlation %v ", f.Value)
		} else if f.Name == "distribution" {
			args += fmt.Sprintf("--distribution %v ", f.Value)
		}
	})

	args = fmt.Sprintf("%s %s", commonArgs, args)
	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_NETDELAY, Status: b.FS_READY, Command: args, CreateTime: time.Now(), Timeout: timeout}
	return &fault
}

func createNetworkLoss(flags *pflag.FlagSet) *b.Fault {
	timeout, _ := flags.GetInt("timeout")
	commonArgs := fmt.Sprintf("%s create loss %s", netTcBinFile, buildNetCommonArgs(flags))
	args := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "percent" {
			args += fmt.Sprintf("--percent %v ", f.Value)
		} else if f.Name == "correlation" {
			args += fmt.Sprintf("--correlation %v ", f.Value)
		}
	})

	args = fmt.Sprintf("%s %s", commonArgs, args)
	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_NETLOSS, Status: b.FS_READY, Command: args, CreateTime: time.Now(), Timeout: timeout}
	return &fault
}

func buildNetCommonArgs(flags *pflag.FlagSet) string {
	args := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "interface" {
			args += fmt.Sprintf("--interface %v ", f.Value)
		} else if f.Name == "dest-ip" {
			args += fmt.Sprintf("--dest-ip %v ", f.Value)
		} else if f.Name == "dest-port" {
			args += fmt.Sprintf("--dest-port %v ", f.Value)
		}
	})
	return args
}

func getNetFaultInterface(fault *b.Fault) string {
	args := strings.Split(fault.Command, " ")
	for i, v := range args {
		if v == "--interface" {
			return args[i+1]
		}
	}
	return ""
}
