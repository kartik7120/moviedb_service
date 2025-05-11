package helper

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Conn *gorm.DB
}

func ConnectToDB() (*gorm.DB, error) {
	var count int64
	dsn := os.Getenv("DSN")

	for {
		// conn, err := openDB(dsn)
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		if err != nil {
			log.Debug("Postgres is not ready yet...")
			count++
		} else {
			log.Info("Connected to Postgres Successfully")
			return db, nil
		}

		if count > 10 {
			log.Error("Could not connect to Postgres")
			return nil, err
		}

		log.Debug("Backing off for 2 seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
