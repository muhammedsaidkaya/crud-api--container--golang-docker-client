package service_layer

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"main/dto"
)

type ContainerServiceInterface interface {
	GetAll(limit string) ([]types.Container, error)
	GetByID(id string) (types.Container, error)
	Create(newContainerInput dto.ContainerInput) (container.ContainerCreateCreatedBody, error)
	Delete(id string) error
}
