package controller

import (
	"github.com/foxfurry/go_dining_hall/internal/dto"
	"github.com/foxfurry/go_dining_hall/internal/infrastructure/logger"
	"github.com/foxfurry/go_dining_hall/internal/service/supervisor"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

type IController interface {
	RegisterDiningRouter(c *gin.Engine)
	Initialize(menu int)
}

type DiningController struct {
	super supervisor.ISupervisor
}

func (ctrl *DiningController) Initialize(menu int) {
	ctrl.super.GenerateTables(viper.GetInt("table_num"), menu)
	ctrl.super.InitializeTables()
	ctrl.super.GenerateWaiter(viper.GetInt("waiter_num"))
	ctrl.super.StartWaiters()
}

func NewDiningController() IController {
	return &DiningController{super: &supervisor.DiningSupervisor{}}
}

func (ctrl *DiningController) distribution(c *gin.Context) {
	var data dto.Distribution

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Panic(err)
	}
	logger.LogWaiterF(data.WaiterID, "Table %d was served!", data.TableID)
	ctrl.super.FreeTable(data.TableID)
}
