package initialization

import (
	"github.com/BurntSushi/toml"
	"log"
)

type configuration struct {
	Setting setting
	RedisSetting redis
}

type setting struct {
	Ip string
	Port string
	LogDir string
}

type redis struct {
	Exists bool
	Ip string
	Port string
	Password string
}

var Configuration configuration

func init() {
	confFile := "conf.toml"
	if _, err := toml.DecodeFile(confFile, &Configuration); err != nil {
		log.Fatal(err)
	}
}
