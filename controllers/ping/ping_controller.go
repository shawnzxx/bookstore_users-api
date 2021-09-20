package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_utils-go/logger"
	"net"
	"net/http"
	"os"
)

const (
	AuthServiceHost = "AUTH_SERVICE_HOST"
)

var (
	_logger = logger.GetLogger()
	ipv4    string
	port    int
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
	host := os.Getenv(AuthServiceHost)

	//print out service IPs
	ips, err := net.LookupHost(host)
	_logger.Printf("oauth api LookupHost return: %v, error: %v\n", ips, err)
}
