package bolt

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/asdine/storm/v3"
	"log"
)

var DB *Storage

type Storage struct {
	Ipdb *storm.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	return &Storage{Ipdb: db}, nil
}

func (s *Storage) Close() error {
	return s.Ipdb.Close()
}

func (s *Storage) SaveIpCache(ipCache storage.IpCache) error {
	return s.Ipdb.Save(&ipCache)
}

func (s *Storage) GetIpCache(ip string) (*storage.IpCache, error) {
	var ipCache storage.IpCache
	err := s.Ipdb.One("Ip", ip, &ipCache)
	if err != nil {
		return nil, err
	}

	return &ipCache, nil
}

func (s *Storage) UpdateCache(ipCache storage.IpCache) error {
	if _, err := s.GetIpCache(ipCache.Ip); err != nil {
		return s.SaveIpCache(ipCache)
	}
	return s.Ipdb.Update(&ipCache)

	//TODO further check
}

func (s *Storage) SaveIpCacheAsync(ipCache storage.IpCache, resultChan chan error) {
	go func() {
		resultChan <- s.Ipdb.Save(&ipCache)
	}()
}

func (s *Storage) GetIpCacheAsync(ip string, resultChan chan *storage.IpCache, errChan chan error) {
	go func() {
		var ipCache storage.IpCache
		err := s.Ipdb.One("Ip", ip, &ipCache)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- &ipCache
		}
	}()
}

func (s *Storage) UpdateCacheAsync(ipCache storage.IpCache, resultChan chan error) {
	go func() {
		if _, err := s.GetIpCache(ipCache.Ip); err != nil {
			resultChan <- s.Ipdb.Save(&ipCache)
		} else {
			resultChan <- s.Ipdb.Update(&ipCache)
		}
	}()
}

func (s *Storage) UpdateServiceInfo(host string, port int, services *cmd.PortInfo) error {
	ipCache, err := s.GetIpCache(host)
	if err != nil {
		return err
	}

	for i := 0; i < len(ipCache.Services); i++ {
		if ipCache.Services[i].Port == port {
			ipCache.Services[i] = &storage.PortInfoStore{
				PortInfo: services,
				Banner:   ipCache.Services[i].Banner,
			}
			break
		}
	}
	return s.Ipdb.Update(ipCache)
}

func InitDatabase() {
	// 创建一个新的存储实例
	var err error
	DB, err = NewStorage(cmd.Config.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDatabase() {
	defer func(db *Storage) {
		err := db.Close()
		if err != nil {
			//TODO add errorLog
		}
	}(DB)
}
