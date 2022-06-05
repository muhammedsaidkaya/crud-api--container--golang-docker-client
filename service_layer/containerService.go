package service_layer

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/dto"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/repository_layer"
	"regexp"
	"strconv"
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
	exposePortMatched, _ := regexp.MatchString(`^\d+$`, newContainerInput.ExposePort)
	publishPortMatched, _ := regexp.MatchString(`^\d+$`, newContainerInput.PublishPort)
	if exposePortMatched && publishPortMatched {
		exposePort, err := nat.NewPort("tcp", newContainerInput.ExposePort)
		if err != nil {
			return container.ContainerCreateCreatedBody{}, errors.New(fmt.Sprintf("Port is not appropriate: %v", err))
		}
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
	} else {
		return container.ContainerCreateCreatedBody{}, errors.New(fmt.Sprintf("Ports is not appropriate for ^0-9+$ regex"))
	}
}

func (_c containerService) Delete(id string) error {
	return _c.repository.DeleteByID(id)
}
