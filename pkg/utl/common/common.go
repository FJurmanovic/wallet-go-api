package common

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type RouteGroups struct {
	Api               *gin.RouterGroup
	Auth              *gin.RouterGroup
	Wallet            *gin.RouterGroup
	WalletHeader      *gin.RouterGroup
	Transaction       *gin.RouterGroup
	TransactionType   *gin.RouterGroup
	Subscription      *gin.RouterGroup
	SubscriptionType  *gin.RouterGroup
	TransactionStatus *gin.RouterGroup
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Error occured. %v", err)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
