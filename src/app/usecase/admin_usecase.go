//go:generate mockery --output mock --name AdminUsecase
package usecase

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	liberrors "github.com/kujilabo/cocotola-translator-api/src/lib/errors"
	"github.com/kujilabo/cocotola-translator-api/src/lib/log"
)

type AdminUsecase interface {
	FindTranslationsByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error)

	FindTranslationByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)

	FindTranslationByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error)

	AddTranslation(ctx context.Context, param service.TranslationAddParameter) error

	UpdateTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error

	RemoveTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) error
}

type AdminPresenter interface {
	WriteTranslations(ctx context.Context, translations []domain.Translation) error
	WriteTranslation(ctx context.Context, translation domain.Translation) error
}

type adminUsecase struct {
	rf service.RepositoryFactory
}

func NewAdminUsecase(rf service.RepositoryFactory) AdminUsecase {
	return &adminUsecase{rf: rf}
}

func (u *adminUsecase) FindTranslationsByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	customRepo := u.rf.NewCustomTranslationRepository(ctx)
	customResults, err := customRepo.FindByFirstLetter(ctx, lang2, firstLetter)
	if err != nil {
		return nil, err
	}
	azureRepo := u.rf.NewAzureTranslationRepository(ctx)
	azureResults, err := azureRepo.FindByFirstLetter(ctx, lang2, firstLetter)
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

func (u *adminUsecase) FindTranslationByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	customRepo := u.rf.NewCustomTranslationRepository(ctx)
	customResult, err := customRepo.FindByTextAndPos(ctx, lang2, text, pos)
	if err == nil {
		return customResult, nil
	}
	if !errors.Is(err, service.ErrTranslationNotFound) {
		return nil, err
	}

	azureRepo := u.rf.NewAzureTranslationRepository(ctx)
	azureResult, err := azureRepo.FindByTextAndPos(ctx, lang2, text, pos)
	if err != nil {
		return nil, err
	}
	return azureResult, nil
}

func (u *adminUsecase) FindTranslationByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error) {
	logger := log.FromContext(ctx)
	customRepo := u.rf.NewCustomTranslationRepository(ctx)
	customResults, err := customRepo.FindByText(ctx, lang2, text)
	if err != nil {
		return nil, err
	}
	azureRepo := u.rf.NewAzureTranslationRepository(ctx)
	azureResults, err := azureRepo.FindByText(ctx, lang2, text)
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
	customRepo := u.rf.NewCustomTranslationRepository(ctx)
	if err := customRepo.Add(ctx, param); err != nil {
		return err
	}
	return nil
}

func (u *adminUsecase) UpdateTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	customRepo := u.rf.NewCustomTranslationRepository(ctx)

	translationFound := true
	if _, err := customRepo.FindByTextAndPos(ctx, lang2, text, pos); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			translationFound = false
		} else {
			return err
		}
	}

	if translationFound {
		if err := customRepo.Update(ctx, lang2, text, pos, param); err != nil {
			return err
		}
		return nil
	}

	paramToAdd, err := service.NewTransalationAddParameter(text, pos, lang2, param.GetTranslated())
	if err != nil {
		return err
	}
	if err := customRepo.Add(ctx, paramToAdd); err != nil {
		return err
	}
	return nil
}

func (u *adminUsecase) RemoveTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) error {
	customRepo := u.rf.NewCustomTranslationRepository(ctx)
	if err := customRepo.Remove(ctx, lang2, text, pos); err != nil {
		return liberrors.Errorf("failed to customRepo.Remove in adminUsecase.RemoveTranslation. err: %w", err)
	}
	return nil
}
