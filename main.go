package main

import (
	"bufio"
	"fmt"
	"github.com/lzk97224/igo/islice"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	defer clean()
	printHelp()
	printServers()
	for {
		command, err := parseCommand(getCommand())
		if err != nil {
			print(err)
			continue
		}
		fun, ok := RunCommandMap[command.Name]
		if ok {
			fun(command.Args...)
		} else {
			switch command.Name {
			case "q":
				return
			default:
				entrySsh(string(command.Name))
			}
		}
	}
}

var RunCommandMap = map[CommandName]func(...string){}

func init() {
	RunCommandMap[Print] = printServers
	RunCommandMap[Help] = printHelp
	RunCommandMap[SCPUP] = scpUp
	RunCommandMap[SCPDOWN] = scpDown
}

type CommandName string

const (
	Print   CommandName = "p"
	Help    CommandName = "h"
	SCPUP   CommandName = "up"
	SCPDOWN CommandName = "down"
)

type Command struct {
	Name CommandName
	Args []string
}

func parseCommand(command string) (*Command, error) {
	command = strings.TrimSpace(command)
	coms := strings.Split(command, " ")
	if len(coms) <= 0 {
		return nil, fmt.Errorf("参数为空")
	}

	coms = islice.Filter(coms, func(item string) bool {
		return strings.TrimSpace(item) == ""
	})

	return &Command{Name: CommandName(coms[0]), Args: coms[1:]}, nil
}

func entrySsh(index string) {
	atoi, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println(fmt.Errorf("输入错误:%w", err))
		return
	}

	hostConfig := HostConfigList[atoi]
	err = hostConfig.Dail()
	if err != nil {
		fmt.Println(fmt.Errorf("连接失败:%w", err))
	}
}

func scpUp(args ...string) {
	if len(args) != 3 {
		fmt.Println("参数错误")
		return
	}
	atoi, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(fmt.Errorf("输入错误:%w", err))
		return
	}

	config := HostConfigList[atoi]
	c := exec.Command(
		"scp",
		"-P", fmt.Sprintf("%v", config.Port),
		fmt.Sprintf("\"%v\"", args[0]),
		fmt.Sprintf("\"%v@%v:%v\"", config.User, config.Host, args[2]),
	)

	exeShell(c, config.Pass)
}

func scpDown(args ...string) {
	//down 1 /tmp/sss.txt /sss
	if len(args) != 3 {
		fmt.Println("参数错误")
		return
	}
	atoi, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(fmt.Errorf("输入错误:%w", err))
		return
	}

	config := HostConfigList[atoi]
	c := exec.Command(
		"scp",
		"-P", fmt.Sprintf("%v", config.Port),
		fmt.Sprintf("\"%v@%v:%v\"", config.User, config.Host, args[1]),
		fmt.Sprintf("\"%v\"", args[2]),
	)

	fmt.Println(c.String())

	exeShell(c, config.Pass)
}

func printServers(...string) {
	fmt.Println("----------服务器----------")
	for index, hostConfig := range HostConfigList {
		fmt.Println(index, hostConfig.Name)
	}
	fmt.Println("-------------------------")
}
func printHelp(...string) {
	fmt.Println()
	fmt.Println("帮助信息----'h'")
	fmt.Println("服务列表----'p'")
	fmt.Println("退出输入----'q'")
	fmt.Println()
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
