package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/Jarvis-Sui/chaos-os/binding"
	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func initProcessPause(flags *pflag.FlagSet) *binding.Fault {
	timeout, _ := flags.GetInt("timeout")
	commonArgs := fmt.Sprintf("%s create pause %s", procBinFile, buildProcCommonArgs(flags))

	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_PROCPAUSE, Status: b.FS_READY, Command: commonArgs, CreateTime: time.Now(), Timeout: timeout}
	return &fault
}

func buildProcCommonArgs(flags *pflag.FlagSet) string {
	args := ""
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "pid" {
			pids, _ := flags.GetIntSlice("pid")
			if len(pids) != 0 {
				strPids := []string{}
				for _, p := range pids {
					strPids = append(strPids, fmt.Sprintf("%d", p))
				}
				args += fmt.Sprintf("--pid %s", strings.Join(strPids, ","))
			}
		} else if f.Name == "pattern" {
			args += fmt.Sprintf("--pattern %v ", f.Value)
		}
	})

	return args
}
