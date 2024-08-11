package main

import (
	"log"
	"os"

	"github.com/raspiantoro/gormseeder"
	"github.com/raspiantoro/gormseeder/example/with-gormseed/db/seeder/seeds"
	"github.com/raspiantoro/gormseeder/gormseed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	var command *string

	if len(os.Args) > 1 {
		command = &os.Args[1]
	}

	dsn := "host=localhost user=admin password=admin dbname=marketplace port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	seeds := gormseed.Load(&seeds.Seeds{})

	seeder := gormseeder.New(db, seeds)

	if command != nil && *command == "rollback" {
		if err = seeder.Rollback(); err != nil {
			log.Fatalln(err)
		}
	} else {
		if err = seeder.Seed(); err != nil {
			log.Fatalln(err)
		}
	}
}
