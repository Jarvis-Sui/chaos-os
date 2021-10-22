package manager

import (
	"fmt"
	"time"

	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
)

func initMemStress(flags *pflag.FlagSet) *b.Fault {
	timeout, _ := flags.GetInt("timeout")

	nWorker, _ := flags.GetInt("worker-num")
	bytes, _ := flags.GetString("bytes")

	args := fmt.Sprintf("%s create vm %d --bytes %s --timeout %d", memBinFile, nWorker, bytes, timeout)

	fault := b.Fault{Uid: uuid.NewString(), Type: b.FT_MEMSTRESS, Status: b.FS_READY, Command: args, CreateTime: time.Now(), Timeout: timeout}
	return &fault
}
