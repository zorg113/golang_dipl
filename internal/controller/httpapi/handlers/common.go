package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
)

var validatePattern *regexp.Regexp
var ipAlreadyExist = errors.New("IP address already exists")

func init() {
	validatePattern = regexp.MustCompile(`(?m)^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
}

func initHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
}

func ValidateIP(network entity.IpNetwork) bool {
	if !isCorrectIP(network.Ip) && !isCorrectMask(network.Mask) {
		return false
	}
	return true
}

func isCorrectIP(ip string) bool {
	return validatePattern.MatchString(ip)
}

func isCorrectMask(mask string) bool {
	return validatePattern.MatchString(mask)
}

func ValidateRequest(request entity.Request) bool {
	if request.Login == "" || request.Password == "" {
		return false
	}
	if !isCorrectIP(request.Ip) {
		return false
	}
	return true
}
