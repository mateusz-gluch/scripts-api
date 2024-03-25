package configurators

import "os"

type Config struct {
	Host        string
	Environment string
	RootMessage string
	Version     string
	AssetsUrl   string
	MountPath   string
}

type EnvConfig struct {
}

func (cfg *EnvConfig) GetConfig() Config {
	return Config{
		Host:        os.Getenv("HOST"),
		Environment: os.Getenv("ENVIRONMENT"),
		RootMessage: os.Getenv("ROOT_MESSAGE"),
		Version:     os.Getenv("VERSION"),
		AssetsUrl:   os.Getenv("ASSETS_URL"),
		MountPath:   os.Getenv("MOUNT_PATH"),
	}
}
