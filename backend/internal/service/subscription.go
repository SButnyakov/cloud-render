package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/external"
	"time"
)

type SubscriptionService struct {
	subscriptionProvider SubscriptionProvider
	config               *config.Config
}

type SubscriptionProvider interface {
	GetExpireDate(uid int64) (*time.Time, error)
}

func NewSubscriptionService(subscriptionProvider SubscriptionProvider, config *config.Config) *SubscriptionService {
	return &SubscriptionService{
		subscriptionProvider: subscriptionProvider,
		config:               config,
	}
}

func (s *SubscriptionService) GetExpireDateWithUserInfo(id int64) (*dto.UserInfoDTO, error) {
	userDataChan := make(chan *dto.GetUserDTO)
	expDateChan := make(chan *time.Time)

	go s.asyncGetUserInfo(userDataChan, s.config.External.SSOUserInfo, id)
	go s.asyncGetExpireDate(expDateChan, id)

	userDataDTO, ok := <-userDataChan
	if !ok {
		return nil, ErrExternalError
	}

	expDate, ok := <-expDateChan
	if !ok {
		return nil, ErrFailedToGetSubscription
	}

	return &dto.UserInfoDTO{
		Login:          userDataDTO.Login,
		Email:          userDataDTO.Email,
		ExpirationDate: expDate,
	}, nil
}

func (s *SubscriptionService) asyncGetExpireDate(out chan<- *time.Time, id int64) {
	time, err := s.subscriptionProvider.GetExpireDate(id)
	if err != nil {
		close(out)
		return
	}
	out <- time
}

func (s *SubscriptionService) asyncGetUserInfo(out chan<- *dto.GetUserDTO, url string, id int64) {
	userDTO, err := external.UserInfo(url, id)
	if err != nil {
		close(out)
		return
	}
	out <- userDTO
}
