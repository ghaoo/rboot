package rboot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
)

const defaultScriptDir = "scripts"

var scripts = make(map[string]script)

type script struct {
	Name        string            `yaml:"name"`
	Ruleset     map[string]string `yaml:"ruleset"`
	Usage       map[string]string `yaml:"usage"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	Command     []command         `yaml:"command"`
}

type command struct {
	Dir string   `yaml:"dir"`
	Cmd []string `yaml:"cmd"`
}

// 设置脚本
func setupScript(bot *Robot, in *Message) (msg []*Message) {
	rule := in.Header.Get("rule")

	scp := scripts[rule]

	for _, sc := range scp.Command {
		for _, c := range sc.Cmd {
			args := strings.Split(c, " ")

			out, err := runCommand(sc.Dir, args[0], args[1:]...)
			if err != nil {
				return NewMessages(err.Error())
			}

			msg = append(msg, NewMessage(out, in.From))
		}
	}

	return msg
}

// 注册脚本
func registerScript() error {
	scpDir := os.Getenv("SCRIPT_DIR")

	if scpDir == "" {
		scpDir = defaultScriptDir
	}

	scps, err := allScripts(scpDir)
	if err != nil {
		return err
	}

	if len(scps) <= 0 {
		return fmt.Errorf("no script found")
	}

	for _, scp := range scps {
		if len(scp.Command) <= 0 {
			continue
		}

		// 注册插件到脚本
		RegisterPlugin(scp.Name, Plugin{
			Action:      setupScript,
			Ruleset:     scp.Ruleset,
			Usage:       scp.Usage,
			Description: scp.Description,
		})

		scripts[scp.Name] = scp
	}

	return nil
}

func allScripts(dir string) ([]script, error) {

	scpFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}

	var scps = make([]script, 0)

	for _, file := range scpFiles {
		data, err := loadScript(file)
		if err != nil {
			continue
		}

		var scp = script{}
		err = yaml.Unmarshal(data, &scp)
		if err != nil {
			continue
		}

		scps = append(scps, scp)
	}

	return scps, nil
}

func loadScript(file string) ([]byte, error) {
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
		return "", fmt.Errorf("run script error: %v - %q", err, string(output))
	}

	return string(output), nil
}

func init() {

	err := registerScript()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"mod": `rboot`,
		}).Errorf("register script err: %v", err)
	}

	RegisterPlugin("refresh_script", Plugin{
		Action: func(bot *Robot, incoming *Message) []*Message {
			err := registerScript()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"mod": `rboot`,
				}).Error(err)
				return NewMessages(err.Error(), incoming.From)
			}

			return NewMessages("更新成功！", incoming.From)
		},
		Ruleset: map[string]string{
			"refresh": `^!refresh scripts`,
		},
		Usage: map[string]string{
			"!refresh scripts": "重新加载插件YML配置文件",
		},
		Description: "当插件配置有变化时更新插件YML配置文件",
	})
}
