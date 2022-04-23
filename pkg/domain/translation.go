//go:generate mockery --output mock --name Translation
package domain

import (
	"time"

	lib "github.com/kujilabo/cocotola-translator-api/pkg_lib/domain"
)

type Translation interface {
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetText() string
	GetPos() WordPos
	GetLang() Lang2
	GetTranslated() string
	GetProvider() string
}

type translation struct {
	Version    int `validate:"required,gte=1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Text       string `validate:"required"`
	Pos        WordPos
	Lang2      Lang2
	Translated string
	Provider   string
}

func NewTranslation(version int, createdAt time.Time, updatedAt time.Time, text string, pos WordPos, lang Lang2, translated, provider string) (Translation, error) {
	m := &translation{
		// ID:         id,
		Version:    version,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Text:       text,
		Pos:        pos,
		Lang2:      lang,
		Translated: translated,
		Provider:   provider,
	}

	return m, lib.Validator.Struct(m)
}

// func (t *translation) GetID() TranslationID {
// 	return t.ID
// }

func (t *translation) GetVersion() int {
	return t.Version
}

func (t *translation) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *translation) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

func (t *translation) GetText() string {
	return t.Text
}

func (t *translation) GetPos() WordPos {
	return t.Pos
}

func (t *translation) GetLang() Lang2 {
	return t.Lang2
}

func (t *translation) GetTranslated() string {
	return t.Translated
}

func (t *translation) GetProvider() string {
	return t.Provider
}
