package main

import (
	"scripts-api/configurators"
	ctrl "scripts-api/controllers"
	"scripts-api/docs"
	"scripts-api/handlers"
	"time"

	log "github.com/elmodis/go-libs/api/logging"
	"github.com/elmodis/go-libs/api/monitoring"
	"github.com/elmodis/go-libs/caches"
	cli "github.com/elmodis/go-libs/clients"
	"github.com/elmodis/go-libs/fileengines"
	"github.com/elmodis/go-libs/models/properties"
	"github.com/elmodis/go-libs/parsers"
	"github.com/elmodis/go-libs/repositories"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	swaggergin "github.com/swaggo/gin-swagger"
)

var engine, metrics *gin.Engine

var data *ctrl.ScriptDataController
var misc *ctrl.MiscController

var (
	assetCacheExpiration = 30 * time.Minute
	eventCategories      = []string{
		"machine", "data", "diagnostics", "maintenance",
		"system", "anomaly", "ai_maintenance", "ai_energy"}
)

func init() {

	c := &configurators.EnvConfig{}
	cfg := c.GetConfig()

	logger := log.Configure(cfg.Environment)

	if cfg.Environment == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.New()
	engine.Use(gin.Recovery())
	engine.Use(log.Middleware(logger))
	engine.Use(monitoring.Middleware("scripts"))

	metrics = gin.New()
	metrics.Use(gin.Recovery())

	assetsApi := cli.NewApiClient(cfg.AssetsUrl, "assets-api", logger)
	assetsCache := caches.NewExpirationCache(new(any), "assets-cache", assetCacheExpiration, logger)
	assetsCApi := cli.NewCacheClient(assetsApi, assetsCache, "assets", logger)

	assets := repositories.NewPropertiesRepo(properties.Asset{}, assetsCApi, "assets/%s", nil, "assets", logger)

	scriptData := &repositories.ScriptDataRepository{
		Engine: &fileengines.CSVEngine{RootDir: cfg.MountPath},
	}

	filterParser := map[string]parsers.Parser[[]string]{
		"category": parsers.NewSequenceEnumParser(eventCategories, "category", logger),
	}

	misc = ctrl.NewMiscController(cfg.Misc)
	data = ctrl.NewScriptDataController(scriptData, assets, filterParser, logger)

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
	engine.GET("/events-summary/data", handlers.EventsSummary(data))
	engine.GET("/online-summary/data", handlers.OnlineSummary(data))

	go metrics.Run(":8081")
	engine.Run(":8080")
}
