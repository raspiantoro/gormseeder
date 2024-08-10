# Gormseeder

Inspired by [gormigrate](https://github.com/go-gormigrate/gormigrate), Gormseeder is a helper for the Gorm ORM. However, instead of handling database migrations, Gormseeder is used as a database seeder to populate initial data for your database. 

Gormseeder is a powerful and flexible library designed to enhance the Gorm ORM in Golang by providing an easy-to-use framework for database seeding. As an essential companion to Gorm, Gormseeder helps developers automate the process of populating databases with initial data, making it easier to set up development and testing environments, or to deploy production-ready databases with pre-filled data.

Key Features:
- Simple Seeder Management: Gormseeder streamlines the creation and management of seeder files, allowing developers to define and organize their data seeds in a clean and maintainable way. The library handles the execution of seeders in a consistent order, ensuring that your database is populated correctly every time.

- Rollback Support: Mistakes happen, and when they do, Gormseeder has you covered with its rollback functionality. This feature allows you to revert specific seeders or all of them, giving you full control over the state of your database at any point in time.

- Customizable Seeder Logic: Gormseeder provides the flexibility to define complex seeder logic, accommodating everything from simple data inserts to more sophisticated data relationships. You can customize your seeders to meet the unique needs of your application.

- CLI Utility: The accompanying CLI tool, **Gormseed**, simplifies the process of generating and managing seeder files. With a few commands, you can create new seeders, run them, or roll them back, making it easy to maintain your seed data over the course of your project's development.

- Environment-Specific Seeding: Gormseeder supports environment-specific seeders, allowing you to define different seed data for development, testing, and production environments. This ensures that each environment is populated with the appropriate data for its purpose.

Use Cases:
- Initial Data Population: Quickly populate your database with the necessary initial data when setting up a new environment.
- Testing: Create specific datasets for testing purposes, ensuring your tests run against consistent and predictable data.
- Demo Data: Seed your production database with demo data, enabling you to showcase features or onboard new users with pre-populated content.

Gormseeder is particularly useful for developers looking to enhance their Gorm-based applications with robust, organized, and maintainable data seeding capabilities. By automating the process of database seeding, Gormseeder allows you to focus more on building your application and less on managing your database's data state. Whether you are developing a new application or maintaining an existing one, Gormseeder provides the tools you need to ensure your database is always in the right state.

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
