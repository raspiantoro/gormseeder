package seeds

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// don't rename the function name
func (s *Seeds) SeedMoreProduct(db *gorm.DB) error {
	fmt.Println("mode product")

	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	products := []*Product{
		{
			ID:   uuid.MustParse("31dadeea-5702-4d84-a9f3-2d783c12be82"),
			Name: "Apple 2024 Macbook Air",
		},
	}

	for _, product := range products {
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&product)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *Seeds) RollbackMoreProduct(db *gorm.DB) error {
	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	productIDs := []uuid.UUID{
		uuid.MustParse("31dadeea-5702-4d84-a9f3-2d783c12be82"),
	}

	result := db.Unscoped().Delete(&Product{}, productIDs)

	return result.Error
}
