package main

import (
	"flag"
	"fmt"
	"github.com/iamsayantan/larasockets/app_managers"
	"github.com/iamsayantan/larasockets/channel_managers"
	"github.com/iamsayantan/larasockets/config"
	"github.com/iamsayantan/larasockets/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

var (
	defaultConfigLocation = "."
)

func main() {
	configPath := flag.String("config", defaultConfigLocation, "Configuration file path")
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("server.port", "8005")

	viper.AddConfigPath(*configPath)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err.Error())
	}

	var larasocketConfig config.LarasocketsConfig
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("error reading configuration file", zap.String("error", err.Error()))
		return
	}

	err = viper.Unmarshal(&larasocketConfig)
	if err != nil {
		logger.Fatal("error unmarshalling configuration file", zap.String("error", err.Error()))
		return
	}

	if err := larasocketConfig.Validate(); err != nil {
		logger.Fatal("error validating the configuration file", zap.String("error", err.Error()))
		return
	}

	appManager := app_managers.NewConfigManager(larasocketConfig.Apps)
	channelManager := channel_managers.NewLocalManager(appManager, logger)

	srv := server.NewServer(logger, channelManager)
	logger.Info("starting larasockets server", zap.String("port", larasocketConfig.Server.Port))

	if larasocketConfig.Server.TLS {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%s", larasocketConfig.Server.Port),
			larasocketConfig.Server.Certificate,
			larasocketConfig.Server.Key,
			srv,
		)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", larasocketConfig.Server.Port), srv)
	}

	if err != nil {
		logger.Fatal("error starting server", zap.String("error", err.Error()))
	}
}
