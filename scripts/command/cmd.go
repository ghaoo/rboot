package command

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const defaultCmdDir = "command"

var command = make(map[string]Cmd)

type Cmd struct {
	Name    string              `yaml:"name"`
	Rule    string              `yaml:"rule"`
	Usage   string              `yaml:"usage"`
	Version string              `yaml:"version"`
	Cmd     map[string][]string `yaml:"cmd"`
}

func setup(bot *rboot.Robot, in *rboot.Message) []*rboot.Message {
	rule := in.Header.Get("rule")

	cmd := command[rule]

	for name, args := range cmd.Cmd {
		out, err := runCommand(args[0], args[1:]...)
		if err != nil {
			bot.Outgoing(rboot.NewMessage(err.Error()))
		}

		bot.Outgoing(rboot.NewMessage(out, in.From))
		bot.Outgoing(rboot.NewMessage(name+"执行完毕...", in.From))
	}

	return nil
}

func registerCommand() error {
	cmdDir := os.Getenv("COMMAND_DIR")

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
	var usage = ""
	var desc = "命令执行脚本"
	for _, cmd := range cmds {
		command[cmd.Name] = cmd

		ruleset[cmd.Name] = cmd.Rule
		usage += "\n> " + cmd.Usage + "\n\n"
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
			return nil, err
		}

		var cmd = Cmd{}
		err = yaml.Unmarshal(data, &cmd)
		if err != nil {
			return nil, err
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

func runCommand(command string, args ...string) (string, error) {

	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command: %v: %q", err, string(output))
	}

	return string(output), nil
}

func init() {
	if err := registerCommand(); err != nil {
		log.Println(err)
	}
}
