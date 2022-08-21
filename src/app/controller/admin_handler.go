package controller

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-translator-api/src/app/controller/converter"
	"github.com/kujilabo/cocotola-translator-api/src/app/controller/entity"
	handlerhelper "github.com/kujilabo/cocotola-translator-api/src/app/controller/helper"
	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	"github.com/kujilabo/cocotola-translator-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-translator-api/src/lib/log"
)

type AdminHandler interface {
	FindTranslationsByFirstLetter(c *gin.Context)
	FindTranslationByTextAndPos(c *gin.Context)
	FindTranslationsByText(c *gin.Context)
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

// FindTranslationsByFirstLetter godoc
// @Summary     find translations with first letter
// @Description find translations with first letter
// @Tags        translator
// @Accept      json
// @Produce     json
// @Param       param body entity.TranslationFindParameter true "parameter to find translations"
// @Success     200 {object} entity.TranslationFindResponse
// @Failure     400
// @Failure     401
// @Router      /v1/admin/find [post]
// @Security    BasicAuth
func (h *adminHandler) FindTranslationsByFirstLetter(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslations")

	handlerhelper.HandleFunction(c, func() error {
		param := entity.TranslationFindParameterHTTPEntity{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		if len(param.Letter) != 1 {
			c.Status(http.StatusBadRequest)
			return nil
		}

		results, err := h.adminUsecase.FindTranslationsByFirstLetter(ctx, domain.Lang2JA, param.Letter)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationFindResposne(ctx, results)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

// FindTranslationByTextAndPos godoc
// @Summary     find translations with text and pos
// @Description find translations with text and pos
// @Tags        translator
// @Accept      json
// @Produce     json
// @Param       text path string true "text"
// @Param       pos path int true "pos"
// @Success     200 {object} entity.Translation
// @Failure     400
// @Failure     401
// @Router      /v1/admin/text/{text}/pos/{pos} [get]
// @Security    BasicAuth
func (h *adminHandler) FindTranslationByTextAndPos(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslationByTextAndPos")

	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
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

// FindTranslationsByText godoc
// @Summary     find translations with text
// @Description find translations with text
// @Tags        translator
// @Accept      json
// @Produce     json
// @Param       text path string true "text"
// @Success     200 {object} entity.Translation
// @Failure     400
// @Failure     401
// @Router      /v1/admin/text/{text} [get]
// @Security    BasicAuth
func (h *adminHandler) FindTranslationsByText(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetStringFromPath(c, "text")
		results, err := h.adminUsecase.FindTranslationByText(ctx, domain.Lang2JA, text)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationFindResposne(ctx, results)
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
		param := entity.TranslationAddParameterHTTPEntity{}
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
		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
		if err != nil {
			return err
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}

		param := entity.TranslationUpdateParameterHTTPEntity{}
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
		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
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
