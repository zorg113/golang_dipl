package common

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
)

var validatePattern *regexp.Regexp

var IPAlreadyExist = errors.New("IP address already exists") //nolint:revive,stylecheck

func init() {
	validatePattern = regexp.MustCompile(`(?m)^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
}

func InitHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
}

func ValidateIP(network entity.IPNetwork) bool {
	if !isCorrectIP(network.IP) || !isCorrectMask(network.Mask) {
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
	if !isCorrectIP(request.IP) {
		return false
	}
	return true
}
