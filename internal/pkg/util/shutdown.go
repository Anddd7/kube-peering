package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitElegantExit() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 启动一个 goroutine 监听退出信号
	go func() {
		sig := <-done
		fmt.Printf("接收到退出信号 %v，开始清理...\n", sig)

		// 执行清理工作
		// ...

		fmt.Println("清理工作完成，程序退出")
		os.Exit(0)
	}()

	// 等待退出信号
	fmt.Println("程序已启动，等待退出信号...")
	<-done
}
