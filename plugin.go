// +build darwin dragonfly freebsd linux netbsd openbsd

package rboot

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultPluginDir = "plugins"

// plugin 支持脚本语言执行的脚本插件
type plugin struct {
	Name string
	Script
}

func (plug plugin) runScript(bot *Robot, in *Message) []*Message {
	rule := in.Header.Get("rule")
	cmd := exec.Command(rule, plug.Ruleset[rule])

	out, err := cmd.Output()
	if err != nil {
		return NewMessages(err.Error(), in.From)
	}

	return NewMessages(string(out), in.From)
}

// 注册插件
func registerPlugin() {
	dir := os.Getenv("PLUGIN_DIR")

	if dir == "" {
		dir = defaultPluginDir
	}

	plugs, err := loadPlugins(dir)
	if err != nil {
		log.Error(err)
		return
	}

	if len(plugs) <= 0 {
		log.Warn("no plug-in found")
		return
	}

	for _, plug := range plugs {
		if len(plug.Ruleset) > 0 {
			RegisterScripts(plug.Name, Script{
				Action:      plug.runScript,
				Ruleset:     plug.Ruleset,
				Usage:       plug.Usage,
				Description: plug.Description,
			})
		}
	}
}

// loadPlugins 加载所有Plugin
func loadPlugins(dir string) ([]plugin, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	var plugins = make([]plugin, 0)

	for _, file := range files {
		// 获取配置内容
		conf, err := readConf(file)
		if err != nil {
			log.Errorf("plugin %s exec failure，err: %v", file, err)
			continue
		}

		var p plugin
		if err = json.Unmarshal(conf, &p); err != nil {
			log.Errorf("plugin %s register failed, err: %v", file, err)
			continue
		}

		plugins = append(plugins, p)
	}

	return plugins, nil
}

// readConf 读取配置文件内容
func readConf(file string) ([]byte, error) {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		return nil, err
	}

	// 每个脚本文件需要有一个 init 命令解析函数，用来将脚本信息注册到rboot
	cmd := exec.Command(file, "init")

	return cmd.Output()
}

func init() {
	registerPlugin()
}
