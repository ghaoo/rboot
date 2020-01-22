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

	for _, plug := range plugs {
		if len(plug.Command) <= 0 {
			log.Warnf("插件脚本 %s 命令集为空，跳过注册", plug.Name)
			continue
		}

		// 注册插件到脚本
		RegisterScripts(plug.Name, Script{
			Action:      setupPlugin,
			Ruleset:     plug.Ruleset,
			Usage:       plug.Usage,
			Description: plug.Description,
		})

		bot.plugins[plug.Name] = plug
	}

	return nil
}

func allPlug(dir string) ([]plugin, error) {

	plugFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}

	var plugs = make([]plugin, 0)

	for _, file := range plugFiles {
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
			"refresh": `^!refresh plugin`,
		},
		Usage: map[string]string{
			"!refresh plugin": "重新加载插件YML配置文件",
		},
		Description: "当插件配置有变化时可运行命令`!refresh plugin`更新插件YML配置文件",
	})
}
