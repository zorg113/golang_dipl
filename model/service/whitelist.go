package service

import (
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
)

type WhiteListStore interface {
	Add(prefix, mask string) error
	Remove(prefix, mask string) error
	GetIPs() ([]entity.IpNetwork, error)
}

type WhiteList struct {
	stor WhiteListStore
	log  *zerolog.Logger
}

func NewWhiteList(stor WhiteListStore, log *zerolog.Logger) *WhiteList {
	return &WhiteList{stor: stor, log: log}
}

func (w *WhiteList) AddIP(network entity.IpNetwork) error {
	w.log.Info().Msg("Get prefix")
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = w.stor.Add(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteList) DeleteIP(network entity.IpNetwork) error {
	w.log.Info().Msg("Get prefix")
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = w.stor.Remove(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteList) GetIPs() ([]entity.IpNetwork, error) {
	ipList, err := w.stor.GetIPs()
	if err != nil {
		return nil, err
	}
	return ipList, nil
}
