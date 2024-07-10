package configurators

import (
	"os"
	"strconv"

	"github.com/elmodis/go-libs/api"
)

type Config struct {
	Host        string
	BasePath    string
	Environment string
	Misc        *api.MiscConfig
	AssetsUrl   string
	MountPath   string
	Postgres    *api.DatabaseConfig
	EventsTable string
	OnlineTable string
}

type EnvConfig struct {
}

func (cfg *EnvConfig) GetConfig() Config {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	return Config{
		Host:        os.Getenv("HOST"),
		BasePath:    os.Getenv("BASE_PATH"),
		Environment: os.Getenv("ENVIRONMENT"),

		Misc: &api.MiscConfig{
			RootMessage: os.Getenv("ROOT_MESSAGE"),
			Version:     os.Getenv("VERSION"),
			PingValue:   1,
		},

		AssetsUrl: os.Getenv("ASSETS_URL"),
		MountPath: os.Getenv("MOUNT_PATH"),

		Postgres: &api.DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     port,
			Login:    os.Getenv("POSTGRES_LOGIN"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
			DBtype:   "postgres",
		},

		EventsTable: os.Getenv("EVENTS_TABLE"),
		OnlineTable: os.Getenv("ONLINE_TABLE"),
	}
}
