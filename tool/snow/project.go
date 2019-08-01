package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// project project config
type project struct {
	Name       string
	Path       string
	ModuleName string // 支持项目的自定义module名 （go.mod init）
}

const (
	_tplTypeReadme              = iota
	_tplTypeGitignore
	_tplTypeGoMod
	_tplTypeMain
	_tplTypeEnv
	_tplTypeEnvExample
	_tplTypeLog
	_tplTypeCacheKey
	_tplTypeBannerListCache
	_tplTypeBannerListCacheTest
	_tplTypeConsoleKernel
	_tplTypeConsoleTest
	_tplTypeCommand
	_tplTypeConstantCommon
	_tplTypeConstantErrorCode
	_tplTypeConstantLogType
	_tplTypeControllerBase
	_tplTypeControllerTest
	_tplTypeEntity
	_tplTypeFormatter
	_tplTypeFormatterTest
	_tplTypeMiddleWare
	_tplTypeRoute
	_tplTypeJobBase
	_tplTypeJobKernel
	_tplTypeJobTest
	_tplTypeModel
	_tplTypeModelTest
	_tplTypeService
	_tplTypeUtil
	_tplTypeBootstrap
	_tplTypeConfig
	_tplTypeOption
	_tplTypeBuildBin
	_tplTypeBuildShell
)

var (
	p project
	// files type => path
	files = map[int]string{
		// init project
		_tplTypeReadme:     "/README.md",
		_tplTypeGitignore:  "/.gitignore",
		_tplTypeGoMod:      "/go.mod",
		_tplTypeMain:       "/main.go",
		_tplTypeEnv:        "/.env",
		_tplTypeEnvExample: "/.env.example",
		_tplTypeLog:        "/logs/.gitignore",
		//init caches
		_tplTypeCacheKey:            "/app/caches/cache_key.go",
		_tplTypeBannerListCache:     "/app/caches/bannerlistcache/banner_list.go",
		_tplTypeBannerListCacheTest: "/app/caches/bannerlistcache/banner_list_test.go",
		//init console
		_tplTypeConsoleKernel: "/app/console/kernel.go",
		_tplTypeConsoleTest:   "/app/console/test.go",
		_tplTypeCommand:       "/app/console/command.go",
		//init constant
		_tplTypeConstantCommon:    "/app/constants/common/common.go",
		_tplTypeConstantErrorCode: "/app/constants/errorcode/error_code.go",
		_tplTypeConstantLogType:   "/app/constants/logtype/log_type.go",
		//init http
		_tplTypeControllerBase: "/app/http/controllers/base.go",
		_tplTypeControllerTest: "/app/http/controllers/test.go",
		_tplTypeEntity:         "/app/http/entities/test.go",
		_tplTypeFormatter:      "/app/http/formatters/bannerformatter/banner.go",
		_tplTypeFormatterTest:  "/app/http/formatters/bannerformatter/banner_test.go",
		_tplTypeMiddleWare:     "/app/http/middlewares/server_recovery.go",
		_tplTypeRoute:          "/app/http/routes/route.go",
		//init job
		_tplTypeJobBase:   "/app/jobs/basejob/base_job.go",
		_tplTypeJobKernel: "/app/jobs/kernel.go",
		_tplTypeJobTest:   "/app/jobs/test.go",
		//init model
		_tplTypeModel:     "/app/models/bannermodel/banner.go",
		_tplTypeModelTest: "/app/models/bannermodel/banner_test.go",
		//init service
		_tplTypeService: "/app/services/bannerservice/banner.go",
		//init util
		_tplTypeUtil: "/app/utils/.gitkeep",
		//init bootstrap
		_tplTypeBootstrap: "/bootstrap/bootstrap.go",
		//init config
		_tplTypeConfig: "/config/config.go",
		_tplTypeOption: "/config/option.go",
		//init build
		_tplTypeBuildBin:   "/build/bin/.gitignore",
		_tplTypeBuildShell: "/build/shell/build.sh",
	}
	// tpls type => content
	tpls = map[int]string{
		_tplTypeReadme:              _tplReadme,
		_tplTypeGitignore:           _tplGitignore,
		_tplTypeGoMod:               _tplGoMod,
		_tplTypeMain:                _tplMain,
		_tplTypeEnv:                 _tplEnv,
		_tplTypeEnvExample:          _tplEnv,
		_tplTypeLog:                 _tplLog,
		_tplTypeCacheKey:            _tplCacheKey,
		_tplTypeBannerListCache:     _tplBannerListCache,
		_tplTypeBannerListCacheTest: _tplBannerListCacheTest,
		_tplTypeConsoleKernel:       _tplConsoleKernel,
		_tplTypeConsoleTest:         _tplConsoleTest,
		_tplTypeCommand:             _tplCommand,
		_tplTypeConstantCommon:      _tplConstantCommon,
		_tplTypeConstantErrorCode:   _tplConstantErrorCode,
		_tplTypeConstantLogType:     _tplConstantLogType,
		_tplTypeControllerBase:      _tplControllerBase,
		_tplTypeControllerTest:      _tplControllerTest,
		_tplTypeEntity:              _tplEntity,
		_tplTypeFormatter:           _tplFormatter,
		_tplTypeFormatterTest:       _tplFormatterTest,
		_tplTypeMiddleWare:          _tplMiddleWare,
		_tplTypeRoute:               _tplRoute,
		_tplTypeJobBase:             _tplJobBase,
		_tplTypeJobKernel:           _tplJobKernel,
		_tplTypeJobTest:             _tplJobTest,
		_tplTypeModel:               _tplModel,
		_tplTypeModelTest:           _tplModelTest,
		_tplTypeService:             _tplService,
		_tplTypeUtil:                _tplUtil,
		_tplTypeBootstrap:           _tplBootstrap,
		_tplTypeConfig:              _tplConfig,
		_tplTypeOption:              _tplOption,
		_tplTypeBuildBin:            _tplBuildBin,
		_tplTypeBuildShell:          _tplBuildShell,
	}
)

func create() (err error) {
	if err = os.MkdirAll(p.Path, 0755); err != nil {
		return
	}
	for t, v := range files {
		i := strings.LastIndex(v, "/")
		if i > 0 {
			dir := v[:i]
			if err = os.MkdirAll(p.Path+dir, 0755); err != nil {
				return
			}
		}
		if err = write(p.Path+v, tpls[t], p); err != nil {
			return
		}
	}
	return
}

func write(name, tpl string, data interface{}) (err error) {
	body, err := parse(tpl, data)
	if err != nil {
		return
	}
	return ioutil.WriteFile(name, body, 0644)
}

func parse(s string, data interface{}) ([]byte, error) {
	t, err := template.New("").Parse(s)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
