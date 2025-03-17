package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type WhiteList struct {
	service *service.WhiteList
	log     *zerolog.Logger
}

func NewWhiteList(service *service.WhiteList, log *zerolog.Logger) *WhiteList {
	return &WhiteList{service: service, log: log}
}

func (wl *WhiteList) AddIP(w http.ResponseWriter, r *http.Request) {
	wl.log.Info().Msg("Add IP in whitelist by POST /auth/whitelist")
	initHeaders(w)
	var inIP entity.IpNetwork
	err := json.NewDecoder(r.Body).Decode(&inIP)
	if err != nil {
		wl.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateIP(inIP) {
		wl.log.Info().Msg("Invalid IP format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = wl.service.AddIP(inIP)
	if err != nil {
		if err.Error() == ipAlreadyExist.Error() {
			wl.log.Info().Msg("IP already exist in white list")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				wl.log.Error().Err(err).Msg("Failed to write response")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		wl.log.Error().Err(err).Msg("Failed to add IP to white list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (wl *WhiteList) DeleteIP(w http.ResponseWriter, r *http.Request /* router params */) {
	wl.log.Info().Msg("Remove IP from whitelist by DELETE /auth/whitelist/remove called")
	initHeaders(w)
	var inIP entity.IpNetwork
	err := json.NewDecoder(r.Body).Decode(&inIP)
	if err != nil {
		wl.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateIP(inIP) {
		wl.log.Info().Msg("Invalid IP format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = wl.service.DeleteIP(inIP)
	if err != nil {
		wl.log.Error().Err(err).Msg("Failed to remove IP from white list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (wl *WhiteList) GetIPs(w http.ResponseWriter, r *http.Request) {
	wl.log.Info().Msg("Get white list by GET /auth/whitelist called")
	initHeaders(w)
	ipList, err := wl.service.GetIPs()
	if err != nil {
		wl.log.Error().Err(err).Msg("Failed to get white list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(ipList)
	if err != nil {
		wl.log.Error().Err(err).Msg("Failed to marshal response to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
