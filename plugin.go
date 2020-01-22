package rboot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
)

const defaultPlugDir = "plugins"

type plugin struct {
	Name    string            `yaml:"name"`
	Rule    string            `yaml:"rule"`
	Usage   map[string]string `yaml:"usage"`
	Version string            `yaml:"version"`
	Command []command         `yaml:"command"`
}

type command struct {
	Dir string   `yaml:"dir"`
	Cmd []string `yaml:"cmd"`
}

func (bot *Robot) registerPlugin() error {
	plugDir := os.Getenv("PLUGIN_DIR")

	if plugDir == "" {
		plugDir = defaultPlugDir
	}

	plugs, err := allPlug(plugDir)
	if err != nil {
		return err
	}

	if len(plugs) <= 0 {
		return fmt.Errorf("no plug-in found")
	}

	var _ruleset = make(map[string]string)
	var _usage = make(map[string]string)
	var _desc = "脚本插件"
	for _, plug := range plugs {
		if len(plug.Command) <= 0 {
			log.Warnf("插件脚本 %s 命令集为空，跳过注册", plug.Name)
			continue
		}
		bot.plugins[plug.Name] = plug
		_ruleset[plug.Name] = plug.Rule

		for _rule, _explain := range plug.Usage {
			_usage[_rule] = _explain
		}
	}

	if len(_ruleset) > 0 {
		RegisterScripts("plug", Script{
			Action:      setupPlugin,
			Ruleset:     _ruleset,
			Usage:       _usage,
			Description: _desc,
		})
	}

	return nil
}

func allPlug(dir string) ([]plugin, error) {

	cmdFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}

	var plugs = make([]plugin, 0)

	for _, file := range cmdFiles {
		data, err := loadPlugin(file)
		if err != nil {
			log.Errorln(err)
			continue
		}

		var plug = plugin{}
		err = yaml.Unmarshal(data, &plug)
		if err != nil {
			log.Println(err)
			continue
		}

		plugs = append(plugs, plug)
	}

	return plugs, nil
}

func loadPlugin(file string) ([]byte, error) {
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
		logrus.Error(err, "\n", string(output))
		return "", fmt.Errorf("run plugin error: %v: %q", err, string(output))
	}

	return string(output), nil
}

func init() {

	RegisterScripts("refresh_plugin", Script{
		Action: func(bot *Robot, incoming *Message) []*Message {
			err := bot.registerPlugin()
			if err != nil {
				log.Println(err)
				return NewMessages(err.Error(), incoming.From)
			}

			return NewMessages("更新成功！", incoming.From)
		},
		Ruleset: map[string]string{
			"refresh": `^!refresh yml`,
		},
		Usage: map[string]string{
			"!refresh plugin": "重新加载YML配置文件",
		},
		Description: "当插件有变化时可运行此命令更新",
	})
}
