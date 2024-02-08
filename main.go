package main

import (
	"context"

	"github.com/talgat-ruby/interactive-comments-api/cmd/api"
	"github.com/talgat-ruby/interactive-comments-api/cmd/db"
	"github.com/talgat-ruby/interactive-comments-api/configs"
	"github.com/talgat-ruby/interactive-comments-api/internal/constant"
	"github.com/talgat-ruby/interactive-comments-api/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// conf
	conf, err := configs.NewConfig(ctx)
	if err != nil {
		panic(err)
	}

	// setup logger
	log := logger.New(conf.Env == constant.EnvironmentLocal)

	// configure db service
	d, err := db.New(log.With("service", constant.DB), conf.DB)
	if err != nil {
		log.ErrorContext(
			ctx,
			"initialize service error",
			"service", "database",
			"error", err,
		)
		panic(err)
	}
	log.InfoContext(ctx, "initialize service", "service", "database")

	// configure gateway service
	srv := api.New(log.With("service", constant.Api), conf.Api)
	log.InfoContext(ctx, "initialize service", "service", "api")
	// start gateway service
	srv.Start(ctx, cancel, d)

	<-ctx.Done()
	// Your cleanup tasks go here
	log.InfoContext(ctx, "cleaning up ...")

	log.InfoContext(ctx, "server was successful shutdown.")
}
