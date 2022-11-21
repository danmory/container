package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "fork":
		fork()
	default:
		log.Fatal("bad command")
	}
}

func waitForNetwork() error {
	maxWait := time.Second * 3
	checkInterval := time.Second
	timeStarted := time.Now()
	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}

		if len(interfaces) > 1 {
			return nil
		}
		if time.Since(timeStarted) > maxWait {
			return fmt.Errorf("timeout after %s waiting for network", maxWait)
		}
		time.Sleep(checkInterval)
	}
}

func run() {
	fmt.Println("Starting container...")

	cmd := exec.Command("/proc/self/exe", append([]string{"fork"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET,
	}

	must(cmd.Start())

	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	netsetgoCmd := exec.Command("/usr/local/bin/netsetgo", "-pid", pid)
	must(netsetgoCmd.Run())
	must(cmd.Wait())
}

func fork() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(waitForNetwork())
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
