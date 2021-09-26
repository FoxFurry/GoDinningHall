package controller

import "github.com/gin-gonic/gin"

func (d *DiningController) RegisterDiningRouter(c *gin.Engine) {
	c.POST("/distribution", d.distribution)
}
