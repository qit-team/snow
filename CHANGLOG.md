## Version 1.1.0 (2019-07-25)

### New Features
- 脚手架：new project

### Change
- 核心组件独立成包
- Queue实现驱动插件式支持
- Cache实现驱动插件式支持
- DB的Mysql驱动默认不导入，由业务自行导入


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
