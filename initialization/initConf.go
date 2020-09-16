package initialization

import (
	"github.com/BurntSushi/toml"
	"log"
)

type configuration struct {
	Setting setting
	RedisSetting redisSetting
}

type setting struct {
	Ip string
	Port string
	LogDir string
}

type redisSetting struct {
	Exists bool
	Ip string
	Port string
	Password string
	Db int
}

var Configuration configuration

func init() {
	confFile := "conf.toml"
	if _, err := toml.DecodeFile(confFile, &Configuration); err != nil {
		log.Fatal(err)
	}
}
