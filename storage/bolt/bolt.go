package bolt

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/asdine/storm/v3"
	"log"
)

type Storage struct {
	ipdb *storm.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	return &Storage{ipdb: db}, nil
}

func (s *Storage) Close() error {
	return s.ipdb.Close()
}

func (s *Storage) SaveIpCache(ipCache storage.IpCache) error {
	return s.ipdb.Save(&ipCache)
}

func (s *Storage) GetIpCache(ip string) (*storage.IpCache, error) {
	var ipCache storage.IpCache
	err := s.ipdb.One("Ip", ip, &ipCache)
	if err != nil {
		return nil, err
	}

	return &ipCache, nil
}

func InitDatabase(config cmd.Configs) {
	// 创建一个新的存储实例
	var err error
	cmd.DB, err = NewStorage(config.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *Storage) {
		err := db.Close()
		if err != nil {
			//TODO add errorLog
		}
	}(cmd.DB)
}
