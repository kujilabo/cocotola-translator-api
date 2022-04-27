//go:generate mockery --output mock --name AzureTranslationRepository
package service

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
)

var ErrAzureTranslationAlreadyExists = errors.New("azure translation already exists")

type AzureTranslation struct {
	Pos        domain.WordPos
	Target     string
	Confidence float64
}

func (t *AzureTranslation) ToTranslation(lang2 domain.Lang2, text string) (domain.Translation, error) {
	return domain.NewTranslation(1, time.Now(), time.Now(), text, t.Pos, lang2, t.Target, "azure")
}

type TranslationSearchCondition struct {
	PageNo   int
	PageSize int
}

type TranslationSearchResult struct {
	TotalCount int64
	Results    [][]AzureTranslation
}

type AzureTranslationRepository interface {
	Add(ctx context.Context, lang2 domain.Lang2, text string, result []AzureTranslation) error

	Find(ctx context.Context, lang2 domain.Lang2, text string) ([]AzureTranslation, error)

	FindByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)

	FindByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error)

	FindByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error)

	Contain(ctx context.Context, lang2 domain.Lang2, text string) (bool, error)
}
