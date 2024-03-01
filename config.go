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
	//createExShell()
	createHostConfigList()
}

func clean() {
	//os.Remove(expectShellPath)
}

func createExShell(cmd string) (string, error) {
	var expectShell = `#!/usr/bin/expect

set pass [lindex $argv 0]
set timeout 20

spawn %v
expect {
"*yes/no" { send "yes\r"; exp_continue }
"*password:" { send "$pass\r";exp_continue }
"*login:*" { interact }
"*0%%*" { interact }
}
`
	temp, err := os.CreateTemp(os.TempDir(), "sh")
	if err != nil {
		return "", err
	}
	_, err = temp.WriteString(fmt.Sprintf(expectShell,
		cmd,
	))
	if err != nil {
		return "", err
	}
	return temp.Name(), nil
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
