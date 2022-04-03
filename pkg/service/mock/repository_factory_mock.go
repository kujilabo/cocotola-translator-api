package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-translator-api/pkg/service"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewAzureTranslationRepository(ctx context.Context) (service.AzureTranslationRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.AzureTranslationRepository), args.Error(1)
}

func (m *RepositoryFactoryMock) NewCustomTranslationRepository(ctx context.Context) (service.CustomTranslationRepository, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.CustomTranslationRepository), args.Error(1)
}
