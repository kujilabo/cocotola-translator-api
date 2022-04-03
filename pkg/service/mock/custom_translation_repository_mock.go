package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
)

type CustomTranslationRepositoryMock struct {
	mock.Mock
}

func (m *CustomTranslationRepositoryMock) Add(ctx context.Context, param service.TranslationAddParameter) error {
	args := m.Called(ctx, param)
	return args.Error(0)
}

func (m *CustomTranslationRepositoryMock) Update(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	args := m.Called(ctx, lang, text, pos, param)
	return args.Error(0)
}

func (m *CustomTranslationRepositoryMock) Remove(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) error {
	args := m.Called(ctx, lang, text, pos)
	return args.Error(0)
}

func (m *CustomTranslationRepositoryMock) FindByText(ctx context.Context, lang domain.Lang2, text string) ([]domain.Translation, error) {
	args := m.Called(ctx, lang, text)
	return args.Get(0).([]domain.Translation), args.Error(1)
}

func (m *CustomTranslationRepositoryMock) FindByTextAndPos(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	args := m.Called(ctx, lang, text, pos)
	return args.Get(0).(domain.Translation), args.Error(1)
}

func (m *CustomTranslationRepositoryMock) FindByFirstLetter(ctx context.Context, lang domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	args := m.Called(ctx, lang, firstLetter)
	return args.Get(0).([]domain.Translation), args.Error(1)
}

func (m *CustomTranslationRepositoryMock) Contain(ctx context.Context, lang domain.Lang2, text string) (bool, error) {
	args := m.Called(ctx, lang, text)
	return args.Bool(0), args.Error(1)
}
