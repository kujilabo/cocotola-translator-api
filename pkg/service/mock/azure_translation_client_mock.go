package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
)

type AzureTranslationClientMock struct {
	mock.Mock
}

func (m *AzureTranslationClientMock) DictionaryLookup(ctx context.Context, text string, fromLang, toLang domain.Lang2) ([]service.AzureTranslation, error) {
	args := m.Called(ctx, text, fromLang, toLang)
	return args.Get(0).([]service.AzureTranslation), args.Error(1)
}
