package main

import (
	"log"
	"os"
	"os/exec"
)

func crossBuildStart() {
	err := os.Remove("/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Link("/usr/bin/resin-xbuild", "/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
}

func crossBuildEnd() {
	err := os.Remove("/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Link("/bin/sh.real", "/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
}

func runShell() {
	cmd := exec.Command("/usr/bin/qemu-arm-static", append([]string{"-0", "/bin/sh", "/bin/sh"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	switch os.Args[0] {
	case "cross-build-start":
		crossBuildStart()
	case "cross-build-end":
		crossBuildEnd()
	case "/bin/sh":
		crossBuildEnd()
		runShell()
		crossBuildStart()
	}
}
