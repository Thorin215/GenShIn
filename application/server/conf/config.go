package conf

import (
	"gopkg.in/ini.v1"
)

var Conf = new(Config)

type Config struct {
	*MysqlConfig  `ini:"mysql"`
	*ServerConfig `ini:"server"`
}

type MysqlConfig struct {
	Conn string `ini:"conn"`
}

type ServerConfig struct {
	Host string `ini:"host"`
	Port string `ini:"port"`
}

func Init() error {
	return ini.MapTo(Conf, "config.ini")
}
