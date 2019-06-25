package main

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/golang/sync/errgroup"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Eric-GreenComb/palletone/config"
	myrouter "github.com/Eric-GreenComb/palletone/router"
)

var (
	g errgroup.Group
)

func main() {

	if config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(Cors())

	/* api base */
	myrouter.SetupBaseRouter(router)

	// Palletone 调用接口
	myrouter.SetupPalletoneRouter(router)

	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	gracehttp.Serve(server)
}
