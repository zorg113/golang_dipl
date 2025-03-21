package service

import (
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
)

type WhiteListStore interface {
	AddIP(prefix, mask string) error
	DeleteIP(prefix, mask string) error
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
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = w.stor.AddIP(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteList) DeleteIP(network entity.IpNetwork) error {
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = w.stor.DeleteIP(prefix, network.Mask)
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
