package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var DefaultCiphers = []string{
	"aes128-ctr",
	"aes192-ctr",
	"aes256-ctr",
	"aes128-gcm@openssh.com",
	"chacha20-poly1305@openssh.com",
	"arcfour256",
	"arcfour128",
	"arcfour",
	"aes128-cbc",
	"3des-cbc",
	"blowfish-cbc",
	"cast128-cbc",
	"aes192-cbc",
	"aes256-cbc",
}

func main() {
	for {
		printServers()
		command := getCommand()

		if command == "q" {
			break
		}

		atoi, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println(fmt.Errorf("输入错误:%w", err))
			continue
		}

		hostConfig := HostConfigList[atoi]
		err = hostConfig.Dail()
		if err != nil {
			fmt.Println(fmt.Errorf("连接失败:%w", err))
		}
	}
	clean()
}

func printServers() {
	fmt.Println("----------服务器----------")
	for index, hostConfig := range HostConfigList {
		fmt.Println(index, hostConfig.Name)
	}
	fmt.Println("-------------------------")
	fmt.Println("退出输入'q'\n")
}
func getCommand() string {
	var command string
	for {
		fmt.Print("请输入：")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command = scanner.Text()
		command = strings.TrimSpace(command)
		if len(command) >= 1 {
			break
		}
	}
	return command
}

func dialSSHUseCommand(host string, port int, user string, key string) error {

	var args []string
	args = append(args, "-o", "ServerAliveInterval=60")
	args = append(args, "-o", "ConnectTimeout=10")
	args = append(args, "-p", fmt.Sprintf("%d", port))
	args = append(args, fmt.Sprintf("%s@%s", user, host))
	if len(key) >= 1 {
		args = append(args, "-i", key)
	}

	cmd := exec.Command("/usr/bin/ssh", args...)
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

func dialSShWithPassword(host string, port int, user, pass string) error {

	cmd := exec.Command("/usr/bin/expect", expectShellPath, user, host, fmt.Sprintf("%d", port), pass)
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

func dialSShWithPubKey_bak(host string, port int, user, key string) error {
	file, err := os.ReadFile(key)
	if err != nil {
		return fmt.Errorf("读文件失败:%w", err)
	}
	privateKey, err := ssh.ParsePrivateKey(file)

	return dialSSH(host, port, user, ssh.PublicKeys(privateKey))
}
func dialSShWithSSHPubKey_bak(host string, port int, user string) error {

	pk := "/Users/lizhikui/.ssh/id_ed25519"

	file, err := os.ReadFile(pk)
	if err != nil {
		return fmt.Errorf("读文件失败:%w", err)
	}
	key, err := ssh.ParsePrivateKey(file)

	return dialSSH(host, port, user, ssh.PublicKeys(key))
}

func dialSSH(host string, port int, user string, authMethod ssh.AuthMethod) error {

	cientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		Timeout:         time.Second * 30,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cientConfig.SetDefaults()
	cientConfig.Ciphers = append(cientConfig.Ciphers, DefaultCiphers...)

	dial, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), cientConfig)

	if err != nil {
		return fmt.Errorf("连接ssh失败:%w", err)
	}
	defer dial.Close()

	session, err := dial.NewSession()
	if err != nil {
		return fmt.Errorf("开启会话失败:%w", err)
	}

	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("获取当前终端的大小失败:%w", err)
	}

	err = session.RequestPty("xterm-256color", height, width, modes)
	if err != nil {
		return fmt.Errorf("启动终端失败:%w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	//stdinPipe, err := session.StdinPipe()
	//if err != nil {
	//	return fmt.Errorf("启动终端失败:%w", err)
	//}

	err = session.Shell()
	if err != nil {
		return fmt.Errorf("开始shell失败:%w", err)
	}

	//go func() {
	//	_, err = io.Copy(stdinPipe, os.Stdin)
	//	session.Close()
	//}()

	err = session.Wait()
	if err != nil {
		return fmt.Errorf("等待session结束失败:%w", err)
	}

	return nil
}
