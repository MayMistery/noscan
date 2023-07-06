package bolt

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/asdine/storm/v3"
	"log"
)

type Storage struct {
	db *storm.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveIpCache(ipCache storage.IpCache) error {
	return s.db.Save(&ipCache)
}

func (s *Storage) GetIpCache(ip string) (*storage.IpCache, error) {
	var ipCache storage.IpCache
	err := s.db.One("Ip", ip, &ipCache)
	if err != nil {
		return nil, err
	}

	return &ipCache, nil
}

func CreatDatabase() {
	// 创建一个新的存储实例
	db, err := NewStorage("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *Storage) {
		err := db.Close()
		if err != nil {
			//TODO add errorLog
		}
	}(db)

	// 存储示例数据
	ipCache := storage.IpCache{
		Ip:     "192.168.0.1",
		IpInfo: cmd.IpInfo{},
		Mark:   true,
	}
	err = db.SaveIpCache(ipCache)
	if err != nil {
		log.Fatal(err)
	}

	// 获取存储的数据
	ip := "192.168.0.1"
	cachedIp, err := db.GetIpCache(ip)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cached IP:", cachedIp.Ip)

	// 其他存储操作...
}
