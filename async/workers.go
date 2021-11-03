package async

import (
	"context"
	"fmt"
	"github.com/RichardKnop/machinery/example/tracers"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"

)

func startServer() (*machinery.Server,error) {
	//读取配置文件
	cnf,err := config.NewFromYaml("./config.yml",false)
	if err != nil{
		log.ERROR.Println("config failed ",err)
	}
	//根据配置文件启动
	server,err := machinery.NewServer(cnf)
	if err != nil{
		return nil,err
	}

	//注册tasks
	tasks := map[string]interface{}{
		"sum": Sum,
	}
	return server,server.RegisterTasks(tasks)
}


func Worker() error {
	//消费者名称
	consumerTag := "machineryDemo"
	//跟踪的worker执行成功还是失败
	cleanup,err := tracers.SetupTracer(consumerTag)
	//如果有问题打印出错误日志
	if err != nil{
		log.FATAL.Fatalln(err)
	}
	defer cleanup()

	server,err := startServer()

	if err != nil{

	}

	//告诉machinery我们work的并发是多少
	worker := server.NewWorker(consumerTag,1)

	//运行中的错误日志
	errorhandler := func(err error) {
		log.ERROR.Println("err handler",err)
	}
	//运行前的错误日志
	pretaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("task hanlder for:",signature.Name)
	}
	//运行结尾的错误日志
	posttaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("task end hanlder for:",signature.Name)
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)

	return worker.Launch()
}


//耗时长的任务丢到work
func Send() error {
	cleanup,err := tracers.SetupTracer("sender")
	if err != nil{
		log.FATAL.Fatalln(err)
	}
	defer cleanup()

	server , err := startServer()

	if err != nil{
		return err
	}

	var (
		addTask tasks.Signature
	)

	var initTasks = func() {
		addTask = tasks.Signature{
			//指明哪一个去跑work
			Name: "sum",
			Args: []tasks.Arg{
				{
					Type: "[]int64",
					Value: []int{1,2,3,4,5,6},
				},
			},
		}
	}
	//opentracing分布式多链路的库，用下面的uuid和每个进程绑定
	span,ctx := opentracing.StartSpanFromContext(context.Background(),"send")
	defer span.Finish()

	batchId := uuid.New().String()
	span.SetBaggageItem("batch.id",batchId)
	//这里绑定
	span.LogFields(opentracing_log.String("batch.id",batchId))

	log.INFO.Println("starting batch:",batchId)
	//初始化tasks
	initTasks()
	//真正开始执行任务
	asyncResult,err := server.SendTaskWithContext(ctx,&addTask)

	//asyncResult.Get()  //获取异步结果
	if err != nil{
		return fmt.Errorf("not tasks",err)
	}
	log.INFO.Println(asyncResult)
	return nil
}