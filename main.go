package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
