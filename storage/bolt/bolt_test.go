package bolt

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/lib/scanlib"
	"github.com/MayMistery/noscan/storage"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestOutputDb(t *testing.T) {
	cmd.Config.DBFilePath = "../../data/database.db"
	//InitAsyncDatabase()
	db, err := bolt.Open(cmd.Config.DBFilePath, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		// 遍历所有的桶
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("Bucket:", string(name))

			// 遍历桶中的键值对
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf("  Key: %s, Value: %s\n", k, v)
				return nil
			})
		})
	})

	if err != nil {
		log.Fatal(err)
	}
}

func TestUpdateCache(t *testing.T) {
	// 创建一个临时的测试数据库文件
	testDBPath := "../../data/database.db"
	//defer func() {
	//	// 清理测试数据库文件
	//	err := cleanupDBFile(testDBPath)
	//	if err != nil {
	//		t.Errorf("Failed to cleanup test database file: %v", err)
	//	}
	//}()

	// 创建存储实例并打开临时数据库文件
	db, err := NewStorage(testDBPath)
	if err != nil {
		t.Errorf("Failed to create storage instance: %v", err)
		return
	}
	defer db.Close()

	// 准备测试数据
	ip := "127.0.0.1"
	portInfoStore := &storage.PortInfoStore{
		PortInfo: &cmd.PortInfo{
			Port:       8080,
			Protocol:   "tcp",
			ServiceApp: []string{"http"},
		},
		Banner: &scanlib.Response{},
	}
	ipCache := storage.IpCache{
		Ip:         ip,
		Services:   []*storage.PortInfoStore{portInfoStore},
		DeviceInfo: "Device 1",
		Honeypot:   []string{"Honeypot 1"},
	}

	// 保存测试数据到数据库
	err = db.SaveIpCache(&ipCache)
	if err != nil {
		t.Errorf("Failed to save test data: %v", err)
		return
	}

	// 修改测试数据
	ipCache.DeviceInfo = "Updated Device"

	// 更新缓存
	err = db.UpdateCache(&ipCache)
	if err != nil {
		t.Errorf("Failed to update cache: %v", err)
		return
	}

	// 从数据库中获取更新后的数据
	updatedCache, err := db.GetIpCache(ip)
	if err != nil {
		t.Errorf("Failed to get updated cache: %v", err)
		return
	}

	// 验证数据是否正确更新
	assert.Equal(t, "Updated Device", updatedCache.DeviceInfo, "Device info should be updated")
}

func cleanupDBFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func TestInitDatabase(t *testing.T) {
	cmd.Config.DBFilePath = "../../data/database.db"
	InitAsyncDatabase()
}

func (s *Storage) TestStorage_Async(t *testing.T) {
	// 用于接收结果的channel
	//saveResultChan := make(chan error, 100)
	//getResultChan := make(chan *storage.IpCache, 100)
	//updateResultChan := make(chan error)
	//
	//// Save
	//go s.SaveIpCacheAsync(ipCache, saveResultChan)
	//
	//// Get
	//go s.GetIpCacheAsync(ip, getResultChan, errChan)
	//
	//// Update
	//go s.UpdateCacheAsync(ipCache, updateResultChan)

	// result
	//saveError := <-saveResultChan
	//ipCacheResult := <-getResultChan
	//getError := <-errChan
	//updateError := <-updateResultChan

}

func TestHighConcurrency(t *testing.T) {
	// 创建一个临时的测试数据库文件
	//defer func() {
	//	// 清理测试数据库文件
	//	err := cleanupDBFile(testDBPath)
	//	if err != nil {
	//		t.Errorf("Failed to cleanup test database file: %v", err)
	//	}
	//}()

	// 创建存储实例并打开临时数据库文件
	//db, err := NewStorage(testDBPath)
	cmd.Config.Threads = 50000
	cmd.Config.DBFilePath = "../../data/database.db"
	InitAsyncDatabase()
	ip := "127.0.0.2"
	portInfoStore := &storage.PortInfoStore{
		PortInfo: &cmd.PortInfo{
			Port:       8080,
			Protocol:   "tcp",
			ServiceApp: []string{"http"},
		},
		Banner: &scanlib.Response{},
	}
	ipCache := storage.IpCache{
		Ip:         ip,
		Services:   []*storage.PortInfoStore{portInfoStore},
		DeviceInfo: "Device 1",
		Honeypot:   []string{"Honeypot 1"},
	}
	println(ipCache.Ip)
	for {
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
		UpdateCacheAsync(&ipCache)
	}

	time.Sleep(10000)
}
