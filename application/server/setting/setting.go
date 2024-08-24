package setting

import (
	"gopkg.in/ini.v1"
)

var Conf =new(Config)


type Config struct{
	*MysqlConfig `ini:"mysql"`
}

type MysqlConfig struct{
	User string `ini:"user"`
	Password string `ini:"password"`
	Host string `ini:"host"`
	Port string `ini:"port"`
	Database string `ini:"db"`
}

func Init(file string )error{
	return ini.MapTo(Conf,file)
}