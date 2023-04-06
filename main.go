package main

import (
	"fmt"
	"io"
	"os"

	llog "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"service_template/app"
	"service_template/infra"
	"service_template/logger"
	"service_template/tracer"
)

func main() {
	viper.SetConfigFile("./config/config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		llog.Fatal("Config error", err)
	}

	port := viper.GetUint("port.stat_api")
	if port == 0 {
		port = 8101
	}
	bindHost := fmt.Sprintf(":%d", port)

	// log level
	logLevel := logger.Level(0)
	logLevel.Decode(viper.GetString("log.level"))

	// infra
	infraConfig := infra.Config{
		ServiceName: viper.GetString("infra.service_name"),
		Logger: &logger.Config{
			Level: logLevel,
		},
		Tracer: &tracer.Config{
			AgentAddress: viper.GetString("tracer.agent_address"),
			Sampler: &tracer.SamplerConfig{
				Type:  viper.GetString("tracer.sampler.type"),
				Param: viper.GetFloat64("tracer.sampler.param"),
			},
		},
	}

	// ctx
	var wr io.Writer = os.Stdout
	loggerPath := viper.GetString("log.output")
	if loggerPath != "stdout" {
		logFile, err := os.OpenFile(loggerPath+"ttmserver.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			llog.Fatalln("Error opening log file", err)
		}
		wr = io.MultiWriter(os.Stdout, logFile)
	}

	ictx := infra.Context(infraConfig, wr)

	srv := &app.App{Infra: infraConfig}
	srv.Initialize(ictx)
	srv.Run(ictx, bindHost)
}
