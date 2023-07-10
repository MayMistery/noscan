package bolt

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/utils"
	"github.com/asdine/storm/v3"
	"log"
)

var (
	DB     *Storage
	DBPool *utils.Pool
)

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

func (s *Storage) SaveIpCache(ipCache *storage.IpCache) error {
	return s.Ipdb.Save(ipCache)
}

func (s *Storage) GetIpCache(ip string) (*storage.IpCache, error) {
	var ipCache storage.IpCache
	err := s.Ipdb.One("Ip", ip, &ipCache)
	if err != nil {
		return nil, err
	}

	return &ipCache, nil
}

func (s *Storage) UpdateCache(ipCache *storage.IpCache) error {
	if _, err := s.GetIpCache(ipCache.Ip); err != nil {
		return s.SaveIpCache(ipCache)
	}
	return s.Ipdb.Update(ipCache)

	//TODO further check
}

func UpdateCacheAsync(ipCache *storage.IpCache) {
	DBPool.Push(poolInput{
		action: "UpdateCache",
		args:   ipCache,
	})
}

func (s *Storage) UpdateServiceInfo(serviceInfoArg *serviceInfoInput) error {
	host := serviceInfoArg.host
	port := serviceInfoArg.port
	services := serviceInfoArg.services
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

func UpdateServiceInfoAsync(host string, port int, services *cmd.PortInfo) {
	DBPool.Push(poolInput{
		action: "UpdateServiceInfo",
		args: &serviceInfoInput{
			host, port, services,
		},
	})
}

func (s *Storage) UpdateDeviceInfo(deviceInfoArg *deviceInfoInput) error {
	host := deviceInfoArg.host
	deviceInfo := deviceInfoArg.deviceInfo
	ipCache, err := s.GetIpCache(host)
	if err != nil {
		return s.SaveIpCache(&storage.IpCache{DeviceInfo: deviceInfo})
	}

	ipCache.DeviceInfo = deviceInfo
	return s.UpdateCache(ipCache)
}

func UpdateDeviceInfoAsync(host string, deviceInfo string) {
	DBPool.Push(poolInput{
		action: "UpdateDeviceInfo",
		args: &deviceInfoInput{
			host, deviceInfo,
		},
	})
}

func (s *Storage) UpdateHoneypot(deviceInfoArg *honeypotInput) error {
	host := deviceInfoArg.host
	honeypot := deviceInfoArg.honeypot
	ipCache, err := s.GetIpCache(host)
	if err != nil {
		return s.SaveIpCache(&storage.IpCache{Honeypot: honeypot})
	}

	ipCache.Honeypot = honeypot
	return s.UpdateCache(ipCache)
}

func UpdateHoneypotAsync(host string, honeypot []string) {
	DBPool.Push(poolInput{
		action: "UpdateHoneypot",
		args: &honeypotInput{
			host, honeypot,
		},
	})
}

func InitDatabase() {
	// 创建一个新的存储实例
	var err error
	DB, err = NewStorage(cmd.Config.DBFilePath)
	DBPool = NewDBPool()
	go DBPool.Run()
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

type poolInput struct {
	action string
	args   interface{}
}

type serviceInfoInput struct {
	host     string
	port     int
	services *cmd.PortInfo
}

type deviceInfoInput struct {
	host       string
	deviceInfo string
}

type honeypotInput struct {
	host     string
	honeypot []string
}

func NewDBPool() *utils.Pool {
	dbPool := utils.NewPool(cmd.Config.Threads/10 + 1)
	dbPool.Function = func(input interface{}) {
		in := input.(poolInput)
		var err error
		switch in.action {
		case "UpdateServiceInfo":
			func() {
				serviceInfoArg := in.args.(*serviceInfoInput)
				err = DB.UpdateServiceInfo(serviceInfoArg)
			}()
		case "UpdateCache":
			func() {
				err = DB.UpdateCache(in.args.(*storage.IpCache))
			}()
		case "UpdateDeviceInfo":
			func() {
				err = DB.UpdateDeviceInfo(in.args.(*deviceInfoInput))
			}()
		case "UpdateHoneypot":
			func() {
				err = DB.UpdateHoneypot(in.args.(*honeypotInput))
			}()
		}

		if err != nil {
			HandleDBError(err)
			return
		}
	}
	return dbPool
}

func HandleDBError(err error) {
	//TODO handle db error
}
