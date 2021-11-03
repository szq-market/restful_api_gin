package main

import (
	"context"
	"fmt"
	"gin_project/async"
	views "gin_project/gallery/api"
	usersDB "gin_project/gallery/models"
	"gin_project/gallery/settings"
	"gin_project/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli" //处理命令行的工具
	"go.uber.org/fx"
	"net/http"
	"os"
)

var (
	app *cli.App
)


func init() {
	//处理命令行的工具cli
	app = cli.NewApp()
	//初始化信息
	app.Name = "lufflyweb"
	app.Usage = "Gin rest demo"
	app.Version = "0.0.0"
}




func loadConfig() (*settings.Config,error) {
	return settings.Load()
}

func newServer(lc fx.Lifecycle,cfg *settings.Config) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(),gin.Recovery(),cors.Default())

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d",cfg.ServerConfig.Port),
		Handler: r,
	}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logging.DefaultLogger().Infof("start to reset api server: %d",cfg)
				go srv.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logging.DefaultLogger().Infof("Stopped the api server")
				return srv.Shutdown(ctx)
			},
		})
	return r
}

func printAppInfo(cfg *settings.Config)  {
	logging.DefaultLogger().Infow("app info","config",cfg)
}


func runApplication()  {
	//setup app + run server
	app := fx.New(
		fx.Provide(
			loadConfig,
			settings.NewDatabase,
			usersDB.NewUsersDB,
			views.NewHandler,
			//gin server
			newServer,
			),
			fx.Invoke(
				views.RouteV1,
				printAppInfo,
				),
			)
	app.Run()
}

func main()  {
	//根据参数做出不同的动作
	app.Commands = [] cli.Command{
		{
			Name: "server",
			Usage: "launch Gin Server By boyleGu",
			Action: func(c *cli.Context) error {
				runApplication()
				return nil
			},
		},
		{
			Name: "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				if err := async.Worker();err != nil{
					return cli.NewExitError(err.Error(),1)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)

	//runApplication()
}