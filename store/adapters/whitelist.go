package adapters

import (
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/store/client"
)

const (
	isIPExistWiteList     = `SELECT exists(SELECT 1 FROM whitelist WHERE prefix = $1 AND mask = $2)`
	insertIPInWiteList    = `INSERT INTO whitelist (prefix, mask) VALUES ($1, $2)`
	deleteIPFromWhiteList = `DELETE FROM whitelist WHERE prefix = $1 AND mask = $2`
	getAllIPFromWhiteList = `SELECT prefix, mask FROM whitelist`
)

type WhiteListStorage struct {
	client *client.PostgresSql
}

func NewWhiteListStorage(client *client.PostgresSql) *WhiteListStorage {
	return &WhiteListStorage{client: client}
}

func (w *WhiteListStorage) AddIP(prefix, mask string) error {
	var isExist bool
	err := w.client.Db.QueryRow(isIPExistWiteList, prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return common.IpAlreadyExist
	}
	err = w.client.Db.QueryRow(insertIPInWiteList, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteListStorage) DeleteIP(prefix, mask string) error {
	err := w.client.Db.QueryRow(deleteIPFromWhiteList, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteListStorage) GetIPs() ([]entity.IpNetwork, error) {
	ipNetworkList := make([]entity.IpNetwork, 0, 5)
	err := w.client.Db.Select(&ipNetworkList, getAllIPFromWhiteList)
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
