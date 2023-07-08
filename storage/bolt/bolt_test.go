package bolt

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/boltdb/bolt"
	"log"
	"testing"
)

func TestOutputDb(t *testing.T) {
	var testConfig cmd.Configs
	testConfig.DBFilePath = "../../data/database.db"
	InitDatabase()
	db, err := bolt.Open("database.ipdb", 0666, nil)
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

func TestStorage_UpdateCache(t *testing.T) {

}
