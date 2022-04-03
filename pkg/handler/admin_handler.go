package handler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/handler/converter"
	"github.com/kujilabo/cocotola-translator-api/pkg/handler/entity"
	handlerhelper "github.com/kujilabo/cocotola-translator-api/pkg/handler/helper"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	"github.com/kujilabo/cocotola-translator-api/pkg/usecase"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/log"
)

type AdminHandler interface {
	FindTranslations(c *gin.Context)
	FindTranslationByTextAndPos(c *gin.Context)
	FindTranslationByText(c *gin.Context)
	AddTranslation(c *gin.Context)
	UpdateTranslation(c *gin.Context)
	RemoveTranslation(c *gin.Context)
	ExportTranslations(c *gin.Context)
}

type adminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) AdminHandler {
	return &adminHandler{adminUsecase: adminUsecase}
}

func (h *adminHandler) FindTranslations(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslations")

	handlerhelper.HandleFunction(c, func() error {
		param := entity.TranslationFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.adminUsecase.FindTranslationsByFirstLetter(ctx, domain.Lang2JA, param.Letter)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationFindResposne(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) FindTranslationByTextAndPos(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslationByTextAndPos")

	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetString(c, "text")

		pos, err := ginhelper.GetInt(c, "pos")
		if err != nil {
			return err
		}

		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}
		result, err := h.adminUsecase.FindTranslationByTextAndPos(ctx, domain.Lang2JA, text, wordPos)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationResposne(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) FindTranslationByText(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {

		text := ginhelper.GetString(c, "text")
		results, err := h.adminUsecase.FindTranslationByText(ctx, domain.Lang2JA, text)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationListResposne(ctx, results)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) AddTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		param := entity.TranslationAddParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationAddParameter(ctx, &param)
		if err != nil {
			return err
		}

		if err := h.adminUsecase.AddTranslation(ctx, parameter); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) UpdateTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetString(c, "text")

		pos, err := ginhelper.GetInt(c, "pos")
		if err != nil {
			return err
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}

		param := entity.TranslationUpdateParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationUpdateParameter(ctx, &param)
		if err != nil {
			return err
		}

		if err := h.adminUsecase.UpdateTranslation(ctx, domain.Lang2JA, text, wordPos, parameter); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) RemoveTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetString(c, "text")

		pos, err := ginhelper.GetInt(c, "pos")
		if err != nil {
			return err
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}

		if err := h.adminUsecase.RemoveTranslation(ctx, domain.Lang2JA, text, wordPos); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) ExportTranslations(c *gin.Context) {
	handlerhelper.HandleFunction(c, func() error {
		csvStruct := [][]string{
			{"name", "address", "phone"},
			{"Ram", "Tokyo", "1236524"},
			{"Shaym", "Beijing", "8575675484"},
		}
		b := new(bytes.Buffer)
		w := csv.NewWriter(b)
		if err := w.WriteAll(csvStruct); err != nil {
			return err
		}
		if _, err := c.Writer.Write(b.Bytes()); err != nil {
			return err
		}
		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *adminHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	if errors.Is(err, service.ErrTranslationAlreadyExists) {
		logger.Warnf("adminHandler. err: %v", err)
		c.JSON(http.StatusConflict, gin.H{"message": "Translation already exists"})
		return true
	}
	logger.Errorf("adminHandler. err: %v", err)
	return false
}
