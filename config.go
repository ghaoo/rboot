package rboot

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	DefaultRbootName      = `RBOOT`
	DefaultRbootAlias     = `rboot`
	DefaultRbootAdapter   = `cli`
	DefaultRbootMemorizer = `bolt`
)

type Config struct {
	Name      string
	Alias     string
	Adapter   string
	Memorizer string
	Debug     bool
}

func newConfig() Config {
	conf := Config{}

	if os.Getenv(`RBOOT_NAME`) != `` {
		conf.Name = os.Getenv(`RBOOT_NAME`)
	} else {
		logrus.Warningf(`RBOOT_NAME not set, using default %s`, DefaultRbootName)
		conf.Name = DefaultRbootName
	}

	if os.Getenv(`RBOOT_ALIAS`) != `` {
		conf.Alias = os.Getenv(`RBOOT_ALIAS`)
	} else {
		logrus.Warningf(`RBOOT_ALIAS not set, using default %s`, DefaultRbootAlias)
		conf.Alias = DefaultRbootAlias
	}

	if os.Getenv(`RBOOT_ADAPTER`) != `` {
		conf.Adapter = os.Getenv(`RBOOT_ADAPTER`)
	} else {
		logrus.Warningf(`RBOOT_ADAPTER not set, using default %s`, DefaultRbootAdapter)
		conf.Adapter = DefaultRbootAdapter
	}

	if os.Getenv(`RBOOT_MEMORIZER`) != `` {
		conf.Memorizer = os.Getenv(`RBOOT_MEMORIZER`)
	} else {
		logrus.Warningf(`RBOOT_MEMORIZER not set, using default %s`, DefaultRbootMemorizer)
		conf.Memorizer = DefaultRbootMemorizer
	}

	debug, _ := strconv.ParseBool(os.Getenv(`DEBUG`))
	conf.Debug = debug

	return conf
}
