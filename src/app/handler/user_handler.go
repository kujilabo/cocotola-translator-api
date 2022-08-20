package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-translator-api/src/app/domain"
	handlerhelper "github.com/kujilabo/cocotola-translator-api/src/app/handler/helper"
	"github.com/kujilabo/cocotola-translator-api/src/app/presenter"
	"github.com/kujilabo/cocotola-translator-api/src/app/service"
	"github.com/kujilabo/cocotola-translator-api/src/app/usecase"
	"github.com/kujilabo/cocotola-translator-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-translator-api/src/lib/log"
)

type UserHandler interface {
	DictionaryLookup(c *gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{userUsecase: userUsecase}
}

// DictionaryLookup godoc
// @Summary     dictionary lookup
// @Description dictionary lookup
// @Tags        translator
// @Accept      json
// @Produce     json
// @Param       text query string true "text"
// @Param       pos query int false "pos"
// @Success     200 {object} entity.Translation
// @Failure     400
// @Failure     401
// @Router      /v1/user/dictionary/lookup [get]
// @Security    BasicAuth
func (h *userHandler) DictionaryLookup(c *gin.Context) {
	ctx := c.Request.Context()
	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetStringFromQuery(c, "text")
		if len(text) <= 1 {
			c.Status(http.StatusBadRequest)
			return nil
		}

		posS := ginhelper.GetStringFromQuery(c, "pos")
		if len(posS) == 0 {
			results, err := h.userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, text)
			if err != nil {
				return err
			}

			userPresenter := presenter.NewUserPresenter(c)
			userPresenter.WriteTranslations(ctx, results)
			return nil
		}

		posI, err := strconv.Atoi(posS)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		pos, err := domain.NewWordPos(posI)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.userUsecase.DictionaryLookupWithPos(ctx, domain.Lang2EN, domain.Lang2JA, text, pos)
		if err != nil {
			return err
		}

		userPresenter := presenter.NewUserPresenter(c)
		userPresenter.WriteTranslation(ctx, result)
		return nil
	}, h.errorHandle)
}

func (h *userHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	if errors.Is(err, service.ErrTranslationAlreadyExists) {
		logger.Warnf("userHandler. err: %v", err)
		c.JSON(http.StatusConflict, gin.H{"message": "Translation already exists"})
		return true
	}
	logger.Errorf("userHandler. err: %v", err)
	return false
}
