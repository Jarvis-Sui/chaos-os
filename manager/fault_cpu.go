package manager

import (
	"fmt"
	"time"

	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func initCPUStress(flags *pflag.FlagSet) *b.Fault {
	timeout, _ := flags.GetInt("timeout")

	cpu, _ := flags.GetInt("cpu")
	load, _ := flags.GetInt("load")

	cpuMask, _ := flags.GetString("cpu-mask")
	taskset, _ := flags.GetString("taskset")

	args := fmt.Sprintf("%s create --cpu %d --load %d --timeout %d", cpuBinFile, cpu, load, timeout)

	if taskset != "" {
		args += fmt.Sprintf(" --taskset %s", taskset)
	} else if cpuMask != "" {
		args += fmt.Sprintf(" --cpu-mask %s", cpuMask)
	}

	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_CPUSTRESS, Status: b.FS_READY, Command: args, CreateTime: time.Now(), Timeout: timeout}
	return &fault
}
