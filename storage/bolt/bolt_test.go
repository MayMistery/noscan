package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"testing"
)

func TestDb(t *testing.T) {
	CreatDatabase()
	db, err := bolt.Open("database.db", 0666, nil)
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
