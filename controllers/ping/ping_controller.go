package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_utils-go/logger"
	"net/http"
)

func Ping(c *gin.Context) {
	logger.Info("pong")
	c.String(http.StatusOK, "pong")
}
