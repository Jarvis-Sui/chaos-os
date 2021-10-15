package manager

import (
	"fmt"
	"os/exec"
	"path"

	"encoding/json"

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

	table := database.GetFaultTable()
	if err := table.AddFault(fault); err != nil {
		logrus.WithField("err", err).Errorf("failed to add an item to table %s", table.TableName)
	} else {
		if out, err := execute(fault); err == nil {
			logrus.WithFields(logrus.Fields{"out": out, "fault": fault}).Info("execute fault")
			if err := table.UpdateFaultStatus(fault.Uid, string(b.FS_RUNNING), out); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "table": table.TableName}).Error("failed to update table")
				return
			}
			// prepare to destroy
			binPath := util.GetBinPath()
			args := fmt.Sprintf("nohup /bin/sh -c 'sleep %d; %s fault destroy --id %s' > /dev/null 2>&1 &",
				timeout, binPath, fault.Uid)

			cmd := exec.Command("bash", "-c", args)
			if _, err := cmd.CombinedOutput(); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "out": out, "cmd": args}).Error("failed to execute command")
			}
		} else {
			table.UpdateFaultStatus(fault.Uid, string(b.FS_ERROR), fmt.Sprintf("%s. %s", out, err))
			logrus.WithFields(logrus.Fields{"err": err, "fault": fault}).Error("failed to execute fault")
		}
	}

}

func FaultDestroy(flags *pflag.FlagSet) {
	id, _ := flags.GetString("id")
	table := database.GetFaultTable()
	if fault, err := table.GetFaultById(id); err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "fault_id": id}).Error("failed to get fault")
	} else {
		logrus.WithField("fault", fault).Info("destroy fault")
	}
}

func FaultStatus(flags *pflag.FlagSet) {
	table := database.GetFaultTable()
	faults := table.GetAllFaults()

	for _, fault := range faults {
		s, _ := json.Marshal(fault)
		fmt.Printf("%s\n", s)
	}
}

func execute(fault *b.Fault) (string, error) {
	cc := exec.Command("bash", "-c", fault.Command)
	out, err := cc.CombinedOutput()
	return string(out), err
}
