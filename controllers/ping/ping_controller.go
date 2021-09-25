package ping

import (
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
	c.String(http.StatusOK, "pong")
	host := os.Getenv(AuthServiceHost)

	//print out service IPs
	ips, err := net.LookupHost(host)
	logger.Info("oauth api LookupHost return: %v, error: %v", ips, err)
}
