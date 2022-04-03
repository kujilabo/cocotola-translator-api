package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	lib "github.com/kujilabo/cocotola-translator-api/pkg_lib/domain"
)

// var ErrCustomTranslationNotFound = errors.New("azure translation not found")
// var ErrCustomTranslationAlreadyExists = errors.New("azure translation already exists")

var ErrTranslationAlreadyExists = errors.New("custsomtranslation already exists")

type TranslationAddParameter interface {
	GetText() string
	GetPos() domain.WordPos
	GetLang() domain.Lang2
	GetTranslated() string
}

type translationAddParameter struct {
	Text       string `validate:"required"`
	Pos        domain.WordPos
	Lang2      domain.Lang2
	Translated string
}

func NewTransalationAddParameter(text string, pos domain.WordPos, lang domain.Lang2, translated string) (TranslationAddParameter, error) {
	m := &translationAddParameter{
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
	}

	return m, lib.Validator.Struct(m)
}

func (p *translationAddParameter) GetText() string {
	return p.Text
}

func (p *translationAddParameter) GetPos() domain.WordPos {
	return p.Pos
}

func (p *translationAddParameter) GetLang() domain.Lang2 {
	return p.Lang2
}

func (p *translationAddParameter) GetTranslated() string {
	return p.Translated
}

type TranslationUpdateParameter interface {
	GetTranslated() string
}

type translationUpdateParameter struct {
	Translated string `validate:"required"`
}

func NewTransaltionUpdateParameter(translated string) (TranslationUpdateParameter, error) {
	m := &translationUpdateParameter{
		Translated: translated,
	}

	return m, lib.Validator.Struct(m)
}

func (p *translationUpdateParameter) GetTranslated() string {
	return p.Translated
}

type CustomTranslationRepository interface {
	Add(ctx context.Context, param TranslationAddParameter) (domain.TranslationID, error)

	Update(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos, param TranslationUpdateParameter) error

	Remove(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) error

	FindByText(ctx context.Context, lang domain.Lang2, text string) ([]domain.Translation, error)

	FindByTextAndPos(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)

	FindByFirstLetter(ctx context.Context, lang domain.Lang2, firstLetter string) ([]domain.Translation, error)

	Contain(ctx context.Context, lang domain.Lang2, text string) (bool, error)
}
