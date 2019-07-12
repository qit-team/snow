package server

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/qit-team/snow/config"
	"strconv"
	"github.com/fvbock/endless"
	"syscall"
)

/**
 * 启动gin引擎
 * @wiki https://github.com/fvbock/endless#signals
 */
func runEngine(engine *gin.Engine, addr string, pidPath string) error {
	//设置gin模式
	if config.IsEnvEqual(config.ProdEnv) {
		gin.SetMode(gin.ReleaseMode)
	}

	server := endless.NewServer(addr, engine)
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		if gin.Mode() != gin.ReleaseMode {
			fmt.Printf("Actual pid is %d \n\r", pid)
		}
		writePidFile(pidPath, pid)
	}
	err := server.ListenAndServe()
	return err
}

// Start proxy with config file
func StartHttp(confFile, pidFile string, boot func(config *config.Config) error, registerRoute func(*gin.Engine)) error {
	//加载配置文件
	conf, err := config.Load(confFile)
	if err != nil {
		return err
	}

	//初始化服务信息
	err = initServer()
	if err != nil {
		return fmt.Errorf("init server failed, %s", err.Error())
	}

	//容器初始化
	err = boot(conf)
	if err != nil {
		return fmt.Errorf("container ini failed %s", err.Error())
	}

	//配置路由引擎
	engine := gin.Default()
	registerRoute(engine)
	addr := conf.Api.Host + ":" + strconv.Itoa(conf.Api.Port)
	runEngine(engine, addr, pidFile)

	//因为信号处理由endless接管实现平滑重启和关闭，这里模拟通用的结束信号
	go func() {
		srv.stop <- true
	}()

	//等待停止信号
	waitStop()
	return nil
}
