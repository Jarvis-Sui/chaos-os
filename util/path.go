package util

import (
	"os"
	"path/filepath"
)

var bin = "chaos-os"

func GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if os.Args[0][len(os.Args[0])-len(bin):] != bin {
		dir = filepath.Dir(dir)
	}
	if err != nil {
		panic(err)
	}

	return dir
}

func GetDBFilePath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "chaos-os.dat")
}

func GetExecBinPath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "bin")
}

func GetBinPath() string {
	dir := GetRootPath()
	return filepath.Join(dir, bin)
}

func GetLogPath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "chaos-os.log")
}
