package trace

import (
	"sync"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/qit-team/snow/config"
)

var (
	tracer *go2sky.Tracer
	lock   sync.Mutex
)

func Tracer() (*go2sky.Tracer, error) {
	if tracer == nil {
		// 有err, 不适合用sync.Once做单例
		lock.Lock()
		defer lock.Unlock()
		if tracer == nil {
			err := InitTracer(config.GetConf().ServiceName, config.GetConf().SkyWalkingOapServer)
			if err != nil {
				return nil, err
			}
		}
	}
	return tracer, nil
}

func InitTracer(serviceName, skyWalkingOapServer string) error {
	report, err := reporter.NewGRPCReporter(skyWalkingOapServer)
	if err != nil {
		return err
	}
	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(report))
	if err != nil {
		return err
	}
	return nil
}
