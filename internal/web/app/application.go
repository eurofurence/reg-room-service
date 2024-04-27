package app

import (
	"fmt"
	"log"

	"github.com/eurofurence/reg-room-service/internal/config"
)

type Params struct {
	configFilePath string
	migrateDB      bool
}

func NewParams(configFile string, migrateDB bool) Params {
	return Params{
		configFilePath: configFile,
		migrateDB:      migrateDB,
	}
}

type Application struct {
	Params Params
}

func New(params Params) *Application {
	return &Application{
		Params: params,
	}
}

func (a *Application) Run() error {
	conf, err := config.UnmarshalFromYamlConfiguration(a.Params.configFilePath)
	if err != nil {
		// this is so early we don't know what log format to set up anyway
		log.Fatal(fmt.Sprintf("%-5s [%s] %s", "ERROR", "00000000", fmt.Sprintf("Error while loading configuration: %v", err)))
	}

	fmt.Println("use config to setup db and business logic", *conf)

	// TODO(noroth) construct types to and pass to the server

	// TODO(noroth) start server
	// if err := runServerWithGracefulShutdown(); err != nil {
	// 	return err
	// }

	return nil
}
