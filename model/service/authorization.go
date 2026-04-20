package service

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"golang.org/x/time/rate"
)

type Authorization struct {
	ipBucketStorage       sync.Map // map[string]*RateLimiter
	loginBucketStorage    sync.Map // map[string]*RateLimiter
	passwordBucketStorage sync.Map // map[string]*RateLimiter
	blackList             *BlackList
	whiteList             *WhiteList
	conf                  *config.Config
	log                   *zerolog.Logger
}

func NewAuthorization(bList *BlackList, wList *WhiteList, cfg *config.Config, logger *zerolog.Logger) *Authorization {
	auth := &Authorization{
		blackList: bList,
		whiteList: wList,
		conf:      cfg,
		log:       logger,
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
	isIPInBlackList, err := a.checkIPByNetworkList(request.IP, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIPInBlackList {
		return false, nil
	}
	a.log.Info().Msg("Check IP in white list")
	ipNetworkList, err = a.whiteList.GetIPs()
	if err != nil {
		return false, err
	}
	isLoginInWhiteList, err := a.checkIPByNetworkList(request.IP, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isLoginInWhiteList {
		return true, nil
	}
	a.log.Info().Msg("Check ip in bucket")
	isAllow := true
	allow := a.getPermissionInBucket(request.IP, &a.ipBucketStorage, a.conf.Bucket.IPLimit)
	if !allow {
		isAllow = false
	}
	a.log.Info().Msg("Check login in bucket")
	allow = a.getPermissionInBucket(request.Login, &a.loginBucketStorage, a.conf.Bucket.LoginLimit)
	if !allow {
		isAllow = allow
	}
	a.log.Info().Msg("Check password in bucket")
	allow = a.getPermissionInBucket(request.Password, &a.passwordBucketStorage, a.conf.Bucket.PasswordLimit)
	if !allow {
		isAllow = false
	}
	return isAllow, nil
}

func (a *Authorization) newBucket(limit int) *RateLimiter {
	r := NewRateLimiter(rate.Limit(float64(limit)/time.Duration.Seconds(60*time.Second)), limit)
	return r
}

func (a *Authorization) checkIPByNetworkList(ip string, ipNetworkList []entity.IPNetwork) (bool, error) {
	for _, network := range ipNetworkList {
		prefix, err := GetPrefix(ip, network.Mask)
		if err != nil {
			return false, err
		}
		if prefix == network.IP {
			return true, nil
		}
	}
	return false, nil
}

func (a *Authorization) getPermissionInBucket(request string, bucketStorage *sync.Map, limit int) bool {
	limiter, ok := bucketStorage.Load(request)
	if ok {
		return limiter.(*RateLimiter).Allow()
	}
	actual, _ := bucketStorage.LoadOrStore(request, a.newBucket(limit))
	return actual.(*RateLimiter).Allow()
}

func (a *Authorization) ResetLoginInBucket(login string) bool {
	_, ok := a.loginBucketStorage.LoadAndDelete(login)
	return ok
}

func (a *Authorization) ResetIPBucket(ip string) bool {
	_, ok := a.ipBucketStorage.LoadAndDelete(ip)
	return ok
}

func (a *Authorization) ResetPasswordInBucket(password string) bool {
	_, ok := a.passwordBucketStorage.LoadAndDelete(password)
	return ok
}

func (a *Authorization) deleteUnusedBucket() {
	interval := time.Duration(a.conf.Bucket.ResetBucketInterval) * time.Second
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		a.evictStale(&a.ipBucketStorage)
		a.evictStale(&a.loginBucketStorage)
		a.evictStale(&a.passwordBucketStorage)
	}
}

func (a *Authorization) evictStale(storage *sync.Map) {
	threshold := time.Duration(a.conf.Bucket.ResetBucketInterval) * time.Second

	storage.Range(func(key, value any) bool {
		limiter := value.(*RateLimiter)
		if time.Since(limiter.LastEvent()) > threshold {
			storage.Delete(key)
		}
		return true // продолжать итерацию
	})
}
