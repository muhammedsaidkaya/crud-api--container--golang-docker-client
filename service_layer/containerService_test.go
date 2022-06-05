package service_layer

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/dto"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/mock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestContainerService(t *testing.T) {

	mockRepository := mock.NewMockContainerRepositoryInterface(gomock.NewController(t))

	t.Run("deleteByID method tests", func(t *testing.T) {
		t.Run("should delegate repository deleteById func", func(t *testing.T) {
			id := "dummy"
			mockRepository.EXPECT().
				DeleteByID(id).
				Return(nil).
				Times(1)

			service := NewContainerService(mockRepository)
			err := service.Delete(id)

			assert.Nil(t, err)
		})

		t.Run("should not delegate repository deleteById func", func(t *testing.T) {
			id := "dummy"
			_error := errors.New("container not found")
			mockRepository.EXPECT().
				DeleteByID(id).
				Return(_error).
				Times(1)

			service := NewContainerService(mockRepository)
			err := service.Delete(id)

			assert.Equal(t, err.Error(), _error.Error())
		})
	})

	t.Run("getAll method tests", func(t *testing.T) {
		t.Run("should delegate repository getAll func", func(t *testing.T) {
			limit := "1"
			_limit, _ := strconv.Atoi(limit)
			var containers []types.Container
			containers = append(containers, types.Container{})
			mockRepository.EXPECT().
				GetAllByLimit(_limit).
				Return(containers, nil).
				Times(1)

			service := NewContainerService(mockRepository)
			_containers, err := service.GetAll(limit)

			assert.Equal(t, len(containers), len(_containers))
			assert.Nil(t, err)
		})

		t.Run("should not delegate repository getAll func", func(t *testing.T) {
			limit := "1"
			_limit, _ := strconv.Atoi(limit)
			mockRepository.EXPECT().
				GetAllByLimit(_limit).
				Return(nil, errors.New("couldn't connect to docker daemon socket")).
				Times(1)

			service := NewContainerService(mockRepository)
			_, err := service.GetAll(limit)

			assert.NotNil(t, err)
		})
	})

	t.Run("getById method tests", func(t *testing.T) {
		t.Run("should delegate repository getByID func", func(t *testing.T) {
			id := "dummy"
			mockRepository.EXPECT().
				GetByID(id).
				Return(types.Container{ID: id}, nil).
				Times(1)

			service := NewContainerService(mockRepository)
			container, err := service.GetByID(id)

			assert.Nil(t, err)
			assert.Equal(t, container.ID, id)
		})

		t.Run("should not delegate repository getByID func", func(t *testing.T) {
			id := "dummy"
			_error := errors.New("container not found")
			mockRepository.EXPECT().
				GetByID(id).
				Return(types.Container{}, _error).
				Times(1)

			service := NewContainerService(mockRepository)
			_, err := service.GetByID(id)

			assert.Equal(t, err.Error(), _error.Error())
		})
	})

	t.Run("create method tests", func(t *testing.T) {
		t.Run("should not delegate repository create func", func(t *testing.T) {

			newContainerInput := dto.ContainerInput{
				Name:        "uzumlukek-test",
				Image:       "uzumlukek/docker-client",
				ExposePort:  "3000a",
				PublishPort: "3000",
			}
			service := NewContainerService(mockRepository)
			_, err := service.Create(newContainerInput)

			assert.ErrorContains(t, err, "Ports is not appropriate")
		})

		t.Run("should not delegate repository create func", func(t *testing.T) {

			newContainerInput := dto.ContainerInput{
				Name:        "uzumlukek-test",
				Image:       "uzumlukek/docker-client",
				ExposePort:  "3000",
				PublishPort: "3000",
			}
			mockRepository.EXPECT().
				Create(gomock.Any(), gomock.Any(), "uzumlukek-test").
				Return(container.ContainerCreateCreatedBody{}, nil).
				Times(1)

			service := NewContainerService(mockRepository)
			_, err := service.Create(newContainerInput)

			assert.Nil(t, err)
		})

		t.Run("should not delegate repository create func", func(t *testing.T) {

			newContainerInput := dto.ContainerInput{
				Name:        "uzumlukek-test",
				Image:       "uzumlukek/docker-client",
				ExposePort:  "3000",
				PublishPort: "3000",
			}
			mockRepository.EXPECT().
				Create(gomock.Any(), gomock.Any(), "uzumlukek-test").
				Return(container.ContainerCreateCreatedBody{}, errors.New("start error")).
				Times(1)

			service := NewContainerService(mockRepository)
			_, err := service.Create(newContainerInput)

			assert.NotNil(t, err)
		})
	})
}
