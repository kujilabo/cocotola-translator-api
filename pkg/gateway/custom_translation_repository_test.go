package gateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/gateway"
)

func Test_customTranslationRepository_FindByText(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
	}
}

func Test_customTranslationRepository_FindByFirstLetter(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()
	for _, db := range dbList() {
		type args struct {
			firstLetter string
			lang2       domain.Lang2
		}
		result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from custom_translation")
		assert.NoError(t, result.Error)

		book, err := domain.NewTranslation(1, time.Now(), time.Now(), "book", domain.PosNoun, domain.Lang2JA, "本", "custom")
		assert.NoError(t, err)

		result = db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(fmt.Sprintf("insert into custom_translation (version,text,pos,lang2,translated) values(%d,'%s',%d,'%s','%s')", uint(book.GetVersion()), book.GetText(), int(book.GetPos()), book.GetLang2().String(), book.GetTranslated()))
		assert.NoError(t, result.Error)

		tests := []struct {
			name    string
			args    args
			want    []domain.Translation
			wantErr bool
		}{
			{
				name: "found a record",
				args: args{
					firstLetter: "b",
					lang2:       domain.Lang2JA,
				},
				want: []domain.Translation{
					book,
				},
				wantErr: false,
			},
		}
		r := gateway.NewCustomTranslationRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := r.FindByFirstLetter(bg, tt.args.lang2, tt.args.firstLetter)
				if (err != nil) != tt.wantErr {
					t.Errorf("customTranslationRepository.FindByFirstLetter() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil {
					assert.Equal(t, len(got), len(tt.want))
					assert.Equal(t, got[0].GetTranslated(), tt.want[0].GetTranslated())
				}

			})
		}
	}
}
