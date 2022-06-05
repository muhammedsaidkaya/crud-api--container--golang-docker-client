package controller_layer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/dto"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/logger"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/service_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/tracer"
	"github.com/peteprogrammer/go-automapper"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
	} else {
		var metadata dto.ContainerMetadata
		automapper.Map(container, &metadata)
		c.IndentedJSON(http.StatusOK, metadata)
	}
}

func (_c ContainerController) getAll(c *gin.Context) {
	limit := c.Query("limit")
	_, span := tracer.GetTracer().Start(c.Request.Context(), "getAll", oteltrace.WithAttributes(attribute.String("limit", limit)))
	defer span.End()

	containers, err := _c.service.GetAll(limit)
	if err != nil {
		logger.GetLogger().Panic(fmt.Sprintf("%v", err))
	}
	var metadataList []dto.ContainerMetadata
	for _, container := range containers {
		var metadata dto.ContainerMetadata
		automapper.Map(container, &metadata)
		metadataList = append(metadataList, metadata)
	}
	c.IndentedJSON(http.StatusOK, metadataList)
}

func (_c ContainerController) delete(c *gin.Context) {
	id := c.Param("id")
	err := _c.service.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "container not found")
	} else {
		c.IndentedJSON(http.StatusOK, "container removed")
	}
}

func (_c ContainerController) create(c *gin.Context) {
	var newContainerInput dto.ContainerInput
	if err := c.BindJSON(&newContainerInput); err != nil {
		return
	}
	container, err := _c.service.Create(newContainerInput)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}
	c.IndentedJSON(http.StatusCreated, container)
}
