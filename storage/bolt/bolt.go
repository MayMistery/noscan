package bolt

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/asdine/storm/v3"
	"log"
	"sync"
)

var (
	DB     *Storage
	DBPool *cmd.Pool
)

type Storage struct {
	Ipdb *storm.DB
	//mu   sync.RWMutex
}

func NewStorage(path string) (*Storage, error) {
	db, err := storm.Open(path)
	if err != nil {
		cmd.ErrLog("Open db fail %v", err)
		return nil, err
	}

	return &Storage{Ipdb: db}, nil
}

func (s *Storage) Close() error {
	return s.Ipdb.Close()
}

func (s *Storage) SaveIpCache(ipCache *storage.IpCache) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	err := s.Ipdb.Save(ipCache)
	if err != nil {
		cmd.ErrLog("save ip to db fail %v", err)
	}
	return err
}

func (s *Storage) GetIpCache(ip string) (*storage.IpCache, error) {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

	var ipCache storage.IpCache
	err := s.Ipdb.One("Ip", ip, &ipCache)
	if err != nil {
		//cmd.ErrLog("Get ip from db fail %v", err)
		return nil, err
	}

	return &ipCache, nil
}

func (s *Storage) UpdateCache(ipCache *storage.IpCache) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	if _, err := s.GetIpCache(ipCache.Ip); err != nil {
		return s.SaveIpCache(ipCache)
	}
	err := s.Ipdb.Update(ipCache)

	if err != nil {
		cmd.ErrLog("update db fail %v", err)
	}

	return err
}

// UpdateCacheAsync pushes an UpdateCache action to the DBPool.
func UpdateCacheAsync(ipCache *storage.IpCache) {
	DBPool.Push(poolInput{
		action: "UpdateCache",
		args:   ipCache,
	})
}

// UpdateServiceInfo updates the services of a host in the Bolt database.
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
			ipCache.Services[i] = services
			break
		}
	}
	return s.Ipdb.Update(ipCache)
}

// UpdateServiceInfoAsync pushes an UpdateServiceInfo action to the DBPool.
func UpdateServiceInfoAsync(host string, port int, services *cmd.PortInfo) {
	DBPool.Push(poolInput{
		action: "UpdateServiceInfo",
		args: &serviceInfoInput{
			host, port, services,
		},
	})
}

// UpdateDeviceInfo updates the device info of a host in the Bolt database.
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

// UpdateDeviceInfoAsync pushes an UpdateDeviceInfo action to the DBPool.
func UpdateDeviceInfoAsync(host string, deviceInfo string) {
	DBPool.Push(poolInput{
		action: "UpdateDeviceInfo",
		args: &deviceInfoInput{
			host, deviceInfo,
		},
	})
}

// UpdateHoneypot updates the honeypot of a host in the Bolt database.
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

// UpdateHoneypotAsync pushes an UpdateHoneypot action to the DBPool.
func UpdateHoneypotAsync(host string, honeypot []string) {
	DBPool.Push(poolInput{
		action: "UpdateHoneypot",
		args: &honeypotInput{
			host, honeypot,
		},
	})
}

// InitAsyncDatabase initializes the asynchronous database.
// It opens a new Bolt database, creates a new DBPool, and starts running the DBPool.
func InitAsyncDatabase() *sync.WaitGroup {
	// 创建一个新的存储实例
	var err error
	wg := &sync.WaitGroup{}
	DB, err = NewStorage(cmd.Config.DBFilePath)
	DBPool = NewDBPool()
	wg.Add(1)
	go DBPool.Run()
	if err != nil {
		cmd.ErrLog("fail to InitAsyncDatabase %v", err)
		log.Fatal(err)
	}
	return wg
}

func CloseDatabase(wg *sync.WaitGroup) {
	defer func(db *Storage) {
		wg.Done()
		err := db.Close()
		if err != nil {
			cmd.ErrLog("fail to CloseDatabase %v", err)
			panic(err)
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

// NewDBPool creates a new DBPool.
// It sets the function of the pool to be a function that performs a database action based on the poolInput.
func NewDBPool() *cmd.Pool {
	dbPool := cmd.NewPool(cmd.Config.Threads/10 + 1)
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
		case "UpdateBannerCache":
			func() {
				err = DB.UpdateBannerCache(in.args.(*storage.BannerCache))
			}()
		}

		if err != nil {
			cmd.ErrLog("dbAction unmatched %v", err)
			//HandleDBError(err)
			return
		}
	}
	return dbPool
}

//func HandleDBError(err error) {
//
//
//	//TODO handle db error together
//}

// SaveBannerCache saves a BannerCache instance to the Bolt database.
func (s *Storage) SaveBannerCache(bannerCache *storage.BannerCache) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	err := s.Ipdb.Save(bannerCache)
	if err != nil {
		cmd.ErrLog("save banner to db fail %v", err)
	}
	return err
}

// GetBannerCache retrieves a BannerCache instance from the Bolt database by IP.
func (s *Storage) GetBannerCache(ip string) (*storage.BannerCache, error) {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

	var bannerCache storage.BannerCache
	err := s.Ipdb.One("Ip", ip, &bannerCache)
	if err != nil {
		//cmd.ErrLog("Get ip from db fail %v", err)
		return nil, err
	}

	return &bannerCache, nil
}

// UpdateBannerCache updates a BannerCache instance in the Bolt database.
func (s *Storage) UpdateBannerCache(bannerCache *storage.BannerCache) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	if _, err := s.GetBannerCache(bannerCache.Ip); err != nil {
		return s.SaveBannerCache(bannerCache)
	}
	err := s.Ipdb.Update(bannerCache)

	if err != nil {
		cmd.ErrLog("update banner db fail %v", err)
	}

	return err
}

// UpdateBannerCacheAsync pushes an UpdateBannerCache action to the DBPool.
func UpdateBannerCacheAsync(bannerCache *storage.BannerCache) {
	DBPool.Push(poolInput{
		action: "UpdateBannerCache",
		args:   bannerCache,
	})
}
