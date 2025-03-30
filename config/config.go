package config

import (
	"github.com/spf13/viper"
	"log"
)

// 存储yml配置文件中的相关性信息
type Config struct {
	// 根据yml文件中的相应顺序进行读入
	APP struct {
		Name string
		Port string
	}
	Database struct {
		MySqlDSN     string
		MaxIdleConns int
		MaxOpenConns int
	}
	Redis struct {
		Host     string
		Port     string
		DB       int
		Password string
	}
	DockerDeploy bool
}

// 读取外部yml配置文件
func InitConfig() (*Config, error) {
	// 设置外部配置文件的相关信息
	// 设置文件名称(无需包含拓展名)
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yml")
	// 设置配置文件路径
	viper.AddConfigPath("./config")
	// 读取文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file %v", err)
		return nil, err
	}
	// 初始化配置文件指!针!
	AppConfig := &Config{}
	// 将读取到的信息存储进结构体
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("unable to decode configs %v", err)
		return nil, err
	}
	return AppConfig, nil
}
