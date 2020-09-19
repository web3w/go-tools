package main

import (
	gen "github.com/gisvr/go-tools/kratos-protoc/generator"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "protc"
	app.Usage = "protobuf生成工具"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "bm",
			Usage:       "whether to use BM for generation",
			Destination: &gen.WithBM,
		},
		cli.BoolFlag{
			Name:        "grpc",
			Usage:       "whether to use gRPC for generation",
			Destination: &gen.WithGRPC,
		},
		cli.BoolFlag{
			Name:        "swagger",
			Usage:       "whether to use swagger for generation",
			Destination: &gen.WithSwagger,
		},
		cli.BoolFlag{
			Name:        "ecode",
			Usage:       "whether to use ecode for generation",
			Destination: &gen.WithEcode,
		},
	}
	app.Action = func(c *cli.Context) error {
		return gen.ProtocAction(c)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
