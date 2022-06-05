package controller_layer

import (
	"github.com/gin-gonic/gin"
)

type ContainerControllerInterface interface {
	Setup(router *gin.Engine)
}
