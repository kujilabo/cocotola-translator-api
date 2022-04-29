package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	libD "github.com/kujilabo/cocotola-translator-api/src/lib/domain"
	libG "github.com/kujilabo/cocotola-translator-api/src/lib/gateway"
)

type azureTranslationRepository struct {
	db *gorm.DB
}

type azureTranslationEntity struct {
	Text   string
	Lang2  string
	Result string
}

func (e *azureTranslationEntity) TableName() string {
	return "azure_translation"
}

func NewAzureTranslationRepository(db *gorm.DB) service.AzureTranslationRepository {
	return &azureTranslationRepository{
		db: db,
	}
}

func (r *azureTranslationRepository) Add(ctx context.Context, lang2 domain.Lang2, text string, result []service.AzureTranslation) error {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	entity := azureTranslationEntity{
		Text:   text,
		Lang2:  lang2.String(),
		Result: string(resultBytes),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, service.ErrAzureTranslationAlreadyExists)
	}

	return nil
}

func (r *azureTranslationRepository) Find(ctx context.Context, lang2 domain.Lang2, text string) ([]service.AzureTranslation, error) {
	entity := azureTranslationEntity{}

	if result := r.db.Where(&azureTranslationEntity{
		Text:  text,
		Lang2: lang2.String(),
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrTranslationNotFound
		}

		return nil, result.Error
	}

	result := make([]service.AzureTranslation, 0)
	if err := json.Unmarshal([]byte(entity.Result), &result); err != nil {
		return nil, err
	}

	return result, nil
}
func (r *azureTranslationRepository) FindByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	results, err := r.Find(ctx, lang2, text)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Pos != pos {
			continue
		}

		t, err := r.ToTranslation(lang2, text)
		if err != nil {
			return nil, err
		}

		return t, nil
	}
	return nil, service.ErrTranslationNotFound
}

func (r *azureTranslationRepository) FindByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error) {
	azureTranslations, err := r.Find(ctx, lang2, text)
	if err != nil {
		return nil, err
	}

	translations := make([]domain.Translation, 0)
	for _, r := range azureTranslations {
		t, err := r.ToTranslation(lang2, text)
		if err != nil {
			return nil, err
		}
		translations = append(translations, t)
	}

	return translations, nil
}

func (r *azureTranslationRepository) FindByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	if len(firstLetter) != 1 {
		return nil, libD.ErrInvalidArgument
	}

	matched, err := regexp.Match("^[a-zA-Z]$", []byte(firstLetter))
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, libD.ErrInvalidArgument
	}
	upper := strings.ToUpper(firstLetter) + "%"
	lower := strings.ToLower(firstLetter) + "%"

	entities := []azureTranslationEntity{}
	if result := r.db.Where(&azureTranslationEntity{
		Lang2: lang2.String(),
	}).Where("text like ? OR text like ?", upper, lower).Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.Translation, 0)
	for _, e := range entities {
		result := make([]service.AzureTranslation, 0)
		if err := json.Unmarshal([]byte(e.Result), &result); err != nil {
			return nil, err
		}
		for _, r := range result {
			t, err := r.ToTranslation(lang2, e.Text)
			if err != nil {
				return nil, err
			}
			results = append(results, t)
		}
	}

	return results, nil
}

// func (r *azureTranslationRepository) FindTranslations(ctx context.Context, param *domain.AzureTranslationSearchCondition) (*domain.AzureTranslation, error) {
// 	limit := param.PageSize
// 	offset := (param.PageNo - 1) * param.PageSize
// 	var entities []azureTranslationEntity
// 	if result := r.db.Limit(limit).Offset(offset).Find(&entities); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	var count int64
// 	if result := r.db.Model(azureTranslationEntity{}).Count(&count); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	results := make([][]domain.AzureTranslation, len(entities))
// 	for i, e := range entities {
// 		result := make([]domain.AzureTranslation, 0)
// 		if err := json.Unmarshal([]byte(e.Result), &result); err != nil {
// 			return nil, err
// 		}
// 		results[i] = result
// 	}

// 	return &domain.AzureTranslationSearchResult{
// 		TotalCount: count,
// 		Results:    results,
// 	}, nil
// }

func (r *azureTranslationRepository) Contain(ctx context.Context, lang2 domain.Lang2, text string) (bool, error) {
	entity := azureTranslationEntity{}

	if result := r.db.Where(&azureTranslationEntity{
		Text:  text,
		Lang2: lang2.String(),
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}
