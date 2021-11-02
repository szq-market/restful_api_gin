module gin_project

go 1.16

require (
	github.com/RichardKnop/machinery v1.10.6 //用于异步处理，相当于celery
	github.com/gin-contrib/cors v1.3.1 //用于用户鉴权
	github.com/gin-gonic/gin v1.7.4
	github.com/google/uuid v1.2.0
	github.com/knadh/koanf v1.2.3 //支持各种格式的文本解析
	github.com/opentracing/opentracing-go v1.2.0
	github.com/urfave/cli v1.22.5
	go.uber.org/fx v1.14.2
	go.uber.org/zap v1.19.1 //用于日志输出
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.12
)
