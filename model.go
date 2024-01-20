package main

import (
	"fmt"
	"os"
	"os/exec"
)

type HostConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Key      string `json:"key"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}

func (h *HostConfig) Dail() error {
	host := h.Host
	port := h.Port
	user := h.User
	pass := h.Pass
	key := h.Key
	timeout := h.Timeout
	interval := h.Interval

	if timeout <= 0 {
		timeout = 10
	}
	if interval <= 0 {
		interval = 60
	}

	var cmdPath string
	var args []string
	if len(pass) >= 1 {
		cmdPath = "/usr/bin/expect"
		args = append(args, expectShellPath)
		args = append(args, user)
		args = append(args, host)
		args = append(args, fmt.Sprintf("%d", port))
		args = append(args, pass)
		args = append(args, fmt.Sprintf("%d", timeout))
		args = append(args, fmt.Sprintf("%d", interval))
	} else {
		cmdPath = "/usr/bin/ssh"
		args = append(args, "-o", fmt.Sprintf("ConnectTimeout=%d", timeout))
		args = append(args, "-o", fmt.Sprintf("ServerAliveInterval=%d", interval))
		args = append(args, "-p", fmt.Sprintf("%d", port))
		args = append(args, fmt.Sprintf("%s@%s", user, host))
		if len(key) >= 1 {
			args = append(args, "-i", key)
		}
	}
	cmd := exec.Command(cmdPath, args...)
	return exe(cmd)
}

func exe(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
