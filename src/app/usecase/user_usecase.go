package usecase

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
)

type UserUsecase interface {
	DictionaryLookup(ctx context.Context, fromLang, toLang domain.Lang2, text string) ([]domain.Translation, error)

	DictionaryLookupWithPos(ctx context.Context, fromLang, toLang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)
}

type userUsecase struct {
	rf                     service.RepositoryFactory
	azureTranslationClient service.AzureTranslationClient
}

func NewUserUsecase(rf service.RepositoryFactory, azureTranslationClient service.AzureTranslationClient) UserUsecase {
	return &userUsecase{
		rf:                     rf,
		azureTranslationClient: azureTranslationClient,
	}
}

func (u *userUsecase) selectMaxConfidenceTranslations(ctx context.Context, in []service.AzureTranslation) (map[domain.WordPos]service.AzureTranslation, error) {
	results := make(map[domain.WordPos]service.AzureTranslation)
	for _, i := range in {
		if _, ok := results[i.Pos]; !ok {
			results[i.Pos] = i
		} else if i.Confidence > results[i.Pos].Confidence {
			results[i.Pos] = i
		}
	}
	return results, nil
}

func (u *userUsecase) customDictionaryLookup(ctx context.Context, text string, fromLang, toLang domain.Lang2) ([]domain.Translation, error) {
	// repo, err := u.rf.NewAzureTranslationRepository()
	// if err != nil {
	// 	return nil, err
	// }
	customRepo := u.rf.NewCustomTranslationRepository(ctx)

	customContained, err := customRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if !customContained {
		return nil, service.ErrTranslationNotFound
	}

	customResults, err := customRepo.FindByText(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	return customResults, nil
}

func (u *userUsecase) azureDictionaryLookup(ctx context.Context, fromLang, toLang domain.Lang2, text string) ([]service.AzureTranslation, error) {
	// repo, err := t.repo(t.db)
	// if err != nil {
	// 	return nil, err
	// }

	azureRepo := u.rf.NewAzureTranslationRepository(ctx)
	azureContained, err := azureRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if azureContained {
		azureResults, err := azureRepo.Find(ctx, toLang, text)
		if err != nil {
			return nil, err
		}
		return azureResults, nil
	}

	azureResults, err := u.azureTranslationClient.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	if len(azureResults) == 0 {
		return azureResults, nil
	}

	if err := azureRepo.Add(ctx, toLang, text, azureResults); err != nil {
		return nil, xerrors.Errorf("failed to add auzre_translation. err: %w", err)
	}

	return azureResults, nil
}

func (u *userUsecase) DictionaryLookup(ctx context.Context, fromLang, toLang domain.Lang2, text string) ([]domain.Translation, error) {
	// find translations from custom reopository
	customResults, err := u.customDictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil && !errors.Is(err, service.ErrTranslationNotFound) {
		return nil, err
	}
	// if !errors.Is(err, service.ErrTranslationNotFound) {
	// 	return customResults, err
	// }

	// find translations from azure
	azureResults, err := u.azureDictionaryLookup(ctx, fromLang, toLang, text)
	if err != nil {
		return nil, err
	}
	azureResultMap, err := u.selectMaxConfidenceTranslations(ctx, azureResults)
	if err != nil {
		return nil, err
	}
	makeKey := func(text string, pos domain.WordPos) string {
		return text + "_" + strconv.Itoa(int(pos))
	}
	resultMap := make(map[string]domain.Translation)

	// insert customResults into resultMap
	for _, c := range customResults {
		key := makeKey(c.GetText(), c.GetPos())
		resultMap[key] = c
	}

	// insert azureResultMap into resultMap
	for _, a := range azureResultMap {
		key := makeKey(text, a.Pos)
		if _, ok := resultMap[key]; !ok {
			result, err := domain.NewTranslation(1, time.Now(), time.Now(), text, a.Pos, fromLang, a.Target, "azure")
			if err != nil {
				return nil, err
			}
			resultMap[key] = result
		}
	}

	// convert map to list
	results := make([]domain.Translation, 0)
	for _, v := range resultMap {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].GetPos() < results[j].GetPos() })

	return results, nil
}

func (u *userUsecase) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	results, err := u.DictionaryLookup(ctx, fromLang, toLang, text)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.GetPos() == pos {
			return r, nil
		}
	}
	return nil, service.ErrTranslationNotFound
}
