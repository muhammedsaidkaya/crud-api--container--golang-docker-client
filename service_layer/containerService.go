package service_layer

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"main/dto"
	"main/repository_layer"
	"strconv"
	"strings"
)

type containerService struct {
	repository repository_layer.ContainerRepositoryInterface
}

func NewContainerService(r repository_layer.ContainerRepositoryInterface) ContainerServiceInterface {
	return &containerService{repository: r}
}

func (_c containerService) GetAll(limit string) ([]types.Container, error) {
	_limit, _ := strconv.Atoi(limit)
	containers, err := _c.repository.GetAllByLimit(_limit)
	return containers, err
}

func (_c containerService) GetByID(id string) (types.Container, error) {
	return _c.repository.GetByID(id)
}

func (_c containerService) Create(newContainerInput dto.ContainerInput) (container.ContainerCreateCreatedBody, error) {
	split := strings.Split(newContainerInput.ExposePort, "/")
	exposePort, _ := nat.NewPort(split[1], split[0])
	containerConfig := &container.Config{
		Image: newContainerInput.Image,
		ExposedPorts: nat.PortSet{
			exposePort: struct{}{},
		},
	}
	containerHostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			exposePort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: newContainerInput.PublishPort,
				},
			},
		},
	}
	return _c.repository.Create(containerConfig, containerHostConfig, newContainerInput.Name)
}

func (_c containerService) Delete(id string) error {
	return _c.repository.DeleteByID(id)
}
