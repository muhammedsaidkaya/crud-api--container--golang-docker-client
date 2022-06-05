package controller_layer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/dto"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/logger"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/service_layer"
	"github.com/peteprogrammer/go-automapper"
	"net/http"
)

type ContainerController struct {
	service service_layer.ContainerServiceInterface
}

func NewContainerController(s service_layer.ContainerServiceInterface) ContainerControllerInterface {
	controller := ContainerController{service: s}
	return &controller
}

func (_c ContainerController) Setup(router *gin.Engine) {
	router.GET("/containers", _c.getAll)
	router.GET("/containers/:id", _c.getByID)
	router.POST("/containers", _c.create)
	router.DELETE("/containers/:id", _c.delete)
}

func (_c ContainerController) getByID(c *gin.Context) {
	id := c.Param("id")
	container, err := _c.service.GetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"data": nil, "message": fmt.Sprintf("%v", err)})
	} else {
		var metadata dto.ContainerMetadata
		automapper.Map(container, &metadata)
		c.IndentedJSON(http.StatusOK, gin.H{"data": metadata, "message": "Container found"})
	}
}

func (_c ContainerController) getAll(c *gin.Context) {
	limit := c.Query("limit")

	/*tr := otel.Tracer("gin-gonic")
	_, span := tr.Start(c.Request.Context(), "controller")
	span.SetAttributes(attribute.Key("limit").String(limit))
	span.AddEvent()
	defer span.End()*/

	containers, err := _c.service.GetAll(limit)
	if err != nil {
		logger.GetLogger().Panic(fmt.Sprintf("%v", err))
	}
	metadataList := make([]dto.ContainerMetadata, 0)
	for _, container := range containers {
		var metadata dto.ContainerMetadata
		automapper.Map(container, &metadata)
		metadataList = append(metadataList, metadata)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": metadataList})
}

func (_c ContainerController) delete(c *gin.Context) {
	id := c.Param("id")
	err := _c.service.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"data": id, "message": "Container not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": id, "message": "Container removed"})
	}
}

func (_c ContainerController) create(c *gin.Context) {
	var newContainerInput dto.ContainerInput
	if err := c.BindJSON(&newContainerInput); err != nil {
		return
	}
	container, err := _c.service.Create(newContainerInput)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"data": nil, "message": fmt.Sprintf("%v", err)})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"data": container, "message": "Container created"})
}
