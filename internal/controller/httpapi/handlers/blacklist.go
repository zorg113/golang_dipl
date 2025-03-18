package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

type BlackList struct {
	service *service.BlackList
	log     *zerolog.Logger
}

func NewBlackList(service *service.BlackList, log *zerolog.Logger) *BlackList {
	return &BlackList{service: service, log: log}
}

func (b *BlackList) AddIP(w http.ResponseWriter, r *http.Request /*router params*/) {
	b.log.Info().Msg("Add IP to black list handler by POST /blacklist/add called")
	initHeaders(w)
	var inIP entity.IpNetwork
	err := json.NewDecoder(r.Body).Decode(&inIP)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateIP(inIP) {
		b.log.Info().Msg("Invalid IP format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = b.service.AddIP(inIP)
	if err != nil {
		if err.Error() == ipAlreadyExist.Error() {
			b.log.Info().Msg("IP already exist in black list")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("Failed to write response")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		b.log.Error().Err(err).Msg("Failed to add IP to black list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // ????
}

func (b *BlackList) DeleteIP(w http.ResponseWriter, r *http.Request /* router params */) {
	b.log.Info().Msg("Remove IP from black list handler by DELETE /blacklist/remove called")
	initHeaders(w)
	var inIP entity.IpNetwork
	err := json.NewDecoder(r.Body).Decode(&inIP)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to decode request body: ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateIP(inIP) {
		b.log.Info().Msg("Invalid IP format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = b.service.DeleteIP(inIP)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to remove IP from black list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (b *BlackList) GetIPs(w http.ResponseWriter, r *http.Request /* router params */) {
	b.log.Info().Msg("Get black list IPs handler by GET /blacklist/get called")
	initHeaders(w)
	ips, err := b.service.GetIPs()
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to get black list IPs")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(ips)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to marshal response to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
