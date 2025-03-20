package service

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"golang.org/x/time/rate"
)

type Authorization struct {
	ipBucketStorage       map[string]*RateLimiter
	loginBucketStorage    map[string]*RateLimiter
	passwordBucketStorage map[string]*RateLimiter
	blackList             *BlackList
	whiteList             *WhiteList
	conf                  *config.Config
	log                   *zerolog.Logger
}

func NewAuthorization(bList *BlackList, wList *WhiteList,
	cfg *config.Config, logger *zerolog.Logger) *Authorization {
	ipBucketStorage := make(map[string]*RateLimiter)
	loginBucketStorage := make(map[string]*RateLimiter)
	passwordBucketStorage := make(map[string]*RateLimiter)
	auth := &Authorization{ipBucketStorage: ipBucketStorage,
		loginBucketStorage:    loginBucketStorage,
		passwordBucketStorage: passwordBucketStorage,
		blackList:             bList,
		whiteList:             wList,
		conf:                  cfg,
		log:                   logger,
	}
	go auth.deleteUnusedBucket()
	return auth
}

func (a *Authorization) Authorization(request entity.Request) (bool, error) {
	a.log.Info().Msg("Check ip in black list")
	ipNetworkList, err := a.blackList.GetIPs()
	if err != nil {
		return false, err
	}
	isIpInBlackList, err := a.checkIpByNetworkList(request.Ip, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIpInBlackList {
		return false, nil
	}
	a.log.Info().Msg("Check login in white list")
	ipNetworkList, err = a.whiteList.GetIPs()
	if err != nil {
		return false, err
	}
	isLoginInWhiteList, err := a.checkIpByNetworkList(request.Login, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isLoginInWhiteList {
		return true, nil
	}
	a.log.Info().Msg("Check ip in bucket")
	isAllow := true
	allow := a.getPermissionInBucket(request.Ip, a.ipBucketStorage, a.conf.Bucket.IpLimit)
	if !allow {
		isAllow = false
	}
	a.log.Info().Msg("Chek password in bucket")
	allow = a.getPermissionInBucket(request.Password, a.passwordBucketStorage, a.conf.Bucket.PasswordLimit)
	if !allow {
		isAllow = allow
	}
	a.log.Info().Msg("Check login in bucket")
	allow = a.getPermissionInBucket(request.Login, a.loginBucketStorage, a.conf.Bucket.LoginLimit)
	if !allow {
		isAllow = allow
	}
	return isAllow, nil
}

func (a *Authorization) newBucket(limit int) *RateLimiter {
	r := NewRateLimiter(rate.Limit(float64(limit)/time.Duration.Seconds(60*time.Second)), limit)
	return r
}

func (a *Authorization) checkIpByNetworkList(ip string, ipNetworkList []entity.IpNetwork) (bool, error) {
	for _, network := range ipNetworkList {
		prefix, err := GetPrefix(network.Ip, network.Mask)
		if err != nil {
			return false, err
		}
		if prefix == network.Ip {
			return true, nil
		}
	}
	return false, nil
}

func (a *Authorization) getPermissionInBucket(request string, bucketStorage map[string]*RateLimiter, limit int) bool {
	limiter, ok := bucketStorage[request]
	if !ok {
		bucketStorage[request] = a.newBucket(limit)
		allow := bucketStorage[request].Allow()
		return allow
	}
	allow := limiter.Allow()
	return allow
}

func (a *Authorization) ResetLoginInBucket(login string) bool {
	_, ok := a.loginBucketStorage[login]
	if !ok {
		return false
	}
	delete(a.loginBucketStorage, login)
	return true
}

func (a *Authorization) ResetIpBucket(ip string) bool {
	_, ok := a.ipBucketStorage[ip]
	if !ok {
		return false
	}
	delete(a.ipBucketStorage, ip)
	return true
}

func (a *Authorization) deleteUnusedBucket() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		<-ticker.C
		for ip, lim := range a.ipBucketStorage {
			if time.Since(lim.LastEvent) > time.Duration(a.conf.Bucket.ResetBucketInterval)*time.Second {
				delete(a.ipBucketStorage, ip)
			}
		}
		for login, lim := range a.loginBucketStorage {
			if time.Since(lim.LastEvent) > time.Duration(a.conf.Bucket.ResetBucketInterval)*time.Second {
				delete(a.loginBucketStorage, login)
			}
		}
		for password, lim := range a.passwordBucketStorage {
			if time.Since(lim.LastEvent) > time.Duration(a.conf.Bucket.ResetBucketInterval)*time.Second {
				delete(a.passwordBucketStorage, password)
			}
		}
	}
}
