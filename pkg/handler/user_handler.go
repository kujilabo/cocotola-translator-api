package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	handlerhelper "github.com/kujilabo/cocotola-translator-api/pkg/handler/helper"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	"github.com/kujilabo/cocotola-translator-api/pkg/usecase"
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
	handlerhelper.HandleFunction(c, func() error {
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
