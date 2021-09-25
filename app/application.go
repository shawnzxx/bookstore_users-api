package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_utils-go/app_logger"
	"net"
	"os"
)

var (
	logger          = app_logger.GetLogger()
	router          = gin.Default()
	env, ipv4, port string
)

func StartApplication() {
	printOutServiceInfo()
	mapUrls()
	router.Run(":" + port)
}

func printOutServiceInfo() {
	//get local or container's host name
	hostName, _ := os.Hostname()
	logger.Info("users api's hostname: %s", hostName)
	//print out service IPs
	ips, err := net.LookupHost(hostName)
	if err != nil {
		logger.Info("Can not find ips list for the host %v", err)
	}
	//print out service ip, env, port
	env = os.Getenv("ENV")
	ipv4 = ips[0]
	port = os.Getenv("PORT")
	logger.Info("bookstore_users-api running on %s environment, ip is %s, port is %s", env, ipv4, port)
}
