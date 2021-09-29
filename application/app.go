package application

import (
	"context"
	"encoding/json"
	"github.com/foxfurry/go_dining_hall/internal/dto"
	"github.com/foxfurry/go_dining_hall/internal/http/controller"
	"github.com/foxfurry/go_dining_hall/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
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
	appHandler := gin.Default()

	ctrl := controller.NewDiningController()
	ctrl.RegisterDiningRouter(appHandler)

	app := dinningApp{
		server: &http.Server{
			Addr:    viper.GetString("dinning_host"),
			Handler: appHandler,
		},
	}

	count := app.initialize()
	ctrl.Initialize(count)

	return &app
}

func (d *dinningApp) initialize() int {
	timeOut := time.Second * 5
	kitchenHost := viper.GetString("kitchen_host")
	var err error
	logger.LogMessageF("Trying to reach kitchenHost server on: %v", kitchenHost)

	connClient := http.Client{
		Timeout: timeOut,
	}
	_, err = connClient.Get(kitchenHost)
	if err != nil {
		logger.LogPanic(err.Error())
	}

	req, err := http.Get(kitchenHost + "/menu")
	if err != nil {
		logger.LogPanic(err.Error())
	}
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.LogPanic(err.Error())
	}
	bodyMenu := dto.Menu{}

	if err = json.Unmarshal(bodyData, &bodyMenu); err != nil {
		logger.LogPanic(err.Error())
	}

	logger.LogMessage("Successfully connected!")

	return bodyMenu.ItemsCount
}

func (d *dinningApp) Start() {
	logger.LogMessage("Starting dinning hall server!")

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
