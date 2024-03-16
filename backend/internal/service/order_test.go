package service_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/config"
	mocks "cloud-render/internal/mocks/repository"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"cloud-render/internal/service"
	"database/sql"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)
	mockConfig := &config.Config{
		OrderStatuses: config.OrderStatuses{
			InQueue: "InQueue",
		},
		Redis: config.Redis{
			QueueName: "queue",
		},
	}
	redisClient, mockRedis := redismock.NewClientMock()

	orderService := service.NewOrderService(mockOrderProvider, nil, nil, "", "", mockConfig, redisClient)

	t.Run("Success", func(t *testing.T) {
		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		dto := dto.CreateOrderDTO{
			UserId:     123,
			Format:     "jpg",
			Resolution: "1920x1080",
			File:       file,
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
		}

		mockOrderProvider.EXPECT().Create(gomock.Any()).Return(nil)
		mockRedis.ExpectRPush(mockConfig.Redis.QueueName, gomock.Any())

		err = orderService.CreateOrder(dto)
		assert.NoError(t, err)
	})

	t.Run("OrderProviderError", func(t *testing.T) {
		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		dto := dto.CreateOrderDTO{
			UserId:     123,
			Format:     "jpg",
			Resolution: "1920x1080",
			File:       file,
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
		}

		mockOrderProvider.EXPECT().Create(gomock.Any()).Return(errors.New("error"))

		err = orderService.CreateOrder(dto)
		assert.Error(t, err)
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)

	m := make(map[string]int64, 1)
	m["InQueue"] = 1

	orderService := service.NewOrderService(mockOrderProvider, m, nil, "", "", nil, nil)

	t.Run("Success", func(t *testing.T) {
		dto := dto.UpdateOrderStatusDTO{
			UserId:      int64(123),
			StoringName: "test.jpg",
			Status:      "InQueue",
		}

		mockOrderProvider.EXPECT().UpdateStatus("test.jpg", int64(123), int64(1)).Return(nil)

		err := orderService.UpdateOrderStatus(dto)
		assert.NoError(t, err)
	})

	t.Run("OrderProviderError", func(t *testing.T) {
		mockOrderProvider.EXPECT().UpdateStatus("test.jpg", int64(123), int64(1)).Return(errors.New("error"))

		dto := dto.UpdateOrderStatusDTO{
			UserId:      int64(123),
			StoringName: "test.jpg",
			Status:      "InQueue",
		}

		err := orderService.UpdateOrderStatus(dto)
		assert.Error(t, err)
	})
}

func TestOrderService_UpdateOrderImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)
	mockConfig := &config.Config{
		OrderStatuses: config.OrderStatuses{
			Success: "Success",
		},
		HTTPServer: config.HTTPServer{
			Host: "localhost",
			Port: 8080,
		},
	}
	redisClient, _ := redismock.NewClientMock()

	orderService := service.NewOrderService(mockOrderProvider, nil, nil, "", "", mockConfig, redisClient)

	t.Run("Success", func(t *testing.T) {
		f, err := os.OpenFile("123/test.blend", os.O_WRONLY|os.O_CREATE, 0666)
		require.NoError(t, err)
		require.NotNil(t, f)
		f.Close()

		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove("123/test.jpg")

		dto := dto.UpdateOrderImageDTO{
			UserId: "123",
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
			File: file,
		}

		mockOrderProvider.EXPECT().UpdateStatus("test.blend", int64(123), gomock.Any()).Return(nil)
		mockOrderProvider.EXPECT().UpdateDownloadLink(int64(123), "test.blend", "http://localhost:8080/123/image/download/test.jpg").Return(nil)

		err = orderService.UpdateOrderImage(dto)
		assert.NoError(t, err)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove("123/test.jpg")

		dto := dto.UpdateOrderImageDTO{
			UserId: "invalid",
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
			File: file,
		}

		err = orderService.UpdateOrderImage(dto)
		assert.Error(t, err)
	})

	t.Run("UpdateStatusError", func(t *testing.T) {
		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove("123/test.jpg")

		dto := dto.UpdateOrderImageDTO{
			UserId: "123",
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
			File: file,
		}

		mockOrderProvider.EXPECT().UpdateStatus("test.blend", int64(123), int64(0)).Return(errors.New("error"))

		err = orderService.UpdateOrderImage(dto)
		assert.Error(t, err)
	})

	t.Run("UpdateDownloadLinkError", func(t *testing.T) {
		file, err := ioutil.TempFile("", "testfile")
		assert.NoError(t, err)
		defer os.Remove("123/test.jpg")

		dto := dto.UpdateOrderImageDTO{
			UserId: "123",
			Header: &multipart.FileHeader{
				Filename: "test.jpg",
			},
			File: file,
		}

		mockOrderProvider.EXPECT().UpdateStatus("test.blend", int64(123), int64(0)).Return(nil)
		mockOrderProvider.EXPECT().UpdateDownloadLink(int64(123), "test.blend", "http://localhost:8080/123/image/download/test.jpg").Return(errors.New("error"))

		err = orderService.UpdateOrderImage(dto)
		assert.Error(t, err)
	})
}

