package main

import (
	"os"
	"os/signal"

	"github-release-puller/config"
	"github-release-puller/executor"

	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: github-release-puller <path/to/configuration>")

		return
	}

	configuration, err := os.Open(os.Args[1])
	if err != nil {
		println("Open configuration:", err.Error())

		return
	}

	defer configuration.Close()

	cfg := &config.Config{}

	if err := yaml.NewDecoder(configuration).Decode(cfg); err != nil {
		println("Parse configuration:", err.Error())

		return
	}

	executor.Start(cfg)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	<-signalCh
}
