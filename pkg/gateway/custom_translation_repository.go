package gateway

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	libD "github.com/kujilabo/cocotola-translator-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-translator-api/pkg_lib/gateway"
)

type customTranslationRepository struct {
	db *gorm.DB
}

type customTranslationEntity struct {
	Version    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Text       string
	Pos        int
	Lang2      string
	Translated string
}

func (e *customTranslationEntity) TableName() string {
	return "custom_translation"
}

func (e *customTranslationEntity) toModel() (domain.Translation, error) {
	lang2, err := domain.NewLang2(e.Lang2)
	if err != nil {
		return nil, err
	}

	t, err := domain.NewTranslation(e.Version, e.CreatedAt, e.UpdatedAt, e.Text, domain.WordPos(e.Pos), lang2, e.Translated, "custom")
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewCustomTranslationRepository(db *gorm.DB) service.CustomTranslationRepository {
	return &customTranslationRepository{
		db: db,
	}
}

func (r *customTranslationRepository) Add(ctx context.Context, param service.TranslationAddParameter) error {
	_, span := tracer.Start(ctx, "customTranslationRepository.Add")
	defer span.End()

	entity := customTranslationEntity{
		Version:    1,
		Text:       param.GetText(),
		Lang2:      param.GetLang2().String(),
		Pos:        int(param.GetPos()),
		Translated: param.GetTranslated(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		err := libG.ConvertDuplicatedError(result.Error, service.ErrTranslationAlreadyExists)
		return xerrors.Errorf("failed to Add translation. err: %w", err)
	}

	return nil
}

func (r *customTranslationRepository) Update(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	_, span := tracer.Start(ctx, "customTranslationRepository.Update")
	defer span.End()

	result := r.db.Model(&customTranslationEntity{}).
		Where("lang2 = ? and text = ? and pos = ?",
			lang2.String(), text, int(pos)).
		Updates(map[string]interface{}{
			"translated": param.GetTranslated(),
		})
	if result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, service.ErrTranslationAlreadyExists)
	}

	if result.RowsAffected != 1 {
		return errors.New("Error")
	}

	return nil
}

func (r *customTranslationRepository) Remove(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) error {
	_, span := tracer.Start(ctx, "customTranslationRepository.Remove")
	defer span.End()

	result := r.db.
		Where("lang2 = ? and text = ? and pos = ?",
			lang2.String(), text, int(pos)).
		Delete(&customTranslationEntity{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *customTranslationRepository) FindByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error) {
	_, span := tracer.Start(ctx, "customTranslationRepository.FindByText")
	defer span.End()

	entities := []customTranslationEntity{}
	if result := r.db.Where(&customTranslationEntity{
		Text:  text,
		Lang2: lang2.String(),
	}).Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.Translation, len(entities))
	for i, e := range entities {
		t, err := e.toModel()
		if err != nil {
			return nil, err
		}
		results[i] = t
	}

	return results, nil
}

func (r *customTranslationRepository) FindByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	_, span := tracer.Start(ctx, "customTranslationRepository.FindByTextAndPos")
	defer span.End()

	entity := customTranslationEntity{}
	if result := r.db.Where(&customTranslationEntity{
		Text:  text,
		Lang2: lang2.String(),
		Pos:   int(pos),
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrTranslationNotFound
		}

		return nil, result.Error
	}

	return entity.toModel()
}

func (r *customTranslationRepository) FindByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	_, span := tracer.Start(ctx, "customTranslationRepository.FindByFirstLetter")
	defer span.End()

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

	entities := []customTranslationEntity{}
	if result := r.db.Where(&customTranslationEntity{
		Lang2: lang2.String(),
	}).Where("text like ? OR text like ?", upper, lower).Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.Translation, len(entities))
	for i, e := range entities {
		t, err := e.toModel()
		if err != nil {
			return nil, err
		}
		results[i] = t
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

func (r *customTranslationRepository) Contain(ctx context.Context, lang2 domain.Lang2, text string) (bool, error) {
	_, span := tracer.Start(ctx, "customTranslationRepository.Contain")
	defer span.End()

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
