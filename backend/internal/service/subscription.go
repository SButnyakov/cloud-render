package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/external"
	"cloud-render/internal/models"
	"fmt"
	"time"
)

type SubscriptionService struct {
	subscriptionProvider SubscriptionProvider
	config               *config.Config
	subMap               SubscriptionsMap
	payMap               PaymentsMap
}

type SubscriptionProvider interface {
	GetExpireDate(uid int64) (*time.Time, error)
	Create(subscription models.Subscription, payment models.Payment) error
	Update(subscription models.Subscription, payment models.Payment) error
}

type SubscriptionsMap map[string]int64
type PaymentsMap map[string]int64

func NewSubscriptionService(subscriptionProvider SubscriptionProvider, config *config.Config,
	subMap SubscriptionsMap, payMap PaymentsMap) *SubscriptionService {
	return &SubscriptionService{
		subscriptionProvider: subscriptionProvider,
		config:               config,
		subMap:               subMap,
		payMap:               payMap,
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
		ExpirationDate: *expDate,
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

func (s *SubscriptionService) SubscribeUser(id int64) error {
	pTypeId, ok := s.payMap[s.config.Payments.SubPremiumMonth]
	if !ok {
		return ErrPaymentTypeNotFound
	}

	sTypeId, ok := s.subMap[s.config.Subscriptions.Premium]
	if !ok {
		return ErrSubscriptionTypeNotFound
	}

	fmt.Println("subscribing")

	expireDate, err := s.subscriptionProvider.GetExpireDate(id)

	if err != nil {
		return err
	}

	if expireDate.IsZero() {
		return s.createSubscription(id, pTypeId, sTypeId)
	} else {
		return s.updateSubscription(id, pTypeId, sTypeId, expireDate)
	}
}

func (s *SubscriptionService) createSubscription(id, pType, sType int64) error {
	return s.subscriptionProvider.Create(models.Subscription{
		UserId:     id,
		TypeId:     sType,
		ExpireDate: time.Now().AddDate(0, 1, 0),
	}, models.Payment{
		UserID:   id,
		TypeId:   pType,
		DateTime: time.Now(),
	})
}

func (s *SubscriptionService) updateSubscription(id, pType, sType int64, expDate *time.Time) error {
	return s.subscriptionProvider.Update(models.Subscription{
		UserId:     id,
		TypeId:     sType,
		ExpireDate: expDate.AddDate(0, 1, 0),
	}, models.Payment{
		UserID:   id,
		TypeId:   pType,
		DateTime: time.Now(),
	})
}
