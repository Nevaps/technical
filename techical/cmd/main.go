package main

import (
	"techical/internal/app"
	"techical/internal/config"
	"techical/pkg/logger"

	_ "techical/docs"
)

// @title           CurrencyRate API
// @version         1.0
// @description     This is a service that helps to calculate get rates
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	log := logger.NewLogger()
	appConfig := config.MustLoad()

	application := app.NewApp(log, appConfig)
	application.Run()

}
