package application

import (
	"context"
	"github.com/foxfurry/go_dining_hall/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

type IApp interface {
	Start()
	Shutdown(ctx context.Context)
}

type dinningApp struct {
	server *http.Server
}

func CreateApp() IApp {
	appHandler := gin.New()

	app := dinningApp{
		server: &http.Server{
			Addr:              viper.GetString("dinning_host"),
			Handler:           appHandler,
		},
	}

	return &app
}

func (d *dinningApp) initialize(){
	timeOut := time.Second * 1
	kitchenHost := viper.GetString("kitchen_host")
	var kitchenConn net.Conn
	var err error
	logger.LogMessageF("Trying to reach kitchenHost server on: %v", kitchenHost)

	for {
		kitchenConn, err = net.DialTimeout("tcp", kitchenHost, timeOut)
		if err == nil {
			break
		}
		logger.LogMessage("Could not reach kitchenHost, retrying")
	}
	kitchenConn.Close()

	logger.LogMessage("Successfully connected!")
}

func (d *dinningApp) Start() {
	d.initialize()
	logger.LogMessage("Starting dinning hall server")

	if err := d.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.LogPanicF("Unexpected error while running server: %v", err)
	}

}

func (d *dinningApp) Shutdown(ctx context.Context) {
	if err := d.server.Shutdown(ctx); err != nil {
		logger.LogPanicF("Unexpected error while closing server: %v", err)
	}
	logger.LogMessage("Server terminated successfully")
}
