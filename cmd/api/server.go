package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 初始化并启动服务器
func (app *application) InitServer() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.cfg.APP.Port), // 根据配置文件中的信息初始化服务器端口
		Handler:      app.routers(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	// 创建通道接受优雅退出的错误
	shutdownError := make(chan error)
	// 启动goroutine监听信号实现优雅退出
	go func() {
		// 创建带缓存的通道监听信号
		quit := make(chan os.Signal, 1)
		// 使用指定的Signal.Notify监听服务端发出的信息(interrupt/terminate)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// 开始监听默认阻塞
		// 将监听到的信号存放进变量用于后续日志输出
		s := <-quit
		app.logger.Printf("received signal %v, starting graceful shutdown", s)
		// 创建五秒超时的deadline监听后续服务器资源是否完全关闭
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// 回收资源
		defer cancel()
		err := srv.Shutdown(ctx)
		// 判断优雅退出过程中是否发生错误
		if err != nil {
			shutdownError <- err
			return
		}
		// 完成所有请求的事务(waitGroup)
		app.logger.Println("finishing background tasks...")
		// 模拟完成后台任务
		time.Sleep(3 * time.Second)
		// 没错误将nil存入
		shutdownError <- nil
	}()
	err := srv.ListenAndServe()
	// 在终止服务器后会默认生成http.ErrServerClosed错误
	if !errors.Is(err, http.ErrServerClosed) {
		// 若当前的错误并非服务器关闭
		return err
	}
	// 判断优雅退出过程中是否发生了错误
	err = <-shutdownError
	if err != nil {
		return err
	}
	// 服务器优雅退出成功
	app.logger.Println("stopped server...")
	return nil
}
