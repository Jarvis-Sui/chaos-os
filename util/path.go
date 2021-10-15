package util

import (
	"os"
	"path/filepath"
)

func GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	return dir
}

func GetDBFilePath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "chaosos.dat")
}

func GetExecBinPath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "bin")
}

func GetBinPath() string {
	dir := GetRootPath()
	return filepath.Join(dir, "chaos-os")
}
