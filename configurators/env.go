package configurators

import (
	"os"

	"github.com/elmodis/go-libs/api"
)

type Config struct {
	Host        string
	Environment string
	Misc        *api.MiscConfig
	AssetsUrl   string
	MountPath   string
}

type EnvConfig struct {
}

func (cfg *EnvConfig) GetConfig() Config {
	return Config{
		Host:        os.Getenv("HOST"),
		Environment: os.Getenv("ENVIRONMENT"),

		Misc: &api.MiscConfig{
			RootMessage: os.Getenv("ROOT_MESSAGE"),
			Version:     os.Getenv("VERSION"),
			PingValue:   1,
		},

		AssetsUrl: os.Getenv("ASSETS_URL"),
		MountPath: os.Getenv("MOUNT_PATH"),
	}
}
