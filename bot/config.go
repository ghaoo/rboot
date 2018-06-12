package bot

import (
	"os"
	"io/ioutil"
	"github.com/go-yaml/yaml"
)

var confile = "config.yml"

var yaml_setting = `
# rboot setting

# robot name
name: Rboot

# robot connecter name
connecter: cli

# robot enable plugins
plugins:
 - testing

`

type Config struct {
	Name string `yaml:"name"`
	Connecter string `yaml:"connecter"`
	Plugins []string `yaml:"plugins"`
}


func load() ([]byte, error) {
	_, err := os.Stat(confile)

	if os.IsNotExist(err) {
		createConf()
	}

	file, err := os.Open(confile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func createConf() {
	_, err := os.Stat(confile)
	if os.IsNotExist(err) {
		_, err := os.OpenFile(confile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		var confwrite = []byte(yaml_setting)
		err = ioutil.WriteFile(confile, confwrite, 0666) //写入文件(字节数组)
		if err != nil {
			panic(err)
		}
	}
}

func NewConf() Config {
	data, err := load()
	if err != nil {
		panic("加载配置文件失败"+ err.Error())
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic("解析配置文件失败"+err.Error())
	}

	return c
}