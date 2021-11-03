# restfula_api_gin

a restfula_api template by gin 



需要的依赖包都在go.mod里，直接下载即可

原理为类似nginx的master主进程和work进程



目前启动项目为命令行:

go run main.go server

go run main.go worker



main.bak为之前的整合main文件

# 项目解析

解耦后的项目，主体为gallery



logging是用的zap包 



1、async为异步的包，使用的machinery包

​	(1)异步需要使用redis作为缓存，所以config.yml为连接redis的配置文件



2、api为views接口处理功能和pagination分页



3、models为数据库的表初始化



4、settings为数据库连接和配置信息初始化

