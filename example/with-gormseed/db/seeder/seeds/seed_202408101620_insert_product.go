package seeds

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// don't rename the function name
func (s *Seeds) SeedInsertProduct(db *gorm.DB) error {

	fmt.Println("insert product")

	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	products := []*Product{
		{
			ID:   uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
			Name: "Amazon Fire TV Stick",
		},
		{
			ID:   uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
			Name: "Samsung Galaxy Tab S6",
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

func (s *Seeds) RollbackInsertProduct(db *gorm.DB) error {
	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	productIDs := []uuid.UUID{
		uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
		uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
	}

	result := db.Unscoped().Delete(&Product{}, productIDs)

	return result.Error
}
