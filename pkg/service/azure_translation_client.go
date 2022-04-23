//go:generate mockery --output mock --name AzureTranslationClient
package service

import (
	"context"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
)

type AzureTranslationClient interface {
	DictionaryLookup(ctx context.Context, text string, fromLang, toLang domain.Lang2) ([]AzureTranslation, error)
}