func TestOrderService_SoftDeleteOneOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)

	orderService := service.NewOrderService(mockOrderProvider, nil, nil, "", "", nil, nil)

	t.Run("Success", func(t *testing.T) {
		mockOrderProvider.EXPECT().SoftDelete(int64(123)).Return(nil)

		err := orderService.SoftDeleteOneOrder(123)
		assert.NoError(t, err)
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockOrderProvider.EXPECT().SoftDelete(int64(123)).Return(repository.ErrNoOrdersFound)

		err := orderService.SoftDeleteOneOrder(123)
		assert.ErrorIs(t, err, service.ErrOrderNotFound)
	})

	t.Run("OrderProviderError", func(t *testing.T) {
		mockOrderProvider.EXPECT().SoftDelete(int64(123)).Return(errors.New("error"))

		err := orderService.SoftDeleteOneOrder(123)
		assert.Error(t, err)
	})
}

func TestOrderService_GetManyOrders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)

	m := make(map[int64]string, 2)
	m[1] = "status1"
	m[2] = "status2"

	orderService := service.NewOrderService(mockOrderProvider, nil, m, "", "", nil, nil)

	t.Run("Success", func(t *testing.T) {
		expectedOrders := []models.Order{
			{Id: 1, FileName: "file1", CreationDate: time.Now(), StatusId: 1, DownloadLink: sql.NullString{Valid: true, String: "link1"}},
			{Id: 2, FileName: "file2", CreationDate: time.Now(), StatusId: 2, DownloadLink: sql.NullString{Valid: true, String: "link2"}},
		}
		mockOrderProvider.EXPECT().GetMany(int64(123)).Return(expectedOrders, nil)

		ordersDTO, err := orderService.GetManyOrders(123)
		assert.NoError(t, err)
		assert.Len(t, ordersDTO, 2)
		assert.Equal(t, "file1", ordersDTO[0].Filename)
		assert.Equal(t, "link1", ordersDTO[0].DownloadLink)
		assert.Equal(t, "status1", ordersDTO[0].OrderStatus)
		assert.Equal(t, "file2", ordersDTO[1].Filename)
		assert.Equal(t, "link2", ordersDTO[1].DownloadLink)
		assert.Equal(t, "status2", ordersDTO[1].OrderStatus)
	})

	t.Run("NoOrdersFound", func(t *testing.T) {
		mockOrderProvider.EXPECT().GetMany(int64(123)).Return([]models.Order{}, repository.ErrNoOrdersFound)

		ordersDTO, err := orderService.GetManyOrders(123)
		assert.NoError(t, err)
		assert.Empty(t, ordersDTO)
	})

	t.Run("OrderProviderError", func(t *testing.T) {
		mockOrderProvider.EXPECT().GetMany(int64(123)).Return(nil, errors.New("error"))

		ordersDTO, err := orderService.GetManyOrders(123)
		assert.Error(t, err)
		assert.Nil(t, ordersDTO)
	})
}

func TestOrderService_GetOneOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderProvider := mocks.NewMockOrderProvider(mockCtrl)

	m := make(map[int64]string, 1)
	m[1] = "status1"

	orderService := service.NewOrderService(mockOrderProvider, nil, m, "", "", nil, nil)

	t.Run("Success", func(t *testing.T) {
		expectedOrder := models.Order{
			Id:           1,
			FileName:     "file1",
			CreationDate: time.Now(),
			StatusId:     1,
			DownloadLink: sql.NullString{Valid: true, String: "link1"},
		}
		mockOrderProvider.EXPECT().GetOne(int64(123)).Return(&expectedOrder, nil)

		orderDTO, err := orderService.GetOneOrder(123)
		assert.NoError(t, err)
		assert.NotNil(t, orderDTO)
		assert.Equal(t, int64(1), orderDTO.Id)
		assert.Equal(t, "file1", orderDTO.Filename)
		assert.Equal(t, "link1", orderDTO.DownloadLink)
		assert.Equal(t, "status1", orderDTO.OrderStatus)
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockOrderProvider.EXPECT().GetOne(int64(123)).Return(nil, repository.ErrNoOrdersFound)

		orderDTO, err := orderService.GetOneOrder(123)
		assert.ErrorIs(t, err, service.ErrOrderNotFound)
		assert.Nil(t, orderDTO)
	})

	t.Run("OrderProviderError", func(t *testing.T) {
		mockOrderProvider.EXPECT().GetOne(int64(123)).Return(nil, errors.New("error"))

		orderDTO, err := orderService.GetOneOrder(123)
		assert.Error(t, err)
		assert.Nil(t, orderDTO)
	})
}
