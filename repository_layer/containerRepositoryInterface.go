package repository_layer

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type ContainerRepositoryInterface interface {
	GetAllByLimit(limit int) ([]types.Container, error)
	GetByID(id string) (types.Container, error)
	DeleteByID(id string) error
	Create(config *container.Config, hostConfig *container.HostConfig, name string) (container.ContainerCreateCreatedBody, error)
}
