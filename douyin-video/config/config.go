package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var C = InitConfig()

type Config struct {
	viper       *viper.Viper
	GC          *GrpcConfig
	EtcdConfig  *EtcdConfig
	MysqlConfig *MysqlConfig
	JwtConfig   *JwtConfig
	Oss         *OssConfig
}
type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}
type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
}
type OssConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	Codedir         string
}
type JwtConfig struct {
	AccessExp     int
	RefreshExp    int
	AccessSecret  string
	RefreshSecret string
}
type EtcdConfig struct {
	Addrs []string //存储多个地址 同一名字的不同地址 表示负载均衡
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	workDir, _ := os.Getwd()
	conf.viper.AddConfigPath(workDir + "/config")
	//读取配置文件 这里conf.viper表示的是读取该目录下的config.yaml文件
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// conf.ReadServerConfig()
	// conf.InitZapLog()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.InitMysqlConfig()
	conf.InitJwtConfig()
	conf.InitOssConfig()
	return conf
}

// func (c *Config) InitZapLog() {
// 	//从配置中读取日志配置，初始化日志
// 	lc := &logs.LogConfig{
// 		DebugFileName: c.viper.GetString("zap.debugFileName"),
// 		InfoFileName:  c.viper.GetString("zap.infoFileName"),
// 		WarnFileName:  c.viper.GetString("zap.warnFileName"),
// 		MaxSize:       c.viper.GetInt("maxSize"),
// 		MaxAge:        c.viper.GetInt("maxAge"),
// 		MaxBackups:    c.viper.GetInt("maxBackups"),
// 	}
// 	err := logs.InitLogger(lc)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

// func (c *Config) ReadServerConfig() {
// 	sc := &ServerConfig{}
// 	sc.Name = c.viper.GetString("server.name")
// 	sc.Addr = c.viper.GetString("server.addr")
// 	c.SC = sc
// }

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = gc
}

// func (c *Config) ReadRedisConfig() *redis.Options {
// 	return &redis.Options{
// 		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
// 		Password: c.viper.GetString("redis.password"),
// 		DB:       c.viper.GetInt("redis.db"),
// 	}
// }

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}

func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}
	c.MysqlConfig = mc
}

func (c *Config) InitOssConfig() {
	oss := &OssConfig{
		Endpoint:        c.viper.GetString("oss.endpoint"),
		AccessKeyId:     c.viper.GetString("oss.accessKeyId"),
		AccessKeySecret: c.viper.GetString("oss.accessKeySecret"),
		BucketName:      c.viper.GetString("oss.bucketName"),
		Codedir:         c.viper.GetString("oss.codedir"),
	}
	c.Oss = oss
}

func (c *Config) InitJwtConfig() {
	jc := &JwtConfig{
		AccessExp:     c.viper.GetInt("jwt.accessExp"),
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		RefreshExp:    c.viper.GetInt("jwt.refreshExp"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
	}
	c.JwtConfig = jc
}