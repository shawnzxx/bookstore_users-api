package ping

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_utils-go/app_logger"
	"net"
	"net/http"
	"os"
)

const (
	AuthServiceHost = "AUTH_SERVICE_HOST"
)

var (
	logger = app_logger.GetLogger()
	ipv4   string
	port   int
)

func Ping(c *gin.Context) {
	host := os.Getenv(AuthServiceHost)

	ips, _ := net.LookupHost(host)
	pong := fmt.Sprintf("oauth api ips %v", ips)
	c.String(http.StatusOK, pong)
}
