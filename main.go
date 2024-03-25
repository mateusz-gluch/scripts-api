package main

import (
	"scripts-api/configurators"
	ctrl "scripts-api/controllers"
	"scripts-api/dataengine"
	"scripts-api/docs"
	"scripts-api/handlers"
	"scripts-api/repos"
	"time"

	log "github.com/elmodis/go-libs/api/logging"
	"github.com/elmodis/go-libs/api/monitoring"
	"github.com/elmodis/go-libs/caches"
	cli "github.com/elmodis/go-libs/clients"
	"github.com/elmodis/go-libs/models/enums"
	"github.com/elmodis/go-libs/models/properties"
	"github.com/elmodis/go-libs/parsers"
	repo "github.com/elmodis/go-libs/repositories"
	"github.com/elmodis/go-libs/validators"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	swaggergin "github.com/swaggo/gin-swagger"
)

var engine, metrics *gin.Engine

var data ctrl.ScriptDataController
var misc ctrl.MiscController

func init() {

	c := &configurators.EnvConfig{}
	cfg := c.GetConfig()

	log.Configure()

	if cfg.Environment == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.New()
	engine.Use(gin.Recovery())
	engine.Use(log.Middleware())
	engine.Use(monitoring.Middleware("scripts"))

	metrics = gin.New()
	metrics.Use(gin.Recovery())

	assetsApi := &cli.ApiClient{
		Kind: "Internal",
		Url:  cfg.AssetsUrl,
	}

	exp, _ := time.ParseDuration("30m")
	assetsCache, err := caches.NewExpKVCache[any](caches.Config{
		"duration": exp,
		"label":    "Assets",
	})

	if err != nil {
		logrus.WithField("Error", err).Fatal("Could not initialize Assets Cache")
	}

	assetsCachingApi := &cli.CachingClient{
		Remote: assetsApi,
		Local:  assetsCache,
	}

	assets := &repo.PropertiesAPIRepo[properties.Asset]{
		Engine:   assetsCachingApi,
		Label:    "assets",
		PathSpec: "assets/%s",
	}

	scriptData := &repos.ScriptDataRepository{
		Engine: &dataengine.CSVEngine{RootDir: cfg.MountPath},
	}

	filterParser := map[string]parsers.Parser[[]string]{
		"category": &parsers.SequenceParser{
			Valid:     &validators.EnumValidator{Enum: enums.Category{}},
			Separator: ",",
		},
	}

	assetParser := &parsers.SequenceParser{
		Valid:     &validators.NumericIdValidator{},
		Separator: ",",
	}

	misc = ctrl.MiscController{
		RootMessage: cfg.RootMessage,
		PingValue:   1,
		Version:     cfg.Version,
	}

	data = ctrl.ScriptDataController{
		ScriptRepo:        scriptData,
		AssetRepo:         assets,
		Filter:            filterParser,
		AssetParser:       assetParser,
		OrganizationValid: &validators.NumericIdValidator{},
		Timestamp:         &parsers.TimestampParser{TsValid: &validators.TimestampValidator{}},
	}

	docs.SwaggerInfo.Host = cfg.Host
	docs.SwaggerInfo.BasePath = "/scripts"
}

// @title						Internal Scripts API
// @version					1.0
//
// @contact.name				Elmodis
// @contact.email				mateusz-gluch@elmodis.com
//
// @host						dev-internal-api.elmodis.com
// @BasePath					/
//
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	// metrics
	metrics.GET("/metrics", monitoring.PromHandler())

	// misc and docs
	engine.GET("/docs/*any", swaggergin.WrapHandler(swaggerfiles.Handler))
	engine.GET("/ping", misc.Ping())
	engine.GET("/version", misc.Ver())
	engine.GET("/", misc.Root())

	// counts
	engine.GET("/events-summary/data", handlers.EventsSummary(&data))
	engine.GET("/online-summary/data", handlers.OnlineSummary(&data))

	go metrics.Run(":8081")
	engine.Run(":8080")
}
