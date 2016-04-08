package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	addr := net.TCPAddr{Port: 8080}
	ln, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Println("Fail to listen")
		return
	}

	file, err := ln.File()
	if err != nil {
		fmt.Println("Fail to get fd")
		ln.Close()
		return
	}
	cmd := exec.Command("echo_worker")
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles = append(cmd.ExtraFiles, file)
	cmd.Start()
	ln.Close()

	fmt.Println("Exec echo worker done")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGUSR1)

	for {
		select {
		case s := <-c:
			if s == syscall.SIGINT {
				goto out
			} else {
				fmt.Println("Signal:", s)
			}
		}
	}

out:
	err = cmd.Wait()
	fmt.Printf("Command finished: %v\n", err)
}
