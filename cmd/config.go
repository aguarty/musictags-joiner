package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
		Port string `envconfig:"SERVER_PORT" default:"8090"`
	}
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFile  string `envconfig:"LOG_FILE" default:""`
	Apikey   string `envconfig:"API_KEY" default:""`
}

func (a *application) initConfig(app string) error {
	var conf config
	if err := envconfig.Process(app, &conf); err != nil {
		envconfig.Usage(app, &conf)
		return err
	}
	a.cfg = conf
	return nil
}
