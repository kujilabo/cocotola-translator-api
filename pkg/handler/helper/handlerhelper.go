package handlerhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-translator-api/pkg_lib/log"
)

func HandleFunction(c *gin.Context, fn func() error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	logger.Info("")
	if err := fn(); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}
