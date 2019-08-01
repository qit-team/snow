## Version 1.1.3 (2019-07-26)

### Changes
- 将任务调度jobs子包涉及消息入队的公共调用抽离成独立包，避免环路调用

## Version 1.1.2 (2019-07-26)

### Changes
- 升级qit-team/snow-core v0.1.7->v0.1.8

## Version 1.1.1 (2019-07-25)

### Changes
- 升级qit-team/work v0.3.3->v0.3.4
- 升级qit-team/snow-core v0.1.5->v0.1.7

## Version 1.1.0 (2019-07-25)

### New Features
- 脚手架：new project、new model
- 支持脚本任务执行模式

### Changes
- 核心组件独立成包
- Queue实现驱动插件式导入机制
- Cache实现驱动插件式导入机制


## Version 1.0.0 (2019-07-08)

### New Features
- 支持多模式：HTTP(平滑重启)、队列调度(平滑结束)、任务调度
- 常用组件支持:
   - Database：MySQL、Postgres、Sqlite3、SQL Server、TiDB...
   - Redis
   - Cache: redis
   - Queue: alimns redis
   - Logger
   - Conifg: toml
   - Reuqest and Response
   - Curl
- 包管理：go module
- 目录结构：大致参照laravel
- 单元测试：部分单测
- 调试：delve
