package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	service_mock "github.com/kujilabo/cocotola-translator-api/src/app/service/mock"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
)

func test_userUsecase_DictionaryLookup_init(t *testing.T, ctx context.Context) (*service_mock.AzureTranslationClient, *service_mock.AzureTranslationRepository, *service_mock.CustomTranslationRepository, usecase.UserUsecase) {

	azureTranslationRepo := new(service_mock.AzureTranslationRepository)
	customTranslationRepo := new(service_mock.CustomTranslationRepository)
	rf := new(service_mock.RepositoryFactory)
	rf.On("NewAzureTranslationRepository", ctx).Return(azureTranslationRepo, nil)
	rf.On("NewCustomTranslationRepository", ctx).Return(customTranslationRepo, nil)
	azureTranslationClient := new(service_mock.AzureTranslationClient)
	userUsecase := usecase.NewUserUsecase(rf, azureTranslationClient)

	return azureTranslationClient, azureTranslationRepo, customTranslationRepo, userUsecase
}

func Test_userUsecase_DictionaryLookup_azureRepo(t *testing.T) {
	bg := context.Background()
	azureTranslationClient, azureTranslationRepo, customTranslationRepo, userUsecase := test_userUsecase_DictionaryLookup_init(t, bg)

	// given
	// - azureRepo
	azureRepoResults := []service.AzureTranslation{
		{
			Pos:        domain.PosNoun,
			Target:     "本ar",
			Confidence: 1,
		},
	}
	azureTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(true, nil)
	azureTranslationRepo.On("Find", bg, domain.Lang2JA, "book").Return(azureRepoResults, nil)
	// - azureClient
	azureClientResults := []service.AzureTranslation{}
	azureTranslationClient.On("DictionaryLookup", bg, "book", domain.Lang2EN, domain.Lang2JA).Return(azureClientResults, nil)
	azureTranslationRepo.On("Add", bg, domain.Lang2JA, "book", azureClientResults).Return(nil)
	// - customRepo has no data
	customTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(false, nil)

	// when
	actual, err := userUsecase.DictionaryLookup(bg, domain.Lang2EN, domain.Lang2JA, "book")
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(actual), 1)
	assert.Equal(t, actual[0].GetTranslated(), "本ar")
}

func Test_userUsecase_DictionaryLookup_azureClient(t *testing.T) {
	bg := context.Background()
	azureTranslationClient, azureTranslationRepo, customTranslationRepo, userUsecase := test_userUsecase_DictionaryLookup_init(t, bg)

	// given
	// - azureRepo has no data
	azureTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(false, nil)
	// - azureClient has data
	azureClientResults := []service.AzureTranslation{{
		Pos:        domain.PosNoun,
		Target:     "本ar",
		Confidence: 1,
	}}
	azureTranslationClient.On("DictionaryLookup", bg, "book", domain.Lang2EN, domain.Lang2JA).Return(azureClientResults, nil)
	azureTranslationRepo.On("Add", bg, domain.Lang2JA, "book", azureClientResults).Return(nil)
	// - customRepo has no data
	customTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(false, nil)

	// when
	actual, err := userUsecase.DictionaryLookup(bg, domain.Lang2EN, domain.Lang2JA, "book")
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(actual), 1)
	assert.Equal(t, actual[0].GetTranslated(), "本ar")
}

func Test_userUsecase_DictionaryLookup_azureRepo_azureClient(t *testing.T) {
	bg := context.Background()
	azureTranslationClient, azureTranslationRepo, customTranslationRepo, userUsecase := test_userUsecase_DictionaryLookup_init(t, bg)

	// given
	// - azureRepo has data
	azureRepoResults := []service.AzureTranslation{{
		Pos:        domain.PosNoun,
		Target:     "本ar",
		Confidence: 1,
	}}
	azureTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(true, nil)
	azureTranslationRepo.On("Find", bg, domain.Lang2JA, "book").Return(azureRepoResults, nil)
	azureClientResults := []service.AzureTranslation{{
		Pos:        domain.PosNoun,
		Target:     "本ac",
		Confidence: 1,
	}}
	// - azureClient has data
	azureTranslationClient.On("DictionaryLookup", bg, "book", domain.Lang2EN, domain.Lang2JA).Return(azureClientResults, nil)
	azureTranslationRepo.On("Add", bg, domain.Lang2JA, "book", azureClientResults).Return(nil)
	// - customRepo has no data
	customTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(false, nil)

	// when
	actual, err := userUsecase.DictionaryLookup(bg, domain.Lang2EN, domain.Lang2JA, "book")
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(actual), 1)
	assert.Equal(t, actual[0].GetTranslated(), "本ar")
}

func Test_userUsecase_DictionaryLookup_custom_azureRepo(t *testing.T) {
	bg := context.Background()
	_, azureTranslationRepo, customTranslationRepo, userUsecase := test_userUsecase_DictionaryLookup_init(t, bg)

	// given
	// - custom has data
	bookNoun, err := domain.NewTranslation(1, time.Now(), time.Now(), "book", domain.PosNoun, domain.Lang2JA, "本c", "")
	assert.NoError(t, err)
	customRepoResults := []domain.Translation{bookNoun}
	customTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(true, nil)
	customTranslationRepo.On("FindByText", bg, domain.Lang2JA, "book").Return(customRepoResults, nil)
	// - azureRepo has data
	azureRepoResults := []service.AzureTranslation{
		{
			Pos:        domain.PosNoun,
			Target:     "本ar",
			Confidence: 1,
		},
		{
			Pos:        domain.PosVerb,
			Target:     "予約するar",
			Confidence: 1,
		},
	}
	azureTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(true, nil)
	azureTranslationRepo.On("Find", bg, domain.Lang2JA, "book").Return(azureRepoResults, nil)
	// - azureClient has no data
	customTranslationRepo.On("Contain", bg, domain.Lang2JA, "book").Return(false, nil)
	// when
	actual, err := userUsecase.DictionaryLookup(bg, domain.Lang2EN, domain.Lang2JA, "book")
	assert.NoError(t, err)
	// then
	assert.Equal(t, len(actual), 2)
	assert.Equal(t, actual[0].GetTranslated(), "本c")
	assert.Equal(t, actual[1].GetTranslated(), "予約するar")
}
