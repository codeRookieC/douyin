package douyin_common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 优雅地启停web服务
func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	//gin 官方文档中的优雅启停方式
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		//保证服务端的运行
		log.Printf("%s running in %s \n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	//使用信号量来优雅停止web服务的运行
	quit := make(chan os.Signal)
	//SIGINT 用户发送INTR字符(Ctrl+C)触发
	//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	//Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。 使用单一信号进行通知的信道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s...", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //计时器
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : ", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout!")
	}
	log.Printf("%s stop success...", srvName)
}
