package server

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"github.com/qit-team/snow/pkg/kernel/close"
)

const (
	Version     = "1.0"
	BuildCommit = ""
	BuildDate   = ""
)

type serverInfo struct {
	stop chan bool
}

var srv *serverInfo

//将进程号写入文件
func writePidFile(path string, pidArgs ...int) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	var pid int
	if len(pidArgs) > 0 {
		pid = pidArgs[0]
	} else {
		pid = os.Getpid()
	}
	_, err = fd.WriteString(fmt.Sprintf("%d\n", pid))
	return err
}

//读取文件的进程号
func readPidFile(path string) (int, error) {
	fd, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer fd.Close()

	buf := bufio.NewReader(fd)
	line, err := buf.ReadString('\n')
	if err != nil {
		return -1, err
	}
	line = strings.TrimSpace(line)
	return strconv.Atoi(line)
}

//阻塞等待程序内部的Stop通道信号
func waitStop() {
	<-srv.stop
	fmt.Println("close service")
	close.Free()
}

//处理进程的信号量
func handleSignal(sig os.Signal) {
	switch sig {
	case syscall.SIGINT:
		fallthrough
	case syscall.SIGTERM:
		Stop()
	default:
	}
}

//监听信号量
func registerSignal() {
	go func() {
		var sigs = []os.Signal{
			syscall.SIGHUP,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
			syscall.SIGINT,
			syscall.SIGTERM,
		}
		c := make(chan os.Signal)
		signal.Notify(c, sigs...)
		for {
			sig := <-c //blocked
			handleSignal(sig)
		}
	}()
}

// HandleUserCmd use to stop/reload the proxy service
func HandleUserCmd(cmd string, pidFile string) error {
	var sig os.Signal

	switch cmd {
	case "stop":
		sig = syscall.SIGTERM
	case "restart":
		//目前api使用endless平滑重启，需要传递此信号，其他只需要平滑关闭就可以了
		sig = syscall.SIGHUP
	default:
		return fmt.Errorf("unknown user command %s", cmd)
	}

	pid, err := readPidFile(pidFile)
	if err != nil {
		return err
	}

	fmt.Printf("send %v to pid %d \n", sig, pid)

	proc := new(os.Process)
	proc.Pid = pid
	return proc.Signal(sig)
}

// Stop proxy
func Stop() {
	srv.stop <- true
}

func initServer() error {
	srv = new(serverInfo)
	srv.stop = make(chan bool, 0)

	return nil
}
