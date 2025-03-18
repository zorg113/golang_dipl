package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type Bucket struct {
	service *service.Authorization
	log     *zerolog.Logger
}

func NewBucket(service *service.Authorization, log *zerolog.Logger) *Bucket {
	return &Bucket{service: service, log: log}
}
func (b *Bucket) ResetBucket(w http.ResponseWriter, r *http.Request) {
	b.log.Info().Msg("Reset bucket handler by POST /bucket/reset called")
	initHeaders(w)
	var request entity.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request.Password = "empty"
	isValidate := ValidateRequest(request)
	if !isValidate {
		b.log.Info().Msg("Invalid input request from client")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte("resetLogin=true"))
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	isIpReset := b.service.ResetIpBucket(request.Ip)
	if !isIpReset {
		b.log.Info().Msg("Failed to reset IP bucket")
		_, err = w.Write([]byte("resetIp=false"))
		if err != nil {
            b.log.Error().Err(err).Msg("Failed to write response")
            w.WriteHeader(http.StatusInternalServerError)
        }
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("resetIp=true"))
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
