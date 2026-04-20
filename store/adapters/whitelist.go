//nolint:dupl
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
	client *client.PostgresSQL
}

func NewWhiteListStorage(client *client.PostgresSQL) *WhiteListStorage {
	return &WhiteListStorage{client: client}
}

func (w *WhiteListStorage) AddIP(prefix, mask string) error {
	var isExist bool
	err := w.client.DB.QueryRow(isIPExistWiteList, prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return common.IPAlreadyExist
	}
	_, err = w.client.DB.Exec(insertIPInWiteList, prefix, mask)
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteListStorage) DeleteIP(prefix, mask string) error {
	_, err := w.client.DB.Exec(deleteIPFromWhiteList, prefix, mask)
	if err != nil {
		return err
	}
	return nil
}

func (w *WhiteListStorage) GetIPs() ([]entity.IPNetwork, error) {
	ipNetworkList := make([]entity.IPNetwork, 0, 5)
	err := w.client.DB.Select(&ipNetworkList, getAllIPFromWhiteList)
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
