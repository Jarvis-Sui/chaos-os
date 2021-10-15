package manager

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/Jarvis-Sui/chaos-os/binding"
	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/database"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type faultfn func(*pflag.FlagSet) *b.Fault

var faultCreateFns map[binding.FaultType]faultfn
var binFile = path.Join(util.GetExecBinPath(), "nettc")

func init() {
	faultCreateFns = map[binding.FaultType]faultfn{
		binding.FT_NETLOSS:  createNetworkLoss,
		binding.FT_NETDELAY: createNetworkDelay,
	}

}

func FaultCreate(ft binding.FaultType, flags *pflag.FlagSet) {
	timeout, err := flags.GetInt("timeout")
	if err != nil {
		logrus.Error("timeout parameter not set")
		return
	}

	fault := faultCreateFns[ft](flags)
	fault.Status = b.FS_READY

	table := database.GetFaultTable()
	if err := table.AddFault(fault); err != nil {
		logrus.WithField("err", err).Errorf("failed to add an item to table %s", table.TableName)
	} else {
		if out, err := execute(fault); err == nil {
			logrus.WithFields(logrus.Fields{"out": out, "cmd": fault.Command}).Info("execute command")
			table.UpdateFaultStatus(fault.Uid, string(b.FS_RUNNING), out)

			// prepare to destroy
			binPath := util.GetBinPath()
			args := fmt.Sprintf("nohup /bin/sh -c 'sleep %d; %s fault destroy --id %s' > /dev/null 2>&1 &",
				timeout, binPath, fault.Uid)

			cmd := exec.Command("bash", "-c", args)
			if _, err := cmd.CombinedOutput(); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "out": out, "cmd": args}).Error("failed to execute command")
			} else {
				logrus.WithField("cmd", args).Info("execute command")
			}
		} else {
			table.UpdateFaultStatus(fault.Uid, string(b.FS_ERROR), out)
			logrus.WithFields(logrus.Fields{"err": err, "cmd": fault.Command}).Error("failed to execute command")
		}
	}

}

func FaultDestroy(flags *pflag.FlagSet) {
	id, _ := flags.GetString("id")
	logrus.WithField("id", id).Info("destroying fault")
	// table := database.GetFaultTable()
}

func FaultStatus(flags *pflag.FlagSet) {

}

func execute(fault *b.Fault) (string, error) {
	cc := exec.Command("bash", "-c", fault.Command)

	out, err := cc.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "out": out, "cmd": cc}).Error("failed to execute command")
	} else {
		logrus.WithField("cmd", cc).Info("execute command")
	}
	return string(out), err
}
