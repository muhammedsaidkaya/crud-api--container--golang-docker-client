package repository_layer

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/cache_layer"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/logger"
	"strings"
)

type containerRepository struct {
	cache        cache_layer.CacheInterface
	dockerClient *client.Client
}

func NewContainerRepository(c cache_layer.CacheInterface, dockerClient *client.Client) ContainerRepositoryInterface {
	return &containerRepository{cache: c, dockerClient: dockerClient}
}

func (_c containerRepository) Create(config *container.Config, hostConfig *container.HostConfig, name string) (container.ContainerCreateCreatedBody, error) {
	container, err := _c.dockerClient.ContainerCreate(context.Background(), config, hostConfig, nil, nil, name)
	if err != nil {
		return container, err
	}
	err = _c.dockerClient.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		return container, errors.New(fmt.Sprintf("start error: %v \n", err))
	}
	_c.cache.SetCache(container.ID, container)
	return container, nil
}

func (_c containerRepository) DeleteByID(id string) error {
	err := _c.dockerClient.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true})
	if err == nil {
		_c.cache.DeleteFromCache(id)
	}
	return err
}

func (_c containerRepository) GetByID(id string) (types.Container, error) {
	container, found := _c.cache.GetCache(id)
	if found == false {
		_container, err := _c.getByID(id)
		if err == nil {
			_c.cache.SetCache(id, _container)
		}
		return _container, err
	}
	logger.GetLogger().Info(fmt.Sprintf("Getting from cache: %v", id))
	return container.(types.Container), nil
}

func (_c containerRepository) getByID(id string) (types.Container, error) {
	containerList, _ := _c.getAll()
	for _, cont := range containerList {
		if strings.EqualFold(cont.ID, id) {
			return cont, nil
		}
	}
	return types.Container{}, errors.New("container not found")
}

func (_c containerRepository) GetAllByLimit(limit int) ([]types.Container, error) {
	return _c.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{Limit: limit})
}

func (_c containerRepository) getAll() ([]types.Container, error) {
	return _c.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
}
