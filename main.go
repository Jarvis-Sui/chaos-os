package main

import (
	"fmt"
	"os"

	"github.com/Jarvis-Sui/chaos-os/cmd"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/sirupsen/logrus"
)

func main() {
	dir := util.GetRootPath()
	f, err := os.OpenFile(dir+"/"+"chaosos.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", f)
	}

	defer f.Close()

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(f)
	cmd.Exec()
}
