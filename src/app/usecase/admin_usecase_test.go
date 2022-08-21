//go:generate mockery --output mock --name AdminUsecase
package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	service_mock "github.com/kujilabo/cocotola-translator-api/src/app/service/mock"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	"github.com/stretchr/testify/assert"
)

func matchErrorFunc(expected error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, actual error, args ...interface{}) bool {
		if errors.Is(actual, expected) {
			return true
		}
		return assert.Fail(t, fmt.Sprintf("error type is mismatch. expected: %v, actual: %v", expected, actual))
	}
}

func Test_adminUsecase_RemoveTranslation(t *testing.T) {
	bg := context.Background()

	// given
	customRepo := new(service_mock.CustomTranslationRepository)
	customRepo.On("Remove", anythingOfContext, domain.Lang2JA, "apple", domain.PosNoun).Return(nil)
	customRepo.On("Remove", anythingOfContext, domain.Lang2JA, "orange", domain.PosNoun).Return(service.ErrTranslationNotFound)
	rf := new(service_mock.RepositoryFactory)
	rf.On("NewCustomTranslationRepository", anythingOfContext).Return(customRepo)
	adminUsecase := usecase.NewAdminUsecase(rf)

	type args struct {
		lang2 domain.Lang2
		text  string
		pos   domain.WordPos
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{"word is registered", args{domain.Lang2JA, "apple", domain.PosNoun}, assert.NoError},
		{"word is not registered", args{domain.Lang2JA, "orange", domain.PosNoun}, matchErrorFunc(service.ErrTranslationNotFound)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			err := adminUsecase.RemoveTranslation(bg, tt.args.lang2, tt.args.text, tt.args.pos)

			// then
			tt.assertion(t, err)
		})
	}
}
