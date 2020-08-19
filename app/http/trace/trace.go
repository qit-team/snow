package trace

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/qit-team/snow/config"
)

var tracer *go2sky.Tracer

func Tracer() (*go2sky.Tracer, error) {
	if tracer == nil {
		err := InitTracer(config.GetConf().ServiceName, config.GetConf().SkyWalkingOapServer)
		if err != nil {
			return nil, err
		}
	}

	return tracer, nil
}

func InitTracer(serviceName, skyWalkingOapServer string) error {
	var (
		report go2sky.Reporter
		err    error
	)
	report, err = reporter.NewGRPCReporter(skyWalkingOapServer)
	if err != nil {
		return err
	}

	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(report))
	if err != nil {
		return err
	}
	return nil
}
