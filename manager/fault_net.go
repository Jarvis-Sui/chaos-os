package manager

import (
	"fmt"

	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func createNetworkDelay(flags *pflag.FlagSet) *b.Fault {
	commonArgs := fmt.Sprintf("%s create delay %s", binFile, buildNetCommonArgs(flags))

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
	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_NETDELAY, Status: b.FS_RUNNING, Command: args}
	return &fault
}

func createNetworkLoss(flags *pflag.FlagSet) *b.Fault {
	commonArgs := fmt.Sprintf("%s create loss %s", binFile, buildNetCommonArgs(flags))
	args := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "percent" {
			args += fmt.Sprintf("--percent %v ", f.Value)
		} else if f.Name == "correlation" {
			args += fmt.Sprintf("--correlation %v ", f.Value)
		}
	})

	args = fmt.Sprintf("%s %s", commonArgs, args)
	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_NETDELAY, Status: b.FS_RUNNING, Command: args}
	return &fault
	// uid := uuid.NewString()
	// table := database.GetFaultTable()
	// if err := table.AddFault(&fault); err != nil {
	// 	logrus.WithField("err", err).Errorf("failed to add an item to table %s", table.TableName)
	// 	return uid, err
	// }
	// cmd := fmt.Sprintf("%s %s", commonArgs, args)
	// if err := execute(cmd); err != nil {
	// 	table.UpdateFaultStatus(uid, string(b.FS_ERROR), fmt.Sprintf("%s", err))
	// }
	// return uid, nil
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
