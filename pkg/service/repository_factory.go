//go:generate mockery --output mock --name RepositoryFactory
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewAzureTranslationRepository(ctx context.Context) AzureTranslationRepository

	NewCustomTranslationRepository(ctx context.Context) CustomTranslationRepository
}
