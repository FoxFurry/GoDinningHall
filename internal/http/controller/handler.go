package controller

import (
	"github.com/foxfurry/go_dining_hall/internal/dto"
	"github.com/foxfurry/go_dining_hall/internal/service/supervisor"
	"github.com/gin-gonic/gin"
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
	ctrl.super.GenerateTables(3, menu)
	ctrl.super.InitializeTables()
	ctrl.super.GenerateWaiter(3)
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
	log.Printf("Freeing table %v", data.TableID)
	ctrl.super.FreeTable(data.TableID)
}
