package ymlcmd

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultCmdDir = "command"

var command = make(map[string]Cmd)

type Cmd struct {
	Name    string            `yaml:"name"`
	Rule    string            `yaml:"rule"`
	Usage   map[string]string `yaml:"usage"`
	Version string            `yaml:"version"`
	Command []Command         `yaml:"command"`
}

type Command struct {
	Dir string
	Cmd []string
}

func registerCommand() error {
	cmdDir := os.Getenv("YML_COMMAND_DIR")

	if cmdDir == "" {
		cmdDir = defaultCmdDir
	}

	cmds, err := allCmd(cmdDir)
	if err != nil {
		return err
	}

	if len(cmds) <= 0 {
		return fmt.Errorf("no command found")
	}

	var ruleset = make(map[string]string)
	var usage = make(map[string]string)
	var desc = "YML命令执行脚本"
	for _, cmd := range cmds {
		command[cmd.Name] = cmd
		ruleset[cmd.Name] = cmd.Rule
		for _rule, _explain := range cmd.Usage {
			usage[_rule] = _explain
		}
	}

	if len(ruleset) > 0 {
		rboot.RegisterScripts("cmd", rboot.Script{
			Action:      setup,
			Ruleset:     ruleset,
			Usage:       usage,
			Description: desc,
		})
	}

	return nil
}

func allCmd(dir string) ([]Cmd, error) {

	cmdFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}

	var cmds = make([]Cmd, 0)

	for _, file := range cmdFiles {
		data, err := load(file)
		if err != nil {
			log.Println(err)
			continue
		}

		var cmd = Cmd{}
		err = yaml.Unmarshal(data, &cmd)
		if err != nil {
			log.Println(err)
			continue
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func load(file string) ([]byte, error) {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		return nil, err
	}

	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	return ioutil.ReadAll(fi)
}

func runCommand(dir, command string, args ...string) (string, error) {

	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command: %v: %q", err, string(output))
	}

	return string(output), nil
}

func init() {
	err := registerCommand()
	if err != nil {
		log.Println("register yml cmd err: ", err)
	}
	rboot.RegisterScripts("refresh_yml", rboot.Script{
		Action: func(bot *rboot.Robot, incoming *rboot.Message) []*rboot.Message {
			err := registerCommand()
			if err != nil {
				log.Println(err)
				return rboot.NewMessages(err.Error(), incoming.From)
			}

			return rboot.NewMessages("更新成功！", incoming.From)
		},
		Ruleset: map[string]string{
			"refresh": `^!refresh yml`,
		},
		Usage: map[string]string{
			"!refresh yml": "重新加载YML文件",
		},
		Description: "当command有变化时可运行次命令更新",
	})
}
