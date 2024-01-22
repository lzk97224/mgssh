package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var HostConfigList []HostConfig
var expectShellPath string

func init() {
	createExShell()
	createHostConfigList()
}

func clean() {
	os.Remove(expectShellPath)
}

func createExShell() {
	var expectShell = `#!/usr/bin/expect

set user [lindex $argv 0]
set ip [lindex $argv 1]
set port [lindex $argv 2]
set password [lindex $argv 3]
set timeout [lindex $argv 4]
set interval [lindex $argv 5]

spawn /usr/bin/ssh -o ConnectTimeout=$timeout -o ServerAliveInterval=$interval -p $port $user@$ip
expect {
"*yes/no" { send "yes\r"; exp_continue }
"*password:" { send "$password\r" }
}
interact
`
	temp, err := os.CreateTemp(os.TempDir(), "sh")
	if err != nil {
		panic(err)
	}
	_, err = temp.WriteString(expectShell)
	if err != nil {
		panic(err)
	}
	expectShellPath = temp.Name()
}

func createHostConfigList() {
	readFile, err := os.ReadFile(getConfigFile())
	if err != nil {
		panic(err)
	}
	var config []HostConfig
	err = json.Unmarshal(readFile, &config)
	if err != nil {
		panic(err)
	}
	HostConfigList = config
}

func getConfigFile() string {
	runFile, _ := exec.LookPath(os.Args[0])
	runFilePath, err := filepath.Abs(runFile)

	dir := filepath.Dir(runFilePath)
	filePath, err := fileExists(dir)
	if err == nil {
		return filePath
	}
	log.Printf("%v", err)

	dir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filePath, err = fileExists(dir)
	if err == nil {
		return filePath
	}
	log.Printf("%v", err)

	dir, err = os.Getwd()
	filePath, err = fileExists(dir)
	if err == nil {
		return filePath
	}

	log.Printf("%v", err)
	panic(err)
}
func fileExists(dir string) (string, error) {
	configName := "mgssh_config.json"
	filePath := filepath.Join(dir, configName)
	_, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("配置文件，%v，不存在,%w", filePath, err)
	}
	return filePath, nil
}
