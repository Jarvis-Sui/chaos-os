package manager

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	b "github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/database"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type faultfn func(*pflag.FlagSet) *b.Fault

var faultInitFns map[b.FaultType]faultfn
var netTcBinFile = path.Join(util.GetExecBinPath(), "nettc")
var procBinFile = path.Join(util.GetExecBinPath(), "process")
var cpuBinFile = path.Join(util.GetExecBinPath(), "cpu")
var memBinFile = path.Join(util.GetExecBinPath(), "memory")

func init() {
	faultInitFns = map[b.FaultType]faultfn{
		b.FT_NETLOSS:      initNetworkLoss,
		b.FT_NETDELAY:     initNetworkDelay,
		b.FT_NETREORDER:   initNetworkReorder,
		b.FT_NETDUPLICATE: initNetworkDuplicate,
		b.FT_NETCORRUPT:   initNetworkCorrupt,
		b.FT_PROCPAUSE:    initProcessPause,
		b.FT_CPUSTRESS:    initCPUStress,
		b.FT_MEMSTRESS:    initMemStress,
	}

}

func InitFault(ft b.FaultType, flags *pflag.FlagSet) (*b.Fault, error) {
	_, err := flags.GetInt("timeout")
	if err != nil {
		logrus.Error("timeout parameter not set")
		return nil, err
	}
	fault := faultInitFns[ft](flags)
	return fault, nil
}

func CreateFault(fault *b.Fault) error {
	table := database.GetFaultTable()
	if err := table.AddFault(fault); err != nil {
		logrus.WithField("err", err).Errorf("failed to add an item to table %s", table.TableName)
		return err
	} else {
		if out, err := execute(fault); err == nil {
			logrus.WithFields(logrus.Fields{"out": out, "fault": fault}).Info("execute fault")
			if err := table.UpdateFaultStatus(fault.Uid, string(b.FS_RUNNING), out); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "table": table.TableName}).Error("failed to update table")
				return err
			}
			// prepare to destroy
			binPath := util.GetBinPath()
			args := fmt.Sprintf("nohup /bin/sh -c 'sleep %d; %s fault destroy --id %s' > /dev/null 2>&1 &",
				fault.Timeout, binPath, fault.Uid)

			cmd := exec.Command("bash", "-c", args)
			if _, err := cmd.CombinedOutput(); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "out": out, "cmd": args}).Error("failed to execute command")
				return err
			}
		} else {
			table.UpdateFaultStatus(fault.Uid, string(b.FS_ERROR), fmt.Sprintf("%s. %s", out, err))
			logrus.WithFields(logrus.Fields{"out": out, "err": err, "fault": fault}).Error("failed to execute fault")
			return err
		}
	}

	return nil

}

func DestroyFault(flags *pflag.FlagSet) error {
	id, _ := flags.GetString("id")
	table := database.GetFaultTable()
	if fault, err := table.GetFaultById(id); err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "fault_id": id}).Error("failed to get fault")
		return err
	} else {
		if fault.Status == b.FS_DESTROYED {
			return nil
		}
		logrus.WithField("fault", fault).Info("destroy fault")
		var args string
		if fault.Type == b.FT_NETLOSS || fault.Type == b.FT_NETDELAY || fault.Type == b.FT_NETREORDER || fault.Type == b.FT_NETDUPLICATE || fault.Type == b.FT_NETCORRUPT {
			classMinor := fault.Reason
			device := getNetFaultInterface(fault)
			args = fmt.Sprintf("%s destroy --class-minor %s --interface %s", netTcBinFile, classMinor, device)
		} else if fault.Type == b.FT_PROCPAUSE {
			pids := fault.Reason
			args = fmt.Sprintf("%s destroy --pid %s", procBinFile, pids)
		} else if fault.Type == b.FT_CPUSTRESS || fault.Type == b.FT_MEMSTRESS {
			pid := fault.Reason
			args = fmt.Sprintf("%s destroy --pid %s", cpuBinFile, pid)
		} else {
			logrus.WithField("type", fault.Type).Error("fault type not supported")
			return fmt.Errorf("fault type %s not supported", fault.Type)
		}

		cmd := exec.Command("bash", "-c", args)
		if out, err := cmd.CombinedOutput(); err != nil {
			table.UpdateFaultStatus(fault.Uid, string(b.FS_ERROR), fmt.Sprintf("%s. %s", out, err))
			logrus.WithFields(logrus.Fields{"err": err, "out": out, "fault": fault}).Error("failed to destroy fault")
			return err
		} else {
			if err := table.UpdateFaultStatus(fault.Uid, string(b.FS_DESTROYED), ""); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "table": table.TableName}).Error("failed to update table")
				return err
			}
		}
	}
	return nil
}

func FaultStatus(flags *pflag.FlagSet) []*b.Fault {
	status, _ := flags.GetString("status")
	id, _ := flags.GetString("id")
	limit, _ := flags.GetInt("limit")
	table := database.GetFaultTable()
	faults := table.GetFaults(id, b.FaultStatus(status), limit)
	return faults
}

func execute(fault *b.Fault) (string, error) {
	cc := exec.Command("bash", "-c", fault.Command)
	out, err := cc.CombinedOutput()
	return strings.Trim(string(out), "\n"), err
}
