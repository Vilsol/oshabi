package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var lineRegex = regexp.MustCompile(`(?m)^\s+(.+?)\s=>\s(.+?)\s\(`)

func main() {
	flag.Parse()
	cmd := exec.Command("ldd", flag.Arg(0))

	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	pathEnv := os.Getenv("PATH")
	pathSplit := strings.Split(pathEnv, ";")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	outputDir := path.Join(wd, "out")

	err = os.MkdirAll(outputDir, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	matches := lineRegex.FindAllSubmatch(stdout, -1)
	for _, match := range matches {
		println("Found match:", string(match[0]))
		for _, p := range pathSplit {
			src := path.Join(p, string(match[1]))

			if _, err := os.Stat(src); err != nil {
				continue
			}

			dst := path.Join(outputDir, string(match[1]))
			println("Copying", src, "=>", dst)
			if _, err := copyFile(src, dst); err != nil {
				panic(err)
			}
		}
	}
}

func copyFile(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, errors.Wrap(err, "failed to stat source file")
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, errors.Wrap(err, "failed to open source file for reading")
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create destination file")
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, errors.Wrap(err, "failed to copy file")
}
