/*
Package main gin-admin

Swagger Reference: https://github.com/swaggo/swag#declarative-comments-format

Usage：

	go get -u github.com/swaggo/swag/cmd/swag
	swag init --generalInfo ./cmd/gin-admin/main.go --output ./internal/app/swagger */
package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/LyricTian/gin-admin/v8/internal/app"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

// VERSION Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "8.1.0"

// @title gin-admin
// @version 8.1.0
// @description RBAC scaffolding based on GIN + GORM + CASBIN + WIRE.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LyricTian
// @contact.email tiannianshou@gmail.com
func main() {
	ctx := logger.NewTagContext(context.Background(), "__main__")

	// 创建CLI
	app := cli.NewApp()
	app.Name = "gin-admin"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on GIN + GORM + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	// 提取参数
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "Casbin model configuration(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "Initialize menu's data configuration(.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
			},
		},
		// 将指定的文件映射根据不同的set方法映射进去
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetModelFile(c.String("model")),
				app.SetWWWDir(c.String("www")),
				app.SetMenuFile(c.String("menu")),
				app.SetVersion(VERSION))
		},
	}
}
