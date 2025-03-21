//nolint:dupl
package adapters

import (
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/store/client"
)

const (
	isIPExistBlackList    = `SELECT exists(SELECT 1 FROM blacklist WHERE prefix = $1 AND mask =$2)`
	insertIPInBlackList   = `INSERT INTO blacklist (prefix, mask) VALUES ($1, $2)`
	deleteIPFromBlackList = `DELETE FROM blacklist WHERE prefix = $1 AND mask = $2`
	getAllIPFromBlackList = `SELECT prefix, mask FROM blacklist`
)

type BlackListStorage struct {
	client *client.PostgresSQL
}

func NewBlackListStorage(client *client.PostgresSQL) *BlackListStorage {
	return &BlackListStorage{client: client}
}

func (b *BlackListStorage) AddIP(prefix, mask string) error {
	var isExist bool
	err := b.client.DB.QueryRow(isIPExistBlackList, prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return common.IPAlreadyExist
	}
	err = b.client.DB.QueryRow(insertIPInBlackList, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (b *BlackListStorage) DeleteIP(prefix, mask string) error {
	err := b.client.DB.QueryRow(deleteIPFromBlackList, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (b *BlackListStorage) GetIPs() ([]entity.IPNetwork, error) {
	ipNetworkList := make([]entity.IPNetwork, 0, 5)
	err := b.client.DB.Select(&ipNetworkList, getAllIPFromBlackList)
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
