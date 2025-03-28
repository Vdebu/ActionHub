package main

import (
	"fmt"
	"net/http"
	"time"
)

// 初始化并启动服务器
func (app *application) InitServer() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.cfg.APP.Port), // 根据配置文件中的信息初始化服务器端口
		Handler:      routers(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	err := srv.ListenAndServe()
	return err
}
