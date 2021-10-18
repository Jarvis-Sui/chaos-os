package main

import (
	"fmt"
	"os"
	"os/exec"
)

func execCmd(args string) string {
	cmd := exec.Command("bash", "-c", args)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "%s. %s\n", out, err)
		os.Exit(1)
	} else {
		return string(out)
	}

	return ""
}
