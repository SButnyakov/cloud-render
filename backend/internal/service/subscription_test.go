package service_test

import (
	"cloud-render/internal/lib/config"
	mocks "cloud-render/internal/mocks/repository"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"cloud-render/internal/service"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	SubPremiumMonth = "sub-premium-month"
	Premium         = "premium"
)

var (
	SubMap = service.SubscriptionsMap{}
	PayMap = service.PaymentsMap{}
	Cfg    = config.Config{}
	Log    = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
)

func TestSubscriptionService_SubscribeUser(t *testing.T) {
	defer clearMaps()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSubProvider := mocks.NewMockSubscriptionProvider(mockCtrl)

	subService := service.NewSubscriptionService(mockSubProvider, &Cfg, SubMap, PayMap, Log)

	initConfig()

	err := subService.SubscribeUser(1)
	assert.Error(t, err)
	assert.Equal(t, service.ErrPaymentTypeNotFound, err)

	PayMap[Cfg.Payments.SubPremiumMonth] = 1

	err = subService.SubscribeUser(1)
	assert.Error(t, err)
	assert.Equal(t, service.ErrSubscriptionTypeNotFound, err)

	SubMap[Cfg.Subscriptions.Premium] = 1

	userID := int64(1)

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(nil, repository.ErrSubscriptionNotFound).
		Times(1)
	mockSubProvider.EXPECT().
		Create(gomock.AssignableToTypeOf(models.Subscription{}), gomock.AssignableToTypeOf(models.Payment{})).
		Return(nil).
		Times(1)

	err = subService.SubscribeUser(1)
	assert.NoError(t, err)

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(nil, errors.New("unknown")).
		Times(1)

	err = subService.SubscribeUser(1)
	require.Error(t, err)
	assert.EqualError(t, err, "unknown")

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(&time.Time{}, nil).
		Times(1)
	mockSubProvider.EXPECT().
		Create(gomock.AssignableToTypeOf(models.Subscription{}), gomock.AssignableToTypeOf(models.Payment{})).
		Return(nil).
		Times(1)

	err = subService.SubscribeUser(1)
	assert.NoError(t, err)

	dayBefore := time.Now().AddDate(0, 0, -1)

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(&dayBefore, nil).
		Times(1)
	mockSubProvider.EXPECT().
		Create(gomock.AssignableToTypeOf(models.Subscription{}), gomock.AssignableToTypeOf(models.Payment{})).
		Return(nil).
		Times(1)

	err = subService.SubscribeUser(1)
	assert.NoError(t, err)

	dayAfter := time.Now().AddDate(0, 0, 1)

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(&dayAfter, nil).
		Times(1)
	mockSubProvider.EXPECT().
		Update(gomock.AssignableToTypeOf(models.Subscription{}), gomock.AssignableToTypeOf(models.Payment{})).
		Return(nil).
		Times(1)

	err = subService.SubscribeUser(1)
	assert.NoError(t, err)
}

func TestSubscriptionService_GetExpireDateWithUserInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSubProvider := mocks.NewMockSubscriptionProvider(mockCtrl)

	subService := service.NewSubscriptionService(mockSubProvider, &Cfg, SubMap, PayMap, Log)

	userID := int64(1)

	monthAfter := time.Now().AddDate(0, 1, 0)

	initConfig()

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(&monthAfter, nil).
		Times(1)

	_, err := subService.GetExpireDateWithUserInfo(userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrExternalError)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(nil, errors.New("any")).
		Times(1)
	httpmock.RegisterResponder("GET", "http://localhost:8081/info/1",
		httpmock.NewStringResponder(200, `{"login": "login", "email": "email@gmail.com"}`))

	_, err = subService.GetExpireDateWithUserInfo(userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrFailedToGetSubscription)

	mockSubProvider.EXPECT().
		GetExpireDate(userID).
		Return(&monthAfter, nil).
		Times(1)

	dto, err := subService.GetExpireDateWithUserInfo(userID)
	assert.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, "login", dto.Login)
	assert.Equal(t, "email@gmail.com", dto.Email)
	assert.Equal(t, monthAfter.Day(), dto.ExpirationDate.Day())
	assert.Equal(t, monthAfter.Month(), dto.ExpirationDate.Month())
	assert.Equal(t, monthAfter.Hour(), dto.ExpirationDate.Hour())
}

func initConfig() {
	Cfg.Payments.SubPremiumMonth = SubPremiumMonth
	Cfg.Subscriptions.Premium = Premium
	Cfg.External.SSOUserInfo = "http://localhost:8081/info/%d"
}

func clearMaps() {
	SubMap = service.SubscriptionsMap{}
	PayMap = service.PaymentsMap{}
}
