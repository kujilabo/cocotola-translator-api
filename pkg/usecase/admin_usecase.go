package usecase

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/log"
)

type AdminUsecase interface {
	FindTranslationsByFirstLetter(ctx context.Context, lang domain.Lang2, firstLetter string) ([]domain.Translation, error)

	FindTranslationByTextAndPos(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)

	FindTranslationByText(ctx context.Context, lang domain.Lang2, text string) ([]domain.Translation, error)

	AddTranslation(ctx context.Context, param service.TranslationAddParameter) error

	UpdateTranslation(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error

	RemoveTranslation(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) error
}

type adminUsecase struct {
	rf service.RepositoryFactory
}

func NewAdminUsecase(rf service.RepositoryFactory) AdminUsecase {
	return &adminUsecase{rf: rf}
}

func (u *adminUsecase) FindTranslationsByFirstLetter(ctx context.Context, lang domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResults, err := customRepo.FindByFirstLetter(ctx, lang, firstLetter)
	if err != nil {
		return nil, err
	}
	azureRepo, err := u.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResults, err := azureRepo.FindByFirstLetter(ctx, lang, firstLetter)
	if err != nil {
		return nil, err
	}

	makeKey := func(text string, pos domain.WordPos) string {
		return text + "_" + strconv.Itoa(int(pos))
	}
	resultMap := make(map[string]domain.Translation)
	for _, c := range customResults {
		key := makeKey(c.GetText(), c.GetPos())
		resultMap[key] = c
	}
	for _, a := range azureResults {
		key := makeKey(a.GetText(), a.GetPos())
		if _, ok := resultMap[key]; !ok {
			resultMap[key] = a
		}
	}

	results := make([]domain.Translation, 0)
	for _, v := range resultMap {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].GetText() < results[j].GetText() })

	return results, nil
}

func (u *adminUsecase) FindTranslationByTextAndPos(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResult, err := customRepo.FindByTextAndPos(ctx, lang, text, pos)
	if err == nil {
		return customResult, nil
	}
	if !errors.Is(err, service.ErrTranslationNotFound) {
		return nil, err
	}

	azureRepo, err := u.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResult, err := azureRepo.FindByTextAndPos(ctx, lang, text, pos)
	if err != nil {
		return nil, err
	}
	return azureResult, nil
}

func (u *adminUsecase) FindTranslationByText(ctx context.Context, lang domain.Lang2, text string) ([]domain.Translation, error) {
	logger := log.FromContext(ctx)
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResults, err := customRepo.FindByText(ctx, lang, text)
	if err != nil {
		return nil, err
	}
	azureRepo, err := u.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResults, err := azureRepo.FindByText(ctx, lang, text)
	if err != nil {
		return nil, err
	}

	makeKey := func(text string, pos domain.WordPos) string {
		return text + "_" + strconv.Itoa(int(pos))
	}
	resultMap := make(map[string]domain.Translation)
	for _, c := range customResults {
		key := makeKey(c.GetText(), c.GetPos())
		resultMap[key] = c
	}
	for _, a := range azureResults {
		key := makeKey(a.GetText(), a.GetPos())
		if _, ok := resultMap[key]; !ok {
			resultMap[key] = a
			logger.Infof("translation: %v", a)
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

func (u *adminUsecase) AddTranslation(ctx context.Context, param service.TranslationAddParameter) error {
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	if err := customRepo.Add(ctx, param); err != nil {
		return err
	}
	return nil
}

func (u *adminUsecase) UpdateTranslation(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	translationFound := true
	if _, err := customRepo.FindByTextAndPos(ctx, lang, text, pos); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			translationFound = false
		} else {
			return err
		}
	}

	if translationFound {
		if err := customRepo.Update(ctx, lang, text, pos, param); err != nil {
			return err
		}
		return nil
	}

	paramToAdd, err := service.NewTransalationAddParameter(text, pos, lang, param.GetTranslated())
	if err != nil {
		return err
	}
	if err := customRepo.Add(ctx, paramToAdd); err != nil {
		return err
	}
	return nil
}

func (u *adminUsecase) RemoveTranslation(ctx context.Context, lang domain.Lang2, text string, pos domain.WordPos) error {
	customRepo, err := u.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	if err := customRepo.Remove(ctx, lang, text, pos); err != nil {
		return err
	}
	return nil
}
