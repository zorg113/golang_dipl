package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type Authorization struct {
	service *service.Authorization
	log     *zerolog.Logger
}

func NewAuthorization(service *service.Authorization, log *zerolog.Logger) *Authorization {
	return &Authorization{service: service, log: log}
}

func (a *Authorization) AuthorizationHanler(w http.ResponseWriter, r *http.Request /* ps httprouter.Params*/) {
	a.log.Info().Msg("Authorization handler by POST /auth/check/ called")
	initHeaders(w)
	var request entity.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateRequest(request)
	if !isValidate {
		a.log.Info().Msg("Invalid input request from client")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isAlllowed, err := a.service.Authorization(request)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to authorization")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAlllowed {
		a.log.Info().Msg("Request is not allowed")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok=false"))
		if err != nil {
			a.log.Error().Err(err).Msg("Failed to write response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok=true"))
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
