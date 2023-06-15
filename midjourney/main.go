package main

import (
	"fmt"
	"midjourney/handlers"
	"midjourney/initialization"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	cfg := pflag.StringP("config", "c", "./config.yaml", "api server config file path.")

	pflag.Parse()

	initialization.LoadConfig(*cfg)
	initialization.LoadDiscordClient(handlers.DiscordMsgCreate, handlers.DiscordMsgUpdate)

	ret, err := handlers.UploadImage("/tmp/tangcheng_futuristic_city_c8cf7f74691fb7543f2baa1055edf889_149bd28d-f5c1-492c-89dc-f58d1245a083.png", "/attachments/1112690540476104734/1118861394272583820/tangcheng_futuristic_city_c8cf7f74691fb7543f2baa1055edf889_149bd28d-f5c1-492c-89dc-f58d1245a083.png")
	fmt.Sprintf("%#v %#v", ret, err)

	r := gin.Default()

	r.POST("/v1/trigger/midjourney-bot", handlers.MidjourneyBot)
	r.POST("/v1/trigger/upload", handlers.UploadFile)

	r.Run(fmt.Sprintf(":%s", initialization.GetConfig().MJ_PORT))
}
