package service

import (
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
)

type BlackListStore interface {
	AddIP(prefix, mask string) error
	DeleteIP(Prefix, mask string) error
	GetIPs() ([]entity.IpNetwork, error)
}

type BlackList struct {
	stor BlackListStore
	log  *zerolog.Logger
}

func NewBlackList(stor BlackListStore, log *zerolog.Logger) *BlackList {
	return &BlackList{stor: stor, log: log}
}

func (b *BlackList) AddIP(network entity.IpNetwork) error {
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = b.stor.AddIP(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlackList) DeleteIP(network entity.IpNetwork) error {
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = b.stor.DeleteIP(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlackList) GetIPs() ([]entity.IpNetwork, error) {
	ips, err := b.stor.GetIPs()
	if err != nil {
		return nil, err
	}
	return ips, nil
}
