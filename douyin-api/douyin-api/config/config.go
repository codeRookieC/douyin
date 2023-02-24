package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var C = InitConfig()

type Config struct {
	viper *viper.Viper //读取配置文件用
	SC    *ServerConfig
}

type ServerConfig struct {
	Name string
	Addr string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	//获取当前目录 添加config路径
	workDir, _ := os.Getwd()
	conf.viper.AddConfigPath(workDir + "/config")
	//读取配置文件 这里conf.viper表示的是读取该路径下的config.yaml文件
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	conf.ReadServerConfig()
	return conf
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}
