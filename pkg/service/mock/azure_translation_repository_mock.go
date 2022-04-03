package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
)

type AzureTranslationRepositoryMock struct {
	mock.Mock
}

func (m *AzureTranslationRepositoryMock) Add(ctx context.Context, lang domain.Lang2, text string, result []service.AzureTranslation) error {
	args := m.Called(ctx, lang, text, result)
	return args.Error(0)
}

func (m *AzureTranslationRepositoryMock) Find(ctx context.Context, lang domain.Lang2, text string) ([]service.AzureTranslation, error) {
	args := m.Called(ctx, lang, text)
	return args.Get(0).([]service.AzureTranslation), args.Error(1)
}

func (m *AzureTranslationRepositoryMock) FindByText(ctx context.Context, lang domain.Lang2, text string) ([]domain.Translation, error) {
	args := m.Called(ctx, lang, text)
	return args.Get(0).([]domain.Translation), args.Error(1)
}

func (m *AzureTranslationRepositoryMock) FindByTextAndPos(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	args := m.Called(ctx, lang, text, pos)
	return args.Get(0).(domain.Translation), args.Error(1)
}

func (m *AzureTranslationRepositoryMock) FindByFirstLetter(ctx context.Context, lang domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	args := m.Called(ctx, lang, firstLetter)
	return args.Get(0).([]domain.Translation), args.Error(1)
}

func (m *AzureTranslationRepositoryMock) Contain(ctx context.Context, lang domain.Lang2, text string) (bool, error) {
	args := m.Called(ctx, lang, text)
	return args.Bool(0), args.Error(1)
}
