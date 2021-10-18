package main

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func execCmd(args string) string {
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err != nil {
		logrus.WithFields(logrus.Fields{"out": string(out), "err": err, "cmd": args}).Error("failed to exec command")
		os.Exit(1)
	} else {
		return string(out)
	}

	return ""
}
