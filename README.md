# Gormseeder

Inspired by [gormigrate](https://github.com/go-gormigrate/gormigrate), gormseeders is a helper for the Gorm ORM. However, instead of handling database migrations, gormseeders is used as a database seeder to populate initial data for your database.

## Requirements
```
Go 1.18
```

## How to install
```bash
go get github.com/raspiantoro/gormseeder
```

## Example usage
```go
package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/raspiantoro/gormseeder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "host=localhost user=admin password=admin dbname=product port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	type Product struct {
		gorm.Model
		ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		ProductName string
	}

	s := gormseeders.New(db, []*gormseeders.Seeders{
		{
			Key:  "202408092043", // using the datetime format 'YYYYMMDDhhmm', you can use timestamp string
			Name: "insert product",
			Seed: func(tx *gorm.DB) (err error) {
				products := []*Product{
					{
						ID:          uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
						ProductName: "SAMSUNG 990 PRO SSD NVMe M.2 PCIe Gen4",
					},
					{
						ID:          uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
						ProductName: "SanDisk 128GB Extreme PRO SDXC UHS-I Memory Card",
					},
				}

				for _, product := range products {
					result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&product)
					if result.Error != nil {
						return result.Error
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) (err error) {
				productIDs := []uuid.UUID{
					uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
					uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
				}

				result := tx.Unscoped().Delete(&Product{}, productIDs)

				return result.Error
			},
		},
	})
	if err := s.Seed(); err != nil {
		log.Fatalf("seeders failed: %v", err)
	}

	log.Println("seeders did run successfully")
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
