package presenter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/handler/entity"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	libD "github.com/kujilabo/cocotola-translator-api/src/lib/domain"
)

type userPresenter struct {
	c *gin.Context
}

func NewUserPresenter(c *gin.Context) usecase.UserPresenter {
	return &userPresenter{
		c: c,
	}
}

func (p *userPresenter) WriteTranslations(ctx context.Context, translations []domain.Translation) error {
	response, err := p.toTranslationFindResposne(ctx, translations)
	if err != nil {
		return err
	}

	p.c.JSON(http.StatusOK, response)
	return nil
}

func (p *userPresenter) WriteTranslation(ctx context.Context, translation domain.Translation) error {
	response, err := p.toTranslationResposne(ctx, translation)
	if err != nil {
		return err
	}

	p.c.JSON(http.StatusOK, response)
	return nil
}

func (p *userPresenter) toTranslationFindResposne(ctx context.Context, translations []domain.Translation) (*entity.TranslationFindResponseHTTPEntity, error) {

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

func (p *userPresenter) toTranslationResposne(context context.Context, translation domain.Translation) (*entity.TranslationHTTPEntity, error) {
	e := &entity.TranslationHTTPEntity{
		Lang2:      translation.GetLang2().String(),
		Text:       translation.GetText(),
		Pos:        int(translation.GetPos()),
		Translated: translation.GetTranslated(),
		Provider:   translation.GetProvider(),
	}
	return e, libD.Validator.Struct(e)
}
