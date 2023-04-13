package common

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
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

func Find[T any](lst *[]T, callback func(item *T) bool) *T {
	for _, item := range *lst {
		if callback(&item) {
			return &item
		}
	}
	return nil
}

func Fetch[T any](method string, url string) (*T, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := new(T)
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
