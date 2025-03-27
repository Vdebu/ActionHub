package main

import (
	"ActionHub/config"
	"fmt"
	"log"
)

func main() {
	// 初始化配置文件
	config.InitConfig()
	// 查看配置文件是否解码成功
	fmt.Println(config.AppConfig)
	err := InitServer()
	if err != nil {
		log.Fatal("failed to starting server %v", err)
	}
}
