package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/handler/converter"
	handlerhelper "github.com/kujilabo/cocotola-translator-api/pkg/handler/helper"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	"github.com/kujilabo/cocotola-translator-api/pkg/usecase"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/log"
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

func (h *userHandler) DictionaryLookup(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	handlerhelper.HandleFunction(c, func() error {
		text := ginhelper.GetStringFromQuery(c, "text")
		posS := ginhelper.GetStringFromQuery(c, "pos")
		if len(posS) == 0 {
			results, err := h.userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, text)
			if err != nil {
				return err
			}

			response, err := converter.ToTranslationFindResposne(ctx, results)
			if err != nil {
				return err
			}

			logger.Infof("response: %+v", response)
			c.JSON(http.StatusOK, response)
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

		response, err := converter.ToTranslationFindResposne(ctx, []domain.Translation{result})
		if err != nil {
			return err
		}

		logger.Infof("response: %+v", response)
		c.JSON(http.StatusOK, response)

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
