package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/converters"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	orderProvider    OrderProvider
	statusesStrToInt OrderStatusesMapStringToInt
	statusesIntToStr OrderStatusesMapIntToString
	inputPath        string
	cfg              *config.Config
	redis            *redis.Client
}

type OrderProvider interface {
	Create(order models.Order) error
	GetOne(id int64) (*models.Order, error)
	GetMany(id int64) ([]models.Order, error)
	UpdateStatus(storingName string, uid, statusId int64) error
	SoftDelete(id int64) error
}

type OrderStatusesMapStringToInt map[string]int64
type OrderStatusesMapIntToString map[int64]string

func NewOrderService(orderProvider OrderProvider, statusesStrToInt OrderStatusesMapStringToInt,
	statusesIntToStr OrderStatusesMapIntToString, inputPath string, cfg *config.Config,
	redis *redis.Client) *OrderService {
	return &OrderService{
		orderProvider:    orderProvider,
		statusesStrToInt: statusesStrToInt,
		statusesIntToStr: statusesIntToStr,
		inputPath:        inputPath,
		cfg:              cfg,
		redis:            redis,
	}
}

func (s *OrderService) CreateOrder(dto dto.CreateOrderDTO) error {
	userPath := filepath.Join(s.inputPath, strconv.FormatInt(dto.UserId, 10))
	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create user dir: %w", err)
	}

	storingName := strconv.FormatInt(time.Now().Unix(), 10) + "." + strings.Split(dto.Header.Filename, ".")[1]

	savePath := userPath + "/" + storingName

	f, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, dto.File)
	if err != nil {
		os.Remove(savePath)
		return fmt.Errorf("failed to write into file: %w", err)
	}

	err = s.orderProvider.Create(models.Order{
		FileName:     dto.Header.Filename,
		StoringName:  storingName,
		CreationDate: time.Now(),
		UserId:       dto.UserId,
		StatusId:     s.statusesStrToInt[s.cfg.OrderStatuses.InQueue],
	})
	if err != nil {
		os.Remove(savePath)
		return fmt.Errorf("failed to store new record: %w", err)
	}

	b, err := json.Marshal(models.RedisData{
		Format:     dto.Format,
		Resolution: dto.Resolution,
		SavePath:   savePath,
	})
	if err != nil {
		return fmt.Errorf("failed to create new redis record: %w", err)
	}

	s.redis.RPush(context.Background(), s.cfg.Redis.QueueName, string(b))

	return nil
}

func (s *OrderService) GetOneOrder(id int64) (*dto.GetOrderDTO, error) {
	order, err := s.orderProvider.GetOne(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoOrdersFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return &dto.GetOrderDTO{
		Id:           order.Id,
		Filename:     order.FileName,
		Date:         order.CreationDate,
		OrderStatus:  s.statusesIntToStr[order.StatusId],
		DownloadLink: converters.NullStringToString(order.DownloadLink),
	}, nil
}

func (s *OrderService) GetManyOrders(id int64) ([]dto.GetOrderDTO, error) {
	orders, err := s.orderProvider.GetMany(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoOrdersFound) {
			return []dto.GetOrderDTO{}, nil
		}
		return nil, err
	}

	ordersDTO := make([]dto.GetOrderDTO, len(orders))
	for i, v := range orders {
		ordersDTO[i] = dto.GetOrderDTO{
			Id:           v.Id,
			Filename:     v.FileName,
			Date:         v.CreationDate,
			OrderStatus:  s.statusesIntToStr[v.StatusId],
			DownloadLink: converters.NullStringToString(v.DownloadLink),
		}
	}

	return ordersDTO, nil
}

func (s *OrderService) UpdateOrderStatus(dto dto.UpdateOrderStatusDTO) error {
	err := s.orderProvider.UpdateStatus(dto.StoringName, dto.UserId, s.statusesStrToInt[dto.Status])
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) SoftDeleteOneOrder(id int64) error {
	err := s.orderProvider.SoftDelete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoOrdersFound) {
			return ErrOrderNotFound
		}
		return err
	}
	return nil
}
