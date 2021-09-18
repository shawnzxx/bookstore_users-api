package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_utils-go/logger"
	"os"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	s := fmt.Sprintf("Application running in environment: %s and on port: %s", os.Getenv("ENV"), os.Getenv("PORT"))
	logger.Info(s)
	router.Run(":" + os.Getenv("PORT"))
}
