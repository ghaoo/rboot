package rboot

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

//var confile = "config.yml"

var yaml_setting = `
# Rboot setting

name: Rboot

connecter: cli

# enable plugins
plugins:
 - testing

`

type Config struct {
	Name      string   `yaml:"name"`
	Connecter string   `yaml:"connecter"`
	Plugins   []string `yaml:"plugins"`
}

func load(confpath string) ([]byte, error) {
	_, err := os.Stat(confpath)

	if os.IsNotExist(err) {
		createConf(confpath)
	}

	file, err := os.Open(confpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func createConf(confpath string) {
	_, err := os.Stat(confpath)
	if os.IsNotExist(err) {
		_, err := os.OpenFile(confpath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		var confwrite = []byte(yaml_setting)
		err = ioutil.WriteFile(confpath, confwrite, 0666) //写入文件(字节数组)
		if err != nil {
			panic(err)
		}
	}
}

func NewConf(confpath string) Config {
	data, err := load(confpath)
	if err != nil {
		panic("加载配置文件失败" + err.Error())
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic("解析配置文件失败" + err.Error())
	}

	return c
}
