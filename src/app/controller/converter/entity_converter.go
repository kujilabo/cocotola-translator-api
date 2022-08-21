package converter

import (
	"context"

	"github.com/kujilabo/cocotola-translator-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	libD "github.com/kujilabo/cocotola-translator-api/src/lib/domain"
)

func ToTranslationFindResposne(ctx context.Context, translations []domain.Translation) (*entity.TranslationFindResponseHTTPEntity, error) {

	results := make([]entity.TranslationHTTPEntity, len(translations))
	for i, t := range translations {
		results[i] = entity.TranslationHTTPEntity{
			Lang2:      t.GetLang2().String(),
			Text:       t.GetText(),
			Pos:        int(t.GetPos()),
			Translated: t.GetTranslated(),
			Provider:   t.GetProvider(),
		}
	}

	e := &entity.TranslationFindResponseHTTPEntity{
		Results: results,
	}
	return e, libD.Validator.Struct(e)
}

func ToTranslationResposne(context context.Context, translation domain.Translation) (*entity.TranslationHTTPEntity, error) {
	e := &entity.TranslationHTTPEntity{
		Lang2:      translation.GetLang2().String(),
		Text:       translation.GetText(),
		Pos:        int(translation.GetPos()),
		Translated: translation.GetTranslated(),
		Provider:   translation.GetProvider(),
	}
	return e, libD.Validator.Struct(e)
}

func ToTranslationAddParameter(ctx context.Context, param *entity.TranslationAddParameterHTTPEntity) (service.TranslationAddParameter, error) {
	pos, err := domain.NewWordPos(param.Pos)
	if err != nil {
		return nil, err
	}

	lang2, err := domain.NewLang2(param.Lang2)
	if err != nil {
		return nil, err
	}
	return service.NewTransalationAddParameter(param.Text, pos, lang2, param.Translated)
}

func ToTranslationUpdateParameter(ctx context.Context, param *entity.TranslationUpdateParameterHTTPEntity) (service.TranslationUpdateParameter, error) {
	return service.NewTransaltionUpdateParameter(param.Translated)
}
