package environment

import "github.com/kelseyhightower/envconfig"

const Name = "go-web-scaffold"
const Version = "development"

type Config struct {
	EnableTracesStdout bool `default:"true" split_words:"true"`
}

func Parse() Config {
	var conf Config

	envconfig.MustProcess("go-web-scaffold", &conf)

	return conf
}
