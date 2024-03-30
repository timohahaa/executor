package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("ls", "-la")
	cmd2 := exec.Command("ls", "-la")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: false,
	}
	cmd2.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: false,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	err = cmd2.Start()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("go process pid: %d\n", os.Getpid())
	fmt.Printf("command pid: %d\n", cmd.Process.Pid)
	fmt.Printf("command 2 pid: %d\n", cmd2.Process.Pid)

	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}
	err = cmd2.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
